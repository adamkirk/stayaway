package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Pinger struct {
	connector *Connector
	db string
}

func (p *Pinger) Ping() error {
	db, err := p.connector.GetDB(p.db)

	if err != nil {
		return err
	}

	var result bson.M

	return db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result)
}

func NewPinger(connector *Connector, db string) *Pinger {
	return &Pinger{
		connector: connector,
		db: db,
	}
}
