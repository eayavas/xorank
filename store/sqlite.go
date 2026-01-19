package store

import (
	"database/sql"
	"elo-app/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// Init tables here...
	return &Storage{DB: db}, nil
}

func (s *Storage) GetItems() ([]*models.Item, error) {
	// ... SELECT query logic
	return nil, nil // Placeholder
}

func (s *Storage) SaveVote(winnerID, loserID, username string) error {
	// ... INSERT/UPDATE logic
	return nil
}
