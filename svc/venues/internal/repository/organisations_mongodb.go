package repository

import (
	"context"
	"errors"
	"math"

	"github.com/adamkirk-stayaway/venues/internal/db"
	"github.com/adamkirk-stayaway/venues/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "organisations"

type MongoDbOrganisations struct {
	connector *db.MongoDbConnector
}

func (r *MongoDbOrganisations) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetOrganisationsDb()

	if err != nil {
		return nil, err
	}

	coll := db.Collection(collectionName)

	return coll, nil
}

func (r *MongoDbOrganisations) Paginate(orderBy model.OrganisationSortBy, orderDir model.SortDirection, page int, perPage int) (model.Organisations, model.PaginationResult, error) {
	// No longer than one
	coll, err := r.getCollection()

	if err != nil {
		return nil, model.PaginationResult{}, err
	}

	// Consider estimated count, prefer it to be accurate though and once we use 
	// filters this is no longer viable
	total, err := coll.CountDocuments(context.TODO(), bson.D{})

	mongoLimit := int64(perPage)
	mongoSkip := int64((page - 1)) * mongoLimit

	cursor, err := coll.Find(context.TODO(), bson.D{}, &options.FindOptions{Limit: &mongoLimit, Skip: &mongoSkip})

	orgs := &model.Organisations{}

	if err := cursor.All(context.TODO(), orgs); err != nil {
		return nil, model.PaginationResult{}, err
	}

	totalPages := int(math.Ceil(float64(total)/float64(perPage)))
	return *orgs, model.PaginationResult{
		Page: page,
		PerPage: perPage,
		Total: int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *MongoDbOrganisations) doInsert(org *model.Organisation) (*model.Organisation, error) {
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

func (r *MongoDbOrganisations) doUpdate(org *model.Organisation) (*model.Organisation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	_, err = coll.UpdateOne(context.TODO(), bson.D{{"_id", org.ID}}, org)

	if err != nil {
		return nil, err
	}

	return org, nil
}

func (r *MongoDbOrganisations) Save(org *model.Organisation) (*model.Organisation, error) {
	if org.ID == "" {
		return r.doInsert(org)
	}

	return r.doUpdate(org)

}

func NewMongoDbOrganisations(connector *db.MongoDbConnector ) *MongoDbOrganisations {
	return &MongoDbOrganisations{
		connector: connector,
	}
}