package repository

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues/accommodations"
	"github.com/adamkirk-stayaway/organisations/internal/util"
	"github.com/adamkirk-stayaway/organisations/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbVenueAccommodations struct {
	connector *mongodb.Connector
	cfg       MongoDBRepositoryConfig
}

func (r *MongoDbVenueAccommodations) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetDB(r.cfg.MongoDbDatabase())

	if err != nil {
		return nil, err
	}

	coll := db.Collection(MongoDBCollections.VenueAccommodations)

	return coll, nil
}


func (r *MongoDbVenueAccommodations) getSortColumn(sortBy accommodations.SortBy) (string, error) {
	switch sortBy {
	case accommodations.SortByReference:
		return "reference", nil
	default:
		return "", common.ErrInvalidSortBy{
			Chosen: string(sortBy),
		}
	}
}

func (r *MongoDbVenueAccommodations) filterToBsonD(search accommodations.SearchFilter) bson.D {
	filters := []bson.D{}

	if len(search.VenueID) > 0 {
		// TODO: this field is just a string, it probably wants to be ObbjectID for quicker lookups
		filters = append(filters, bson.D{{"venue_id", bson.D{{"$in", search.VenueID}}}})
	}

	if len(search.VenueTemplateID) > 0 {
		// TODO: this field is just a string, it probably wants to be ObbjectID for quicker lookups
		filters = append(filters, bson.D{{"venue_template_id", bson.D{{"$in", search.VenueTemplateID}}}})
	}

	if search.ReferencePrefix != nil {
		pattern := fmt.Sprintf("^%s\\.*", *search.ReferencePrefix)
		filters = append(filters, bson.D{{"reference", bson.D{{"$regex", pattern}, {"$options", "i"}}}})
	}

	if len(filters) == 0 {
		return bson.D{{}}
	}

	return bson.D{{"$and", filters}}
}

func (r *MongoDbVenueAccommodations) Paginate(p accommodations.PaginationFilter, search accommodations.SearchFilter) (accommodations.Accommodations, common.PaginationResult, error) {
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

	accs := &accommodations.Accommodations{}

	if err := cursor.All(context.TODO(), accs); err != nil {
		return nil, common.PaginationResult{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))

	return *accs, common.PaginationResult{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *MongoDbVenueAccommodations) Get(id string, venueId string) (*accommodations.Accommodation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		slog.Warn("invalid id given", "id", id)
		return nil, common.ErrNotFound{
			ResourceName: "venue",
			ID:           id,
		}
	}

	a := &accommodations.Accommodation{}

	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"_id", objID}},
			bson.D{{"venue_id", venueId}},
		},
	}}

	res := coll.FindOne(context.TODO(), filter)

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "accommodation",
			ID:           fmt.Sprintf("id:%s,venueid:%s", id, venueId),
		}
	}

	err = res.Decode(a)

	return a, err
}

func (r *MongoDbVenueAccommodations) ByReferenceAndVenueID(reference string, venueId string) (*accommodations.Accommodation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	org := &accommodations.Accommodation{}

	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"reference", reference}},
			bson.D{{"venue_id", venueId}},
		},
	}}

	res := coll.FindOne(context.TODO(), filter)

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "venue_accommodation",
			ID:           fmt.Sprintf("name:%s,venueid:%s", reference, venueId),
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbVenueAccommodations) Delete(a *accommodations.Accommodation) error {
	coll, err := r.getCollection()

	if err != nil {
		return err
	}

	objID, err := primitive.ObjectIDFromHex(a.ID)

	if err != nil {
		return err
	}

	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objID}})

	return err
}

func (r *MongoDbVenueAccommodations) Save(a *accommodations.Accommodation) (*accommodations.Accommodation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(a.ID)

	if err != nil {
		return nil, err
	}

	update := *a
	update.ID = ""
	update.Config = nil

	_, err = coll.ReplaceOne(context.TODO(), bson.D{{"_id", objID}}, update, &options.ReplaceOptions{
		Upsert: util.PointTo(true),
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}

func NewMongoDbVenueAccommodations(connector *mongodb.Connector, cfg MongoDBRepositoryConfig) *MongoDbVenueAccommodations {
	return &MongoDbVenueAccommodations{
		connector: connector,
		cfg:       cfg,
	}
}
