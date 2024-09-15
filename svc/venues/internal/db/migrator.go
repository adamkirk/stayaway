package db

import "errors"

type Migrator interface {
	Migrate() error
}

type MongoDbMigrator struct {
	connector *MongoDbConnector
}

func (p *MongoDbMigrator) Migrate() error {
	return errors.New("not iplemented")
}

func NewMongoDbMigrator(connector *MongoDbConnector) Migrator {
	return &MongoDbMigrator{
		connector: connector,
	}
}