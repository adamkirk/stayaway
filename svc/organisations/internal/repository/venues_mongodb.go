package repository

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
	"github.com/adamkirk-stayaway/organisations/internal/util"
	"github.com/adamkirk-stayaway/organisations/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbVenues struct {
	connector *mongodb.Connector
	cfg       MongoDBRepositoryConfig
}

func (r *MongoDbVenues) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetDB(r.cfg.MongoDbDatabase())

	if err != nil {
		return nil, err
	}

	coll := db.Collection(MongoDBCollections.Venues)

	return coll, nil
}

func (r *MongoDbVenues) getSortColumn(sortBy venues.SortBy) (string, error) {
	switch sortBy {
	case venues.SortByName:
		return "name", nil
	case venues.SortBySlug:
		return "slug", nil
	default:
		return "", common.ErrInvalidSortBy{
			Chosen: string(sortBy),
		}
	}
}

func (r *MongoDbVenues) filterToBsonD(search venues.SearchFilter) bson.D {
	var orgFilter bson.D

	if len(search.OrganisationID) > 0 {
		// TODO: this field is just a string, it probably wants to be ObbjectID for quicker lookups
		orgFilter = bson.D{{"organisation_id", bson.D{{"$in", search.OrganisationID}}}}
	}

	if orgFilter == nil {
		return bson.D{{}}
	}

	return bson.D{{"$and", bson.A{
		orgFilter,
	}}}
}

func (r *MongoDbVenues) Paginate(p venues.PaginationFilter, search venues.SearchFilter) (venues.Venues, common.PaginationResult, error) {
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

	orgs := &venues.Venues{}

	if err := cursor.All(context.TODO(), orgs); err != nil {
		return nil, common.PaginationResult{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))

	return *orgs, common.PaginationResult{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *MongoDbVenues) Get(id string, orgId string) (*venues.Venue, error) {
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

	org := &venues.Venue{}

	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"_id", objID}},
			bson.D{{"organisation_id", orgId}},
		},
	}}

	res := coll.FindOne(context.TODO(), filter)

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "venue",
			ID:           id,
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbVenues) BySlugAndOrganisation(slug string, orgId string) (*venues.Venue, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	org := &venues.Venue{}

	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"slug", slug}},
			bson.D{{"organisation_id", orgId}},
		},
	}}

	res := coll.FindOne(context.TODO(), filter)

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "venue",
			ID:           fmt.Sprintf("slug:%s", slug),
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbVenues) Delete(v *venues.Venue) error {
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

func (r *MongoDbVenues) Save(v *venues.Venue) (*venues.Venue, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(v.ID)

	if err != nil {
		return nil, err
	}

	update := *v
	update.ID = ""

	_, err = coll.ReplaceOne(context.TODO(), bson.D{{"_id", objID}}, update, &options.ReplaceOptions{
		Upsert: util.PointTo(true),
	})

	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewMongoDbVenues(connector *mongodb.Connector, cfg MongoDBRepositoryConfig) *MongoDbVenues {
	return &MongoDbVenues{
		connector: connector,
		cfg:       cfg,
	}
}
