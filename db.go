package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type BaseModel struct {
	ID        uint       `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Datastore interface {
	// User functions
	GetUserById(userId uint) (*User, error)
	GetUserByDiscordId(discordId string) (*User, error)
	GetUserOrCreate(discordId string) (*User, error)
	UserModifyMoney(userId uint, amount int64) error
	UserUseDaily(userId uint, amount int64) (int64, error)

	// Cat functions
	AllCatsOfUser(ownerId uint) ([]*Cat, error)
	GetCatByName(ownerId uint, name string) (*Cat, error)
	CreateCatForUser(ownerId uint, cryptoKittyId int, name string, pronoun string) error
	CatNameExists(ownerId uint, name string) (bool, error)
	MarkCatAwayUntil(catId uint, channelId string, until time.Time) error
	UpdateReturningCats() ([]*Cat, error)
}

type DB struct {
	*sqlx.DB
}

type Tx struct {
	*sql.Tx
}

func NewDb(dataSourceName string) (*DB, error) {
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

	return &DB{db}, nil
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

func (db *DB) BeginTransaction() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
