package repository

import (
	"context"
	"fmt"
	"log/slog"

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

func (r *MongoDbVenueAccommodations) ByNameAndVenueID(name string, venueId string) (*accommodations.Accommodation, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, err
	}

	org := &accommodations.Accommodation{}

	filter := bson.D{{
		"$and", bson.A{
			bson.D{{"name", name}},
			bson.D{{"venue_id", venueId}},
		},
	}}

	res := coll.FindOne(context.TODO(), filter)

	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
		return nil, common.ErrNotFound{
			ResourceName: "venue_accommodation",
			ID:           fmt.Sprintf("name:%s,venueid:%s", name, venueId),
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
