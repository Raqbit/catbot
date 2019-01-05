package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
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
	return &DB{db}, nil
}

func (db *DB) BeginTransaction() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}
