package models

import "time"

type Cat struct {
	BaseModel
	OwnerId       uint      `db:"owner_id"`
	Name          string    `db:"name"`
	CryptoKittyID int       `db:"ck_id"`
	Pronoun       string    `db:"pronoun"`
	Hunger        int       `db:"hunger"`
	LastFed       time.Time `db:"last_fed"`
	Away          bool      `db:"away"`
	ReturnTime    time.Time `db:"return_time"`
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
			"ck_id": cryptoKittyId,
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

	err := db.Get(&cat, `select * from cats where owner_id = $1 and name = $2`, ownerId, name)

	if err != nil {
		return nil, err
	}

	return &cat, nil
}