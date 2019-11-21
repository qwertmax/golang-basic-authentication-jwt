package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
)

// User ...
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateHash ...
func (u User) CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// IsHashEqual ...
func (u User) IsHashEqual(hash string) bool {
	if u.Password == hash {
		return true
	}
	return false
}

// Save ...
func (u *User) Save(db *sql.DB) error {
	sql := fmt.Sprintf("INSERT INTO users (email, password) VALUES('%s', '%s');", u.Email, u.CreateHash(u.Password))
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// Get ...
func (u *User) Get(db *sql.DB, email string) error {
	sql := fmt.Sprintf("SELECT email, password FROM users WHERE email = '%s';", email)
	row := db.QueryRow(sql)
	err := row.Scan(&u.Email, &u.Password)
	if err != nil {
		return err
	}
	return nil
}
