package repository

import (
	"github.com/adamkirk-stayaway/organisations/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// Basically the connector has a GetClient which would allow the app access to
// all kinds of dbs, it's a small help but this prevents us accessing anything
// other than the organisations db from the repos.
type MongoDbConnector interface {
	GetOrganisationsDb() (*mongo.Database, error)
}

func getSortDirection(dir model.SortDirection) (int, error) {
	switch dir {
	case model.SortAsc:
		return 1, nil
	case model.SortDesc:
		return -1, nil
	default:
		return 0, model.ErrInvalidSortBy{
			Chosen: string(dir),
		}
	}
}