package models

import (
	"errors"
	"time"
)

type User struct {
	BaseModel
	DiscordId string    `db:"discord_id"`
	Money     int64     `db:"money"`
	LastDaily time.Time `db:"last_daily"`
}

func (u *User) UseDaily(db Queryable, amount int64) (int64, error) {
	rows, err := db.NamedQuery(
		`update users set 
			money = money + :amount,
			last_daily = now()
			where id = :user_id
			returning money`,
		map[string]interface{}{
			"amount":  amount,
			"user_id": u.ID,
		},
	)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows.Next() == false {
		return 0, errors.New("no rows returned")
	}

	var newMoneyVal int64

	err = rows.Scan(&newMoneyVal)

	return newMoneyVal, nil
}

func (u *User) ModifyMoney(db Queryable, amount int64) error {
	_, err := db.NamedExec(
		`update users set 
			money = money + :amount
			where id = :user_id`,
		map[string]interface{}{
			"amount":  amount,
			"user_id": u.ID,
		},
	)

	return err
}

type UserStore struct{}

func (us *UserStore) GetById(db Queryable, userId uint) (*User, error) {
	var user User

	err := db.Get(&user, `select * from users where id=$1 limit 1`, userId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserStore) GetByDiscordId(db Queryable, discordId string) (*User, error) {
	var user User

	err := db.Get(&user, `select * from users where discord_id=$1 limit 1`, discordId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserStore) GetOrCreate(db Queryable, discordId string) (*User, error) {
	user, err := us.GetByDiscordId(db, discordId)

	if err == nil && user != nil {
		return user, err
	}

	if _, err = db.NamedExec(`insert into users (discord_id) values (:discord_id)`,
		map[string]interface{}{
			"discord_id": discordId,
		},
	); err != nil {
		return nil, err
	}

	return us.GetByDiscordId(db, discordId)
}
