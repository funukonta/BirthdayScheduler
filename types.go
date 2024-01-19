package main

import "time"

type User struct {
	Id_user   int       `db:"id_user"`
	Username  string    `db:"username"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	BirthDate time.Time `db:"birthDate"`
	Promo     SpecialPromo
}

type SpecialPromo struct {
	PromoCode  string
	ValidUntil time.Time
	Id_Promo   int
}
