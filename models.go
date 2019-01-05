package main

import (
	"github.com/pkg/errors"
	"time"
)

type Cat struct {
	BaseModel
	OwnerId       uint      `db:"owner_id"`
	Name          string    `db:"name"`
	CryptoKittyID int       `db:"ck_id"`
	Pronoun       string    `db:"pronoun"`
	Hunger        int       `db:"hunger"`
	LastFed       time.Time `db:"last_fed"`
	Away          bool      `db:"away"`
	AwayUntil     time.Time `db:"away_until"`
	AwayChannel   string    `db:"away_channel"`
}

func (db *DB) AllCatsOfUser(ownerId uint) ([]*Cat, error) {
	var cats []*Cat

	err := db.Select(&cats, `select * from cats where owner_id=$1`, ownerId)

	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (db *DB) CreateCatForUser(ownerId uint, cryptoKittyId int, name string, pronoun string) error {
	_, err := db.NamedExec(
		`insert into cats (owner_id, ck_id, name, pronoun)
			   values (:owner_id, :ck_id, :name, :pronoun)`,
		map[string]interface{}{
			"owner_id": ownerId,
			"ck_id":    cryptoKittyId,
			"name":     name,
			"pronoun":  pronoun,
		})

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CatNameExists(ownerId uint, name string) (bool, error) {
	var exists bool

	err := db.Get(&exists, `select exists(select 1 from cats where owner_id = $1 and name = $2)`, ownerId, name)

	if err != nil {
		return true, err
	}

	return exists, nil
}

func (db *DB) GetCatByName(ownerId uint, name string) (*Cat, error) {
	var cat Cat

	rows, err := db.Queryx(`select * from cats where owner_id = $1 and name = $2`, ownerId, name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() == false {
		return nil, nil
	}

	err = rows.StructScan(&cat)

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (db *DB) MarkCatAwayUntil(catId uint, channelId string, until time.Time) error {
	_, err := db.NamedExec(
		`update cats set 
			away = true,
            away_until = :away_until,
            away_channel = :away_channel
			where id = :cat_id`,
		map[string]interface{}{
			"cat_id":       catId,
			"away_channel": channelId,
			"away_until":   until,
		},
	)

	return err
}

func (db *DB) UpdateReturningCats() ([]*Cat, error) {
	var cats []*Cat

	err := db.Select(&cats, `update cats set away = false where away = true and away_until <= now() returning *`)

	if err != nil {
		return nil, err
	}

	return cats, nil
}

type User struct {
	BaseModel
	DiscordId string    `db:"discord_id"`
	Money     int64     `db:"money"`
	LastDaily time.Time `db:"last_daily"`
}

func (db *DB) GetUserById(userId uint) (*User, error) {
	var user User

	err := db.Get(&user, `select * from users where id=$1 limit 1`, userId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserByDiscordId(discordId string) (*User, error) {
	var user User

	err := db.Get(&user, `select * from users where discord_id=$1 limit 1`, discordId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) GetUserOrCreate(discordId string) (*User, error) {
	user, err := db.GetUserByDiscordId(discordId)

	if err == nil && user != nil {
		return user, err
	}

	_, err = db.NamedExec(`insert into users (discord_id) values (:discord_id)`,
		map[string]interface{}{
			"discord_id": discordId,
		},
	)

	if err != nil {
		return nil, err
	}

	user, err = db.GetUserByDiscordId(discordId)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) UserUseDaily(userId uint, amount int64) (int64, error) {
	rows, err := db.NamedQuery(
		`update users set 
			money = money + :amount,
			last_daily = now()
			where id = :user_id
			returning money`,
		map[string]interface{}{
			"amount":  amount,
			"user_id": userId,
		},
	)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows.Next() == false {
		return 0, errors.New("No rows returned")
	}

	var newMoneyVal int64

	err = rows.Scan(&newMoneyVal)

	return newMoneyVal, nil
}

func (db *DB) UserModifyMoney(userId uint, amount int64) error {
	_, err := db.NamedExec(
		`update users set 
			money = money + :amount
			where id = :user_id`,
		map[string]interface{}{
			"amount":  amount,
			"user_id": userId,
		},
	)

	return err
}
