package db

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig interface {
	MongoDbUri() string
	MongoDbDatabase() string
	MongoDbConnectionRetries() int
}

type MongoDbConnector struct {
	orgsDb *mongo.Database
	client *mongo.Client
	uri string
	dbName string
	connectionRetries int
}

func (c *MongoDbConnector) connect() (*mongo.Client, error) {
	if c.client != nil {
		return c.client, nil
	}
	var client *mongo.Client
	var err error

	for i := range c.connectionRetries {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(c.uri).SetServerAPIOptions(serverAPI)
	
		// Create a new client and connect to the server
		client, err = mongo.Connect(context.TODO(), opts)
		
		if err != nil {
			slog.Warn("failed to connect to mongo", "attempt", i+1, "error", err)
			continue;
		}
	}

	if err != nil {
		c.client = client
	}

	return client, err
}

func (c *MongoDbConnector) GetOrganisationsDb() (*mongo.Database, error) {
	if c.orgsDb != nil {
		return c.orgsDb, nil
	}

	client, err := c.connect()

	if err != nil {
		return nil, err
	}

	return client.Database(c.dbName), nil
}

func NewMongoDbConnector(cfg MongoConfig) *MongoDbConnector {
	return &MongoDbConnector{
		uri: cfg.MongoDbUri(),
		dbName: cfg.MongoDbDatabase(),
		connectionRetries: cfg.MongoDbConnectionRetries(),
	}
}
