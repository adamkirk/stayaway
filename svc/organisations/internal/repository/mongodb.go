package repository

import (
	"context"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRepositoryConfig interface {
	MongoDbDatabase() string
}

var MongoDBCollections = struct {
	Organisations               string
	Venues                      string
	Municipalities              string
	AccommodationVenueTemplates string
	VenueAccommodations string
}{
	Organisations:               "organisations",
	Venues:                      "venues",
	Municipalities:              "municipalities",
	AccommodationVenueTemplates: "venue_accommodation_templates",
	VenueAccommodations: "venue_accommodations",
}

type Migration struct {
	name string `bson:"name"`
	up   func(mongodb.MigrationContext) error
	down func(mongodb.MigrationContext) error
}

func (m Migration) Name() string {
	return m.name
}

func (m Migration) Up(ctx mongodb.MigrationContext) error {
	return m.up(ctx)
}

func (m Migration) Down(ctx mongodb.MigrationContext) error {
	return m.down(ctx)
}

func AllMongoDBMigrations(dbName string) []mongodb.Migration {
	return []mongodb.Migration{
		Migration{
			name: "create-unique-indexes-on-orgs-collections",
			up: func(ctx mongodb.MigrationContext) error {

				coll := ctx.Client().Database(dbName).Collection(MongoDBCollections.Organisations)

				indexModel := mongo.IndexModel{
					Keys:    bson.D{{"slug", -1}},
					Options: options.Index().SetUnique(true).SetName("organisations-unique-slug"),
				}
				_, err := coll.Indexes().CreateOne(context.TODO(), indexModel)

				return err
			},
			down: func(ctx mongodb.MigrationContext) error {
				coll := ctx.Client().Database(dbName).Collection(MongoDBCollections.Organisations)

				_, err := coll.Indexes().DropOne(context.TODO(), "organisations-unique-slug")

				return err
			},
		},
		Migration{
			name: "create-index-venues-unique-slug-per-org",
			up: func(ctx mongodb.MigrationContext) error {
				coll := ctx.Client().Database(dbName).Collection(MongoDBCollections.Venues)

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
			down: func(ctx mongodb.MigrationContext) error {
				coll := ctx.Client().Database(dbName).Collection(MongoDBCollections.Venues)

				_, err := coll.Indexes().DropOne(context.TODO(), "venues-unique-slug-per-org")

				return err
			},
		},
	}
}

func getSortDirection(dir common.SortDirection) (int, error) {
	switch dir {
	case common.SortAsc:
		return 1, nil
	case common.SortDesc:
		return -1, nil
	default:
		return 0, common.ErrInvalidSortBy{
			Chosen: string(dir),
		}
	}
}
