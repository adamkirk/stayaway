package migrations

import (
	"context"

	"github.com/adamkirk-stayaway/organisations/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbMigrationConfig interface {
	MongoDbMigrationsDatabase() string
	MongoDbDatabase() string
}

type Migration struct {
	Name string `bson:"name"`
	up   func(*mongo.Client, MongoDbMigrationConfig) error
	down func(*mongo.Client, MongoDbMigrationConfig) error
}

func (m Migration) Up(client *mongo.Client, cfg MongoDbMigrationConfig) error {
	return m.up(client, cfg)
}

func (m Migration) Down(client *mongo.Client, cfg MongoDbMigrationConfig) error {
	return m.down(client, cfg)
}

var AllMigrations = []Migration{
	{
		Name: "create-unique-indexes-on-orgs-collections",
		up: func(client *mongo.Client, cfg MongoDbMigrationConfig) error {
			coll := client.Database(cfg.MongoDbDatabase()).Collection(mongodb.Organisations)

			indexModel := mongo.IndexModel{
				Keys:    bson.D{{"slug", -1}},
				Options: options.Index().SetUnique(true).SetName("organisations-unique-slug"),
			}
			_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)

			return err
		},
		down: func(client *mongo.Client, cfg MongoDbMigrationConfig) error {
			coll := client.Database(cfg.MongoDbDatabase()).Collection(mongodb.Organisations)

			_, err := coll.Indexes().DropOne(context.TODO(), "organisations-unique-slug")

			return err
		},
	},
	{
		Name: "create-index-venues-unique-slug-per-org",
		up: func(client *mongo.Client, cfg MongoDbMigrationConfig) error {
			coll := client.Database(cfg.MongoDbDatabase()).Collection(mongodb.Venues)

			indexModel := mongo.IndexModel{
				Keys: bson.D{
					{"organisation_id", 1},
					{"slug", 1},
				},
				Options: options.Index().SetUnique(true).SetName("venues-unique-slug-per-org"),
			}
			_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)

			return err
		},
		down: func(client *mongo.Client, cfg MongoDbMigrationConfig) error {
			coll := client.Database(cfg.MongoDbDatabase()).Collection(mongodb.Venues)

			_, err := coll.Indexes().DropOne(context.TODO(), "venues-unique-slug-per-org")

			return err
		},
	},
	{
		Name: "create-index-venues-unique-slug-per-org",
		up: func(client *mongo.Client, cfg MongoDbMigrationConfig) error {
			coll := client.Database(cfg.MongoDbDatabase()).Collection(mongodb.Venues)

			indexModel := mongo.IndexModel{
				Keys: bson.D{
					{"organisation_id", 1},
					{"slug", 1},
				},
				Options: options.Index().SetUnique(true).SetName("venues-unique-slug-per-org"),
			}
			_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)

			return err
		},
		down: func(client *mongo.Client, cfg MongoDbMigrationConfig) error {
			coll := client.Database(cfg.MongoDbDatabase()).Collection(mongodb.Venues)

			_, err := coll.Indexes().DropOne(context.TODO(), "venues-unique-slug-per-org")

			return err
		},
	},
}
