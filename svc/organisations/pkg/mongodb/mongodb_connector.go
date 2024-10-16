package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cachedDBs map[string]*mongo.Database

type Connector struct {
	opts               *options.ClientOptions
	connectionAttempts int

	client *mongo.Client
	dbs    cachedDBs
}

func (c *Connector) connect() (*mongo.Client, error) {
	if c.client != nil {
		return c.client, nil
	}
	var client *mongo.Client
	var err error

	for i := 0; i < c.connectionAttempts; i++ {
		// Create a new client and connect to the server
		client, err = mongo.Connect(context.TODO(), c.opts)

		if err != nil {
			continue
		}
	}

	if err != nil {
		c.client = client
	}

	return client, err
}

func (c *Connector) GetClient() (*mongo.Client, error) {
	return c.connect()
}

func (c *Connector) GetDB(name string) (*mongo.Database, error) {
	if db, ok := c.dbs[name]; ok && db != nil {
		return db, nil
	}

	client, err := c.connect()

	if err != nil {
		return nil, err
	}

	db := client.Database(name)

	c.dbs[name] = db

	return db, nil
}

type ConnectorOpt func(c *Connector)

func WithAttempts(attempts int) ConnectorOpt {
	if attempts < 1 {
		panic(errors.New("connection attempts cannot be less than 1"))
	}

	return func(c *Connector) {
		c.connectionAttempts = attempts
	}
}

func NewConnector(mongoOpts *options.ClientOptions, opts ...ConnectorOpt) *Connector {
	c := &Connector{
		opts:               mongoOpts,
		connectionAttempts: 1,
		dbs:                cachedDBs{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
