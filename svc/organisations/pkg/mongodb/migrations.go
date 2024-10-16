package mongodb

import (
	"context"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MigrationContext interface {
	Client() *mongo.Client
	LogInfo(msg string) error
	LogError(msg string) error
}

type migrationContext struct {
	client   *mongo.Client
	logInfo  func(msg string) error
	logError func(msg string) error
}

func (c *migrationContext) Client() *mongo.Client {
	return c.client
}

func (c *migrationContext) LogInfo(msg string) error {
	return c.logInfo(msg)
}

func (c *migrationContext) LogError(msg string) error {
	return c.logError(msg)
}

type migrationRecord struct {
	Name string `bson:"name"`
}

type Migration interface {
	Name() string
	Up(MigrationContext) error
	Down(MigrationContext) error
}

type Migrator struct {
	connector    *Connector
	migrationsDB string
	migrations   []Migration
	infoIo       io.Writer
	errorIo      io.Writer
}

func (m *Migrator) logInfo(msg string) error {
	if m.infoIo == nil {
		return nil
	}

	_, err := fmt.Fprintf(m.infoIo, "[INFO] %s\n", msg)

	return err
}

func (m *Migrator) logError(msg string) error {
	if m.errorIo == nil {
		return nil
	}

	_, err := fmt.Fprintf(m.infoIo, "[ERROR] %s\n", msg)

	return err
}

func (m *Migrator) migrationsToApply(target string) ([]Migration, error) {
	before := []Migration{}

	for _, mig := range m.migrations {
		before = append(before, mig)
		if mig.Name() == target {
			return before, nil
		}
	}

	return nil, fmt.Errorf("didn't find migration: %s", target)
}

func (m *Migrator) migrationsToRevert(target string) ([]Migration, error) {
	after := []Migration{}

	for i := len(m.migrations) - 1; i >= 0; i-- {
		mig := m.migrations[i]
		after = append(after, mig)
		if mig.Name() == target {
			return after, nil
		}
	}

	return nil, fmt.Errorf("didn't find migration: %s", target)
}

func (m *Migrator) Up(to string) error {
	client, err := m.connector.GetClient()
	coll := client.Database(m.migrationsDB).Collection("applied_migrations")

	if err != nil {
		return err
	}

	migs := m.migrations

	if to != "" {
		migs, err = m.migrationsToApply(to)

		if err != nil {
			return err
		}
	}

	migCtx := &migrationContext{
		client:   client,
		logInfo:  m.logInfo,
		logError: m.logError,
	}

	for _, mig := range migs {
		res := coll.FindOne(context.TODO(), bson.D{{"name", mig.Name()}})

		if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
			// have to apply
			m.logInfo(fmt.Sprintf("applying: %s", mig.Name()))

			if err := mig.Up(migCtx); err != nil {
				return fmt.Errorf("Error while applying %s: %w", mig.Name(), err)
			}

			_, err := coll.InsertOne(context.TODO(), migrationRecord{
				Name: mig.Name(),
			})

			if err != nil {
				return err
			}
		} else {
			m.logInfo(fmt.Sprintf("skipping (already applied): %s", mig.Name()))
		}
	}

	return nil
}

func (m *Migrator) Down(to string) error {
	client, err := m.connector.GetClient()
	coll := client.Database(m.migrationsDB).Collection("applied_migrations")

	if err != nil {
		return err
	}

	migs := m.migrations

	if to != "" {
		migs, err = m.migrationsToRevert(to)

		if err != nil {
			return err
		}
	}

	migCtx := &migrationContext{
		client:   client,
		logInfo:  m.logInfo,
		logError: m.logError,
	}

	for _, mig := range migs {
		res := coll.FindOne(context.TODO(), bson.D{{"name", mig.Name()}})

		if res.Err() != nil && res.Err() == mongo.ErrNoDocuments {
			m.logInfo(fmt.Sprintf("skipping (not applied): %s", mig.Name()))
		} else {
			m.logInfo(fmt.Sprintf("reverting: %s", mig.Name()))

			if err := mig.Down(migCtx); err != nil {
				return fmt.Errorf("Error while reverting %s: %w", mig.Name(), err)
			}

			m.logInfo(fmt.Sprintf("reverted: %s", mig.Name()))

			_, err := coll.DeleteOne(context.TODO(), bson.D{{"name", mig.Name()}})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

type MigratorOpt func(*Migrator)

func WithErrorChannel(channel io.Writer) MigratorOpt {
	return func(m *Migrator) {
		m.errorIo = channel
	}
}

func WithInfoChannel(channel io.Writer) MigratorOpt {
	return func(m *Migrator) {
		m.infoIo = channel
	}
}

func NewMigrator(
	connector *Connector,
	migrationsDB string,
	migrations []Migration,
	opts ...MigratorOpt,
) *Migrator {
	m := &Migrator{
		connector:    connector,
		migrationsDB: migrationsDB,
		migrations:   migrations,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}
