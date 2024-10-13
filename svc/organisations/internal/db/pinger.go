package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Pinger interface {
	Ping() error
}

type MongoDbPinger struct {
	connector *MongoDbConnector
}

func (p *MongoDbPinger) Ping() error {
	db, err := p.connector.GetOrganisationsDb()

	if err != nil {
		return err
	}

	var result bson.M

	return db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result)
}

func NewMongoDbPinger(connector *MongoDbConnector) Pinger {
	return &MongoDbPinger{
		connector: connector,
	}
}
