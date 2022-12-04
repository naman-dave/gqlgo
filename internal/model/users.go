package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/naman-dave/gqlgo/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64  `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Mobilenumber string `json:"mobilenumber"`
	Passkey      string `json:"passkey"`
	Token        string `json:"token"`
}

func (user *User) Create() error {

	userID, _, _ := GetUserIdByUsername(user.Mobilenumber)
	if userID != 0 {
		return fmt.Errorf("mobilenumber already exists, try login")
	}

	statement, err := db.Db.Prepare("INSERT INTO Users(Mobilenumber,Passkey) VALUES($1,$2)")
	if err != nil {
		log.Println(err)
		return err
	}
	defer statement.Close()

	hashedPassword, err := HashPassword(user.Passkey)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = statement.Exec(user.Mobilenumber, hashedPassword)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (user *User) Authenticate() (bool, error) {
	statement, err := db.Db.Prepare("select Passkey from Users WHERE Mobilenumber = $1")
	if err != nil {
		log.Fatal(err)
	}

	row := statement.QueryRow(user.Mobilenumber)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	return CheckPasswordHash(user.Passkey, hashedPassword), nil
}

func (u *User) UpdateToken() error {

	updateSQL := ` UPDATE USERS
	SET Token = $1
	WHERE Mobilenumber = $2 `

	_, err := db.Db.Exec(updateSQL, u.Token, u.Mobilenumber)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func RemoveToken(userID int) error {

	updateSQL := ` UPDATE USERS
	SET Token = $1
	WHERE id = $2 `

	_, err := db.Db.Exec(updateSQL, "", userID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserIdByUsername(mobilenumber string) (int, string, error) {
	statement, err := db.Db.Prepare("select ID, token from Users WHERE Mobilenumber = $1 AND (token != '') ")
	if err != nil {
		log.Println(err)
		return 0, "", err
	}
	row := statement.QueryRow(mobilenumber)

	var Id int
	var token string

	err = row.Scan(&Id, &token)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			return 0, "", err
		}
		return 0, "", err
	}

	return Id, token, nil
}
