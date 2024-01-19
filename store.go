package main

import (
	"database/sql"
	"fmt"
	"os"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres() (*Postgres, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("USER")
	dbname := os.Getenv("DBNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := os.Getenv("SSLMODE")
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s", host, user, dbname, password, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

// Migration table and creation of dummy data
func (p *Postgres) Init() error {
	err := p.createTable()
	if err != nil {
		return err
	}
	// err = p.generateDummy()
	return nil
}

func (p *Postgres) createTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id_user INT PRIMARY KEY,
		username VARCHAR(20) UNIQUE,
		name VARCHAR(50),
		password VARCHAR(60),
		birthDate TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS promo (
		id_promo INT PRIMARY KEY,
		promoName VARCHAR(50),
		validFrom TIMESTAMP,
		validUntil TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS special_promo (
		promoCode VARCHAR(20) PRIMARY KEY,
		id_user INT REFERENCES users(id_user),
		id_promo INT REFERENCES promo(id_promo)
	);
	
	`
	_, err := p.db.Exec(query)
	return err
}
