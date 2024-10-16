package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type PrimitiveObjectIDGenerator struct{}

func (g *PrimitiveObjectIDGenerator) Generate() string {
	return primitive.NewObjectID().Hex()
}

func NewPrimitiveObjectIDGenerator() *PrimitiveObjectIDGenerator {
	return &PrimitiveObjectIDGenerator{}
}
