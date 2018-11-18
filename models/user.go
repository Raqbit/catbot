package models

import (
	"time"
)

type User struct {
	BaseModel
	DiscordId string    `db:"discord_id"`
	Money     int64     `db:"money"`
	LastDaily time.Time `db:"last_daily"`
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
	_, err := db.NamedExec(
		`update users set 
			money = money + :amount,
			last_daily = now()
			where id = :user_id`,
		map[string]interface{}{
			"amount":  amount,
			"user_id": userId,
		},
	)

	if err != nil {
		return 0, err
	}

	var newMoneyVal int64

	err = db.Get(&newMoneyVal, `select money from users where id = $1`, userId)

	if err != nil {
		return 0, err
	}

	return newMoneyVal, nil
}

func (db *DB) UserRemoveMoney(userId uint, amount int64) error {
	_, err := db.NamedExec(
		`update users set 
			money = money - :amount,
			last_daily = now()
			where id = :user_id`,
		map[string]interface{}{
			"amount":  amount,
			"user_id": userId,
		},
	)

	if err != nil {
		return err
	}

	return nil
}