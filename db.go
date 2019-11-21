package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// Storage ...
type Storage struct {
	DB *sql.DB
}

// DBCredentials ...
type DBCredentials struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
}

// Connect ...
func (s *Storage) Connect(params DBCredentials) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		params.Host, params.Port, params.User, params.Password, params.DbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	s.DB = db
	return nil
}

// GetDBCredentials ...
func GetDBCredentials(fileName string) DBCredentials {
	c := DBCredentials{
		Host:     "0.0.0.0",
		Port:     5432,
		User:     "postgres",
		Password: "my-secret-password",
		DbName:   "golang",
	}

	_, ok := os.LookupEnv("DB_HOST")
	if ok {
		c.Host = os.Getenv("DB_HOST")
	}

	_, ok = os.LookupEnv("DB_PORT")
	if !ok {
		port, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			log.Println(err.Error())
		}
		c.Port = port
	}

	_, ok = os.LookupEnv("DB_USER")
	if !ok {
		c.User = os.Getenv("DB_USER")
	}

	_, ok = os.LookupEnv("DB_PASSWORD")
	if !ok {
		c.Password = os.Getenv("DB_PASSWORD")
	}

	_, ok = os.LookupEnv("DB_NAME")
	if !ok {
		c.DbName = os.Getenv("DB_NAME")
	}

	return c
}

// Init ...
func (s *Storage) Init() error {
	sql := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email character varying(255),
		password character varying(255)
	);`
	_, err := s.DB.Exec(sql)
	if err != nil {
		log.Printf("Error create PG table\n%s\n", err)
		return err
	}

	sql = `CREATE UNIQUE INDEX IF NOT EXISTS users_pkey ON users(id int4_ops);`
	_, err = s.DB.Exec(sql)
	if err != nil {
		log.Printf("Error create PG index\n%s\n", err)
		return err
	}

	sql = `CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users(email text_ops);`
	_, err = s.DB.Exec(sql)
	if err != nil {
		log.Printf("Error create PG index\n%s\n", err)
		return err
	}

	return nil
}
