package models

import (
	"database/sql"
	"time"
)

type Cat struct {
	BaseModel
	OwnerId       uint      `db:"owner_id"`
	Name          string    `db:"name"`
	CryptoKittyID int       `db:"ck_id"`
	Hunger        int       `db:"hunger"`
	LastFed       time.Time `db:"last_fed"`
	Away          bool      `db:"away"`
	AwayUntil     time.Time `db:"away_until"`
	AwayChannel   string    `db:"away_channel"`
}

func (c *Cat) MarkAwayUntil(db Queryable, channelId string, until time.Time) error {
	_, err := db.NamedExec(
		`update cats set 
			away = true,
            away_until = :away_until,
            away_channel = :away_channel
			where id = :cat_id`,
		map[string]interface{}{
			"cat_id":       c.ID,
			"away_channel": channelId,
			"away_until":   until,
		},
	)

	return err
}

type CatStore struct{}

func (cs *CatStore) GetAllCatsOfUser(db Queryable, owner *User) ([]*Cat, error) {
	var cats []*Cat

	err := db.Select(
		&cats,
		`select * from cats where owner_id = $1`,
		owner.ID,
	)

	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (cs *CatStore) GetByName(db Queryable, owner *User, name string) (*Cat, error) {
	var cat Cat

	err := db.Get(
		&cat,
		`select * from cats where owner_id = $1 and name = $2`,
		owner.ID,
		name,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (cs *CatStore) CreateForUser(db Queryable, owner *User, cryptoKittyId int, name string) error {
	_, err := db.NamedExec(
		`insert into cats (owner_id, ck_id, name)
			   values (:owner_id, :ck_id, :name)`,
		map[string]interface{}{
			"owner_id": owner.ID,
			"ck_id":    cryptoKittyId,
			"name":     name,
		})

	if err != nil {
		return err
	}

	return nil
}

func (cs *CatStore) CatNameExists(db Queryable, owner *User, name string) (bool, error) {
	var exists bool

	err := db.Get(
		&exists,
		`select exists(select 1 from cats where owner_id = $1 and name = $2)`,
		owner.ID,
		name,
	)

	if err != nil {
		return true, err
	}

	return exists, nil
}

func (cs *CatStore) UpdateReturning(db Queryable) ([]*Cat, error) {
	var cats []*Cat

	err := db.Select(
		&cats,
		`update cats set away = false where away = true and away_until <= now() returning *`,
	)

	if err != nil {
		return nil, err
	}

	return cats, nil
}
