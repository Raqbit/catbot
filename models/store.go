package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Queryable interface {
	Select(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
}

type Transaction interface {
	Queryable
	Commit() error
	Rollback() error
}

type Datastore interface {
	Queryable
	BeginTransaction() (Transaction, error)
}

type BaseModel struct {
	ID        uint       `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

var Users UserStore
var Cats CatStore
