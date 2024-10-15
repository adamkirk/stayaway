package repository

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/domain/accommodations"
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/util"
	"github.com/adamkirk-stayaway/organisations/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbVenueAccommodationTemplates struct {
	connector *mongodb.Connector
	cfg MongoDBRepositoryConfig
}

func (r *MongoDbVenueAccommodationTemplates) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetDB(r.cfg.MongoDbDatabase())

	if err != nil {
		return nil, err
	}

	coll := db.Collection(MongoDBCollections.AccommodationVenueTemplates)

	return coll, nil
}

func (r *MongoDbVenueAccommodationTemplates) getSortColumn(sortBy accommodations.SortBy) (string, error) {
	switch sortBy {
	case accommodations.SortByName:
		return "name", nil
	default:
		return "", common.ErrInvalidSortBy{
			Chosen: string(sortBy),
		}
	}
}

func (r *MongoDbVenueAccommodationTemplates) filterToBsonD(search accommodations.SearchFilter) bson.D {
	filters := []bson.D{}

	if len(search.VenueID) > 0 {
		// TODO: this field is just a string, it probably wants to be ObbjectID for quicker lookups
		filters = append(filters, bson.D{{"venue_id", bson.D{{"$in", search.VenueID}}}})
	}

	if search.NamePrefix != nil {
		pattern := fmt.Sprintf("^%s\\.*", *search.NamePrefix)
		filters = append(filters, bson.D{{"template.name", bson.D{{"$regex", pattern}, {"$options", "i"}}}})
	}

	if len(filters) == 0 {
		return bson.D{{}}
	}

	return bson.D{{"$and", filters}}
}

func (r *MongoDbVenueAccommodationTemplates) Paginate(p accommodations.PaginationFilter, search accommodations.SearchFilter) (accommodations.VenueTemplates, common.PaginationResult, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	sortColumn, err := r.getSortColumn(p.OrderBy)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	sortDir, err := getSortDirection(p.OrderDir)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	opts := options.Find().
		SetLimit(int64(p.PerPage)).
		SetSkip(int64((p.Page - 1)) * int64(p.PerPage)).
		SetSort(bson.D{{sortColumn, sortDir}})

	filter := r.filterToBsonD(search)

	// Consider estimated count, prefer it to be accurate though and once we use
	// filters this is no longer viable
	total, err := coll.CountDocuments(context.TODO(), filter)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	cursor, err := coll.Find(context.TODO(), filter, opts)

	vts := &accommodations.VenueTemplates{}

	if err := cursor.All(context.TODO(), vts); err != nil {
		return nil, common.PaginationResult{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))

	return *vts, common.PaginationResult{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *MongoDbVenueAccommodationTemplates) Save(vt *accommodations.VenueTemplate) (*accommodations.VenueTemplate, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(vt.ID)

	if err != nil {
		return nil, err
	}

	update := *vt
	update.ID = ""

	_, err = coll.ReplaceOne(context.TODO(), bson.D{{"_id", objID}}, update, &options.ReplaceOptions{
		Upsert: util.PointTo(true),
	})


	if err != nil {
		return nil, err
	}

	return vt, nil
}


func (r *MongoDbVenueAccommodationTemplates) ByNameAndVenue(name string, venueId string) (*accommodations.VenueTemplate, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	vt := &accommodations.VenueTemplate{}

	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"template.name", name}},
			bson.D{{"venue_id", venueId}},
		},
	}}

	res := coll.FindOne(context.TODO(), filter)

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "venueaccommodationtemplate",
			ID:           fmt.Sprintf("name:%s,venue_id:%s", name, venueId),
		}
	}

	err = res.Decode(vt)

	return vt, err
}

func (r *MongoDbVenueAccommodationTemplates) Get(id string, venueId string) (*accommodations.VenueTemplate, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		slog.Warn("invalid id given", "id", id, "err", err)
		return nil, common.ErrNotFound{
			ResourceName: "venue",
			ID:           id,
		}
	}

	vt := &accommodations.VenueTemplate{}

	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"_id", objID}},
			bson.D{{"venue_id", venueId}},
		},
	}}

	res := coll.FindOne(context.TODO(), filter)

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "venueaccommodationtemplate",
			ID:           fmt.Sprintf("id:%s,venue_id:%s", id, venueId),
		}
	}

	err = res.Decode(vt)

	return vt, err
}

func (r *MongoDbVenueAccommodationTemplates) Delete(v *accommodations.VenueTemplate) error {
	coll, err := r.getCollection()

	if err != nil {
		return err
	}

	objID, err := primitive.ObjectIDFromHex(v.ID)

	if err != nil {
		return err
	}

	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objID}})

	return err
}

func NewMongoDbVenueAccommodationTemplates(connector *mongodb.Connector, cfg MongoDBRepositoryConfig) *MongoDbVenueAccommodationTemplates {
	return &MongoDbVenueAccommodationTemplates{
		connector: connector,
		cfg: cfg,
	}
}
