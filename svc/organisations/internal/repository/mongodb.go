package repository

import (
	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"go.mongodb.org/mongo-driver/mongo"
)

// Basically the connector has a GetClient which would allow the app access to
// all kinds of dbs, it's a small help but this prevents us accessing anything
// other than the organisations db from the repos.
type MongoDbConnector interface {
	GetOrganisationsDb() (*mongo.Database, error)
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
