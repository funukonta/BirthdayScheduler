package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Storage interface {
	GetBirthdayData() ([]User, error)
	GeneratePromoCode([]User) error
}

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
	err = p.generateDummy()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) createTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id_user INT PRIMARY KEY,
		username VARCHAR(20) UNIQUE,
		name VARCHAR(50),
		email VARCHAR(25),
		birthDate TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS promo (
		id_promo INT PRIMARY KEY,
		promoName VARCHAR(50),
		promoDesc text
	);
	
	CREATE TABLE IF NOT EXISTS special_promo (
		promoCode VARCHAR(20) PRIMARY KEY,
		validUntil TIMESTAMP,
		id_user INT REFERENCES users(id_user),
		id_promo INT REFERENCES promo(id_promo)
	);
	`
	_, err := p.db.Exec(query)
	return err
}

func (p *Postgres) generateDummy() error {
	_, err := p.db.Exec(`INSERT user values (1,'evan','evan','evan@gmail.com','2023-01-19 15:05:04'),
	(2,'roy','roy','roy@gmail.com','2023-01-19 15:05:04'),
	(3,'darmawan','darmawan','darmawan@gmail.com','2023-01-20 15:05:04'),`)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(`INSERT promo values (1,'BirthDay','birthday special promo'),
	(2,'B1G1','Promo buy 1 get 1')`)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetBirthdayData() ([]User, error) {
	now := time.Now()
	rows, err := p.db.Query(`select * from user where birthday=$1`, now.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	result := []User{}
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u)
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	return result, nil
}

func (p *Postgres) GeneratePromoCode(users []User) error {
	for _, user := range users {
		_, err := p.db.Exec(`Insert special_promo values ($1,$2,$3,$4,$5)`, user.Promo.PromoCode, user.Promo.ValidUntil, user.Id_user, user.Promo.Id_Promo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Postgres) GetPromo() (int, error) {
	var id_promo int
	err := p.db.QueryRow(`select id_promo where promoname='BirthDay'`).Scan(&id_promo)
	if err != nil {
		return 0, err
	}

	return id_promo, nil
}
