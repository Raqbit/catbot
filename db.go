package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type database struct {
	*sqlx.DB
}

func (d *database) BeginTransaction() (models.Transaction, error) {
	return d.Beginx()
}

func NewDb(dataSourceName string) (models.Datastore, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = migrateDatabase(db)

	if err != nil {
		logrus.WithError(err).Error("Could not migrate database")
		return nil, err
	}

	return &database{db}, nil
}

func migrateDatabase(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		return err
	}

	currVersion, dirty, err := getDatabaseVersion(m)

	if dirty {
		err = errors.New("database is dirty")
	}

	if err != nil {
		logrus.WithError(err).Error("Could not get database version")
		return err
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		logrus.WithField("version", currVersion).Info("Database up to date")
		return nil
	}

	if err != nil {
		return err
	}

	newVersion, dirty, err := getDatabaseVersion(m)

	if dirty {
		err = errors.New("database is dirty")
	}

	if err != nil {
		logrus.WithError(err).Error("Could not get database version")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"oldVersion": currVersion,
		"newVersion": newVersion,
	}).Info("Updated database")

	return nil
}

func getDatabaseVersion(m *migrate.Migrate) (string, bool, error) {
	version := "None"
	ver, dirty, err := m.Version()

	if err != nil && err != migrate.ErrNilVersion {
		return "", false, err
	}

	if err == nil {
		version = fmt.Sprint(ver)
	}

	return version, dirty, nil
}
