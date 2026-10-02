package postgre

import "fmt"

func GetUserIdByLogIn(username, hashedPassword string) (string, error) {
	rows := db.QueryRow(`
	SELECT id FROM users
	WHERE (username, password) = ($1, $2);`,
		username, hashedPassword)

	var id string
	err := rows.Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func InsertUser(userId, username, hashedPassword, token string) error {
	fmt.Println(userId, username, hashedPassword, token)
	_, err := db.Exec(`
	INSERT INTO users (id, username, password, token)
	VALUES ($1, $2, $3, $4);`,
		userId, username, hashedPassword, token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func RemoveByUserId(userId string) error {
	_, err := db.Exec(`
	DELETE FROM users
	WHERE id = $1;`,
		userId)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetUserByUserId(userId string) (string, string, error) {
	rows := db.QueryRow(`
	SELECT username, token FROM users
	WHERE id = $1;`,
		userId)

	var username string
	var token string
	err := rows.Scan(&username, &token)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	return username, token, nil
}
