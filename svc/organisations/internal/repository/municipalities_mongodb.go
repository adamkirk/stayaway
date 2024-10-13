package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/adamkirk-stayaway/organisations/internal/domain/municipalities"
	"github.com/adamkirk-stayaway/organisations/internal/repository/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbMunicipalities struct {
	connector MongoDbConnector
}

func (r *MongoDbMunicipalities) getCollection() (*mongo.Collection, error) {
	db, err := r.connector.GetOrganisationsDb()

	if err != nil {
		return nil, err
	}

	coll := db.Collection(mongodb.Municipalities)

	return coll, nil
}

func (r *MongoDbMunicipalities) getSortColumn(opt municipalities.SortBy) (string, error) {
	switch opt {
	case municipalities.SortByName:
		return "name", nil
	default:
		return "", common.ErrInvalidSortBy{
			Chosen: string(opt),
		}
	}
}

func (r *MongoDbMunicipalities) filterToBsonD(search municipalities.SearchFilter) bson.D {
	filters := []bson.D{}

	if len(search.Country) > 0 {
		// TODO: this field is just a string, it probably wants to be ObbjectID for quicker lookups
		filters = append(filters, bson.D{{"country", bson.D{{"$in", search.Country}}}})
	}

	if search.NamePrefix != nil {
		// Consider whether this should use the ascii name, probably easier
		pattern := fmt.Sprintf("^%s\\.*", *search.NamePrefix)
		fmt.Printf("%s\n", pattern)
		filters = append(filters, bson.D{{"name", bson.D{{"$regex", pattern}, {"$options", "i"}}}})
	}

	if len(filters) == 0 {
		return bson.D{{}}
	}

	return bson.D{{"$and", filters}}
}

func (r *MongoDbMunicipalities) Paginate(p municipalities.PaginationFilter, search municipalities.SearchFilter) (municipalities.Municipalities, common.PaginationResult, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	sortColumn, err := r.getSortColumn(p.OrderBy)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	sortDir, err := getSortDirection(p.OrderDir)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	opts := options.Find().
		SetLimit(int64(p.PerPage)).
		SetSkip(int64((p.Page - 1)) * int64(p.PerPage)).
		SetSort(bson.D{{sortColumn, sortDir}})

	filter := r.filterToBsonD(search)

	// Consider estimated count, prefer it to be accurate though and once we use
	// filters this is no longer viable
	total, err := coll.CountDocuments(context.TODO(), filter)

	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	cursor, err := coll.Find(context.TODO(), filter, opts)

	municipalities := &municipalities.Municipalities{}

	if err := cursor.All(context.TODO(), municipalities); err != nil {
		return nil, common.PaginationResult{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))

	return *municipalities, common.PaginationResult{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

// func (r *MongoDbMunicipalities) Get(id string, orgId string) (*municipalities.Venue, error) {
// 	coll, err := r.getCollection()

// 	if err != nil {
// 		return nil, err
// 	}

// 	objID, err := primitive.ObjectIDFromHex(id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	org := &municipalities.Venue{}

// 	filter := bson.D{{
// 		"$and", bson.A{
// 			bson.D{{"_id", objID}},
// 			bson.D{{"organisation_id", orgId}},
// 		},
// 	}}

// 	res := coll.FindOne(context.TODO(), filter)

// 	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
// 		return nil, municipalities.ErrNotFound{
// 			ResourceName: "venue",
// 			ID: id,
// 		}
// 	}

// 	err = res.Decode(org)

// 	return org, err
// }

// func (r *MongoDbMunicipalities) Delete(v *municipalities.Venue) (error) {
// 	coll, err := r.getCollection()

// 	if err != nil {
// 		return err
// 	}

// 	objID, err := primitive.ObjectIDFromHex(v.ID)

// 	if err != nil {
// 		return err
// 	}

// 	_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objID}})

// 	return err
// }

func (r *MongoDbMunicipalities) UpdateBatch(batch []municipalities.Municipality) (municipalities.BatchUpdateResult, error) {

	coll, err := r.getCollection()

	if err != nil {
		return municipalities.BatchUpdateResult{}, err
	}

	models := []mongo.WriteModel{}

	for _, m := range batch {
		models = append(
			models,
			mongo.NewReplaceOneModel().SetFilter(bson.D{{"import_id", m.ImportID}}).SetUpsert(true).SetReplacement(m),
		)
	}

	res, err := coll.BulkWrite(context.TODO(), models, options.BulkWrite())

	return municipalities.BatchUpdateResult{
		Created: int(res.UpsertedCount),
		Updated: int(res.ModifiedCount),
	}, err

}

func NewMongoDbMunicipalities(connector MongoDbConnector) *MongoDbMunicipalities {
	return &MongoDbMunicipalities{
		connector: connector,
	}
}
