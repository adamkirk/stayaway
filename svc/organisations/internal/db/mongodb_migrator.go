package db

import (
	"context"
	"fmt"

	"github.com/adamkirk-stayaway/organisations/internal/db/mongodb/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbMigratorConfig interface {
	 MongoDbMigrationsDatabase() string
	 MongoDbDatabase() string
}

type MongoDbMigrator struct {
	connector *MongoDbConnector
	cfg MongoDbMigratorConfig
}

func migrationsToApply(target string) ([]migrations.Migration, error) {
	before := []migrations.Migration{}

	for _, mig := range migrations.AllMigrations {
		before = append(before, mig)
		if mig.Name == target {
			return before, nil
		}
	}

	return nil, fmt.Errorf("didn't find migration: %s", target)
}

func migrationsToRevert(target string) ([]migrations.Migration, error) {
	after := []migrations.Migration{}

	for i := len(migrations.AllMigrations) -1; i >= 0; i-- {
		mig := migrations.AllMigrations[i]
		after = append(after, mig)
		if mig.Name == target {
			return after, nil
		}
	}

	return nil, fmt.Errorf("didn't find migration: %s", target)
}

func (m *MongoDbMigrator) Up(to string) error {
	client, err := m.connector.GetClient()
	coll := client.Database(m.cfg.MongoDbMigrationsDatabase()).Collection("applied_migrations")

	if err != nil {
		return err
	}

	migs := migrations.AllMigrations

	if to != "" {
		migs, err = migrationsToApply(to)

		if err != nil {
			return err
		}
	}

	for _, mig := range migs {
		res := coll.FindOne(context.TODO(), bson.D{{"name", mig.Name}})

		if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
			// have to apply
			fmt.Printf("Applying: %s\n", mig.Name)

			if err := mig.Up(client, m.cfg); err != nil {
				return fmt.Errorf("Error while applying %s: %w", mig.Name, err)
			}

			_, err := coll.InsertOne(context.TODO(), mig)
	
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("Skipping (already applied): %s\n", mig.Name)
		}
	}

	return nil
}

func (m *MongoDbMigrator) Down(to string) error {
	client, err := m.connector.GetClient()
	coll := client.Database(m.cfg.MongoDbMigrationsDatabase()).Collection("applied_migrations")

	if err != nil {
		return err
	}

	migs := migrations.AllMigrations

	if to != "" {
		migs, err = migrationsToRevert(to)

		if err != nil {
			return err
		}
	}


	for _, mig := range migs {
		res := coll.FindOne(context.TODO(), bson.D{{"name", mig.Name}})

		if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
			fmt.Printf("Skipping (not applied): %s\n", mig.Name)
		} else {
			// have to apply
			fmt.Printf("Reverting: %s\n", mig.Name)
			if err := mig.Down(client, m.cfg); err != nil {
				return fmt.Errorf("Error while reverting %s: %w", mig.Name, err)
			}
			_, err := coll.DeleteOne(context.TODO(), mig)
	
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewMongoDbMigrator(connector *MongoDbConnector, cfg MongoDbMigratorConfig) Migrator {
	return &MongoDbMigrator{
		connector: connector,
		cfg: cfg,
	}
}