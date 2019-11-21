package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
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
	var c DBCredentials
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
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
