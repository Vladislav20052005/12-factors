package postgre

import (
	"context"
	"fmt"
)

type Advertisement struct {
	AdId   string `json:"adId"`
	UserId string `json:"userId"`
	File   string `json:"file"`
	Url    string `json:"url"`
}

func GetAdsByUserId(userId string) ([]Advertisement, error) {
	rows, err := db.Query(`
	SELECT id, fileId, url FROM advertisements
	WHERE userId = $1
	ORDER BY createdAt DESC;`,
		userId)
	if err != nil {
		return nil, err
	}
	advs := make([]Advertisement, 0, 10)

	for rows.Next() {
		var adv Advertisement
		err := rows.Scan(&adv.AdId, &adv.File, &adv.Url)
		if err != nil {
			return nil, err
		}
		adv.UserId = userId
		advs = append(advs, adv)
	}
	fmt.Println(advs)

	return advs, nil
}

func GetAdById(adId, userId string) (*Advertisement, error) {
	fmt.Println(adId, userId)
	rows := db.QueryRow(`
	SELECT fileId, url FROM advertisements
	WHERE (id, userId) = ($1, $2);`,
		adId, userId)

	var adv Advertisement

	err := rows.Scan(&adv.File, &adv.Url)
	adv.AdId = adId
	adv.UserId = userId
	if err != nil {
		return nil, err
	}

	return &adv, nil
}

func InsertAd(adId, userId, fileId string, file []byte, url string) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	INSERT INTO files (id, data)
	VALUES ($1, $2);`,
		fileId, file)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO advertisements (id, userId, fileId, url)
	VALUES ($1, $2, $3, $4);`,
		adId, userId, fileId, url)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func EditAdURLOnly(adId, userId, url string) error {
	_, err := db.Exec(`
	UPDATE advertisements
	SET url = $1
	WHERE (id, userId) = ($2, $3);`,
		url, adId, userId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func EditAdNewFile(adId, userId string, file []byte, url string) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	WITH updated_ad AS (
		UPDATE advertisements
		SET url = $1
		WHERE (id, userId) = ($2, $3)
		RETURNING fileId
	)
	UPDATE files
	SET data = $4
	WHERE id = (SELECT fileId FROM updated_ad);`,
		url, adId, userId, file)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func RemoveAdById(adId, userId string) error {
	_, err := db.Exec(`
	DELETE FROM advertisements
	WHERE (id, userId) = ($1, $2);`,
		adId, userId)

	if err != nil {
		return err
	}

	return nil
}

func GetImage(fileId string) ([]byte, error) {
	rows := db.QueryRow(`
	SELECT data FROM files
	WHERE id = $1;`,
		fileId)

	data := make([]byte, 0, 1024)
	err := rows.Scan(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAdsByToken(token string) ([]Advertisement, error) {
	rows, err := db.Query(`
	SELECT advertisements.fileId, advertisements.url
	FROM advertisements
	JOIN users ON users.Id = advertisements.UserId
	WHERE users.Token = $1;`,
		token)

	if err != nil {
		return nil, err
	}

	advs := make([]Advertisement, 0, 10)

	for rows.Next() {
		var adv Advertisement
		err := rows.Scan(&adv.File, &adv.Url)
		if err != nil {
			return nil, err
		}
		advs = append(advs, adv)
	}

	return advs, nil
}
