package repository

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/repository/mongodb"
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbOrganisations struct {
	connector MongoDbConnector
}

func (r *MongoDbOrganisations) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetOrganisationsDb()

	if err != nil {
		return nil, err
	}

	coll := db.Collection(mongodb.Organisations)

	return coll, nil
}

func (r *MongoDbOrganisations) Paginate(orderBy model.OrganisationSortBy, orderDir model.SortDirection, page int, perPage int) (model.Organisations, model.PaginationResult, error) {
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

func (r *MongoDbOrganisations) Get(id string) (*model.Organisation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	org := &model.Organisation{}

	res := coll.FindOne(context.TODO(), bson.D{{"_id", objID}})

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, model.ErrNotFound{
			ResourceName: "organisation",
			ID: id,
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbOrganisations) BySlug(slug string) (*model.Organisation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	org := &model.Organisation{}

	res := coll.FindOne(context.TODO(), bson.D{{"slug", slug}})

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, model.ErrNotFound{
			ResourceName: "organisation",
			ID: fmt.Sprintf("slug:%s", slug),
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbOrganisations) Delete(org *model.Organisation) (error) {
	coll, err := r.getCollection()

	if err != nil {
		return err
	}

	objID, err := primitive.ObjectIDFromHex(org.ID)

	if err != nil {
		return err
	}

	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objID}})

	return err
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

func (r *MongoDbOrganisations) Save(org *model.Organisation) (*model.Organisation, error) {
	if org.ID == "" {
		return r.doInsert(org)
	}

	return r.doUpdate(org)

}

func NewMongoDbOrganisations(connector MongoDbConnector ) *MongoDbOrganisations {
	return &MongoDbOrganisations{
		connector: connector,
	}
}