package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/adamkirk-stayaway/organisations/internal/model"
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


func (r *MongoDbMunicipalities) getSortColumn(opt model.MunicipalitySortBy) (string, error) {
	switch opt {
	case model.MunicipalitySortByName:
		return "name", nil
	default:
		return "", model.ErrInvalidSortBy{
			Chosen: string(opt),
		}
	}
}

func (r *MongoDbMunicipalities) filterToBsonD(search model.MunicipalitySearchFilter) bson.D {
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


func (r *MongoDbMunicipalities) Paginate(p model.MunicipalityPaginationFilter, search model.MunicipalitySearchFilter) (model.Municipalities, model.PaginationResult, error) {
	coll, err := r.getCollection()

	if err != nil {
		return nil, model.PaginationResult{}, err
	}

	sortColumn, err := r.getSortColumn(p.OrderBy)

	if err != nil {
		return nil, model.PaginationResult{}, err
	}

	sortDir, err := getSortDirection(p.OrderDir)

	if err != nil {
		return nil, model.PaginationResult{}, err
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
		return nil, model.PaginationResult{}, err
	}
	
	cursor, err := coll.Find(context.TODO(), filter, opts)

	municipalities := &model.Municipalities{}

	if err := cursor.All(context.TODO(), municipalities); err != nil {
		return nil, model.PaginationResult{}, err
	}

	totalPages := int(math.Ceil(float64(total)/float64(p.PerPage)))

	return *municipalities, model.PaginationResult{
		Page: p.Page,
		PerPage: p.PerPage,
		Total: int(total),
		TotalPages: totalPages,
	}, nil
}

// func (r *MongoDbMunicipalities) Get(id string, orgId string) (*model.Venue, error) {
// 	coll, err := r.getCollection()

// 	if err != nil {
// 		return nil, err
// 	}

// 	objID, err := primitive.ObjectIDFromHex(id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	org := &model.Venue{}

// 	filter := bson.D{{
// 		"$and", bson.A{
// 			bson.D{{"_id", objID}},
// 			bson.D{{"organisation_id", orgId}},
// 		},
// 	}}


// 	res := coll.FindOne(context.TODO(), filter)

// 	if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
// 		return nil, model.ErrNotFound{
// 			ResourceName: "venue",
// 			ID: id,
// 		}
// 	}

// 	err = res.Decode(org)

// 	return org, err
// }

// func (r *MongoDbMunicipalities) Delete(v *model.Venue) (error) {
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

func (r *MongoDbMunicipalities) UpdateBatch(batch []model.Municipality) (model.BatchUpdateResult, error) {
	
	coll, err := r.getCollection()
	
	if err != nil {
		return model.BatchUpdateResult{}, err
	}

	models := []mongo.WriteModel{}

	for _, m := range batch {
		models = append(
			models, 
			mongo.NewReplaceOneModel().SetFilter(bson.D{{"import_id", m.ImportID}}).SetUpsert(true).SetReplacement(m),
		)
	}

	res, err := coll.BulkWrite(context.TODO(), models, options.BulkWrite())

	return model.BatchUpdateResult{
		Created: int(res.UpsertedCount),
		Updated: int(res.ModifiedCount),
	}, err

}

func NewMongoDbMunicipalities(connector MongoDbConnector ) *MongoDbMunicipalities {
	return &MongoDbMunicipalities{
		connector: connector,
	}
}