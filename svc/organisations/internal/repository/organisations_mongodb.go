package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/organisations"
	"github.com/adamkirk-stayaway/organisations/internal/util"
	"github.com/adamkirk-stayaway/organisations/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbOrganisations struct {
	connector *mongodb.Connector
	cfg MongoDBRepositoryConfig
}

func (r *MongoDbOrganisations) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetDB(r.cfg.MongoDbDatabase())

	if err != nil {
		return nil, err
	}

	coll := db.Collection(MongoDBCollections.Organisations)

	return coll, nil
}

func (r *MongoDbOrganisations) Paginate(orderBy organisations.SortBy, orderDir common.SortDirection, page int, perPage int) (organisations.Organisations, common.PaginationResult, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	// Consider estimated count, prefer it to be accurate though and once we use
	// filters this is no longer viable
	total, err := coll.CountDocuments(context.TODO(), bson.D{})

	mongoLimit := int64(perPage)
	mongoSkip := int64((page - 1)) * mongoLimit

	cursor, err := coll.Find(context.TODO(), bson.D{}, &options.FindOptions{Limit: &mongoLimit, Skip: &mongoSkip})

	orgs := &organisations.Organisations{}

	if err := cursor.All(context.TODO(), orgs); err != nil {
		return nil, common.PaginationResult{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return *orgs, common.PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *MongoDbOrganisations) Get(id string) (*organisations.Organisation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	org := &organisations.Organisation{}

	res := coll.FindOne(context.TODO(), bson.D{{"_id", objID}})

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "organisation",
			ID:           id,
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbOrganisations) BySlug(slug string) (*organisations.Organisation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	org := &organisations.Organisation{}

	res := coll.FindOne(context.TODO(), bson.D{{"slug", slug}})

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "organisation",
			ID:           fmt.Sprintf("slug:%s", slug),
		}
	}

	err = res.Decode(org)

	return org, err
}

func (r *MongoDbOrganisations) Delete(org *organisations.Organisation) error {
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

func (r *MongoDbOrganisations) Save(org *organisations.Organisation) (*organisations.Organisation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	objID, err := primitive.ObjectIDFromHex(org.ID)

	if err != nil {
		return nil, err
	}

	update := *org
	update.ID = ""

	_, err = coll.ReplaceOne(context.TODO(), bson.D{{"_id", objID}}, update, &options.ReplaceOptions{
		Upsert: util.PointTo(true),
	})


	if err != nil {
		return nil, err
	}

	return org, nil
}

func NewMongoDbOrganisations(connector *mongodb.Connector, cfg MongoDBRepositoryConfig) *MongoDbOrganisations {
	return &MongoDbOrganisations{
		connector: connector,
		cfg: cfg,
	}
}
