package repository

import "go.mongodb.org/mongo-driver/mongo"

// Basically the connector has a GetClient which would allow the app access to
// all kinds of dbs, it's a small help but this prevents us accessing anything
// other than the organisations db from the repos.
type MongoDbConnector interface {
	GetOrganisationsDb() (*mongo.Database, error)
}