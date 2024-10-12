package repository

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/venues"
	"github.com/adamkirk-stayaway/organisations/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbVenues struct {
	connector MongoDbConnector
}

func (r *MongoDbVenues) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetOrganisationsDb()

	if err != nil {
		return nil, err
	}

	coll := db.Collection(mongodb.Venues)

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

	totalPages := int(math.Ceil(float64(total)/float64(p.PerPage)))

	return *orgs, common.PaginationResult{
		Page: p.Page,
		PerPage: p.PerPage,
		Total: int(total),
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
		return nil, err
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
			ID: id,
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbVenues) BySlugAndOrganisation(slug string, orgId string)(*venues.Venue, error) {
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
			ID: fmt.Sprintf("slug:%s", slug),
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbVenues) Delete(v *venues.Venue) (error) {
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

func (r *MongoDbVenues) doInsert(org *venues.Venue) (*venues.Venue, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	result, err := coll.InsertOne(context.TODO(), org)

	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		org.ID = oid.Hex()
	} else {
		return nil, errors.New("failed to get generated id for document")
	}

	return org, nil
}

func (r *MongoDbVenues) doUpdate(org *venues.Venue) (*venues.Venue, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(org.ID)

	if err != nil {
		return nil, err
	}

	// copy the value, and strip its id
	// Maybe i can do something with the bson conversion somewhere, but otherwise
	// cause the ID is a string it needs to be cast to primitive.ObjectID, and i
	// don't really wanna do that in the model itself, so that it can be agnostic
	// to db driver
	// This seems the simpler option.
	update := *org
	update.ID = ""

	_, err = coll.ReplaceOne(context.TODO(), bson.D{{"_id", objID}}, update)

	if err != nil {
		return nil, err
	}

	return org, nil
}

func (r *MongoDbVenues) Save(org *venues.Venue) (*venues.Venue, error) {
	if org.ID == "" {
		return r.doInsert(org)
	}

	return r.doUpdate(org)

}

func NewMongoDbVenues(connector MongoDbConnector ) *MongoDbVenues {
	return &MongoDbVenues{
		connector: connector,
	}
}