package models

type Item struct {
	ID     string
	Name   string
	Rating float64
	Wins   int
	Losses int
}

type User struct {
	Username string
	Password string
}
