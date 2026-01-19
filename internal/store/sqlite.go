package store

import (
	"database/sql"
	"log"
	"xorank/internal/models"

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

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Storage{DB: db}
	s.initSchema()
	s.seedData()

	return s, nil
}

func (s *Storage) initSchema() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (username TEXT PRIMARY KEY, password TEXT);`,
		`CREATE TABLE IF NOT EXISTS items (id TEXT PRIMARY KEY, name TEXT, rating REAL, wins INT, losses INT);`,
		`CREATE TABLE IF NOT EXISTS votes (username TEXT, pair_key TEXT, PRIMARY KEY (username, pair_key));`,
	}

	for _, q := range queries {
		if _, err := s.DB.Exec(q); err != nil {
			log.Printf("Table creation error: %v", err)
		}
	}
}

func (s *Storage) seedData() {
	// Users
	users := map[string]string{
		"efe":   "1234",
		"guest": "0000",
		"admin": "root",
	}
	for u, p := range users {
		s.DB.Exec("INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)", u, p)
	}

	// Items
	var count int
	s.DB.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if count > 0 {
		return
	}

	items := []models.Item{
		{ID: "1", Name: "Firefox", Rating: 1200},
		{ID: "2", Name: "Chrome", Rating: 1200},
		{ID: "3", Name: "Brave", Rating: 1200},
		{ID: "4", Name: "Edge", Rating: 1200},
		{ID: "5", Name: "Safari", Rating: 1200},
		{ID: "6", Name: "Opera", Rating: 1200},
		{ID: "7", Name: "Vivaldi", Rating: 1200},
		{ID: "8", Name: "Arc", Rating: 1200},
	}

	for _, item := range items {
		s.DB.Exec("INSERT INTO items (id, name, rating, wins, losses) VALUES (?, ?, ?, ?, ?)",
			item.ID, item.Name, item.Rating, 0, 0)
	}
	log.Println("Database is ready with seed data.")
}

// --- METHODS ---

func (s *Storage) GetUser(username, password string) bool {
	var storedPass string
	err := s.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPass)
	if err != nil {
		return false
	}
	return storedPass == password
}

func (s *Storage) GetAllItems() ([]*models.Item, error) {
	rows, err := s.DB.Query("SELECT id, name, rating, wins, losses FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		i := &models.Item{}
		if err := rows.Scan(&i.ID, &i.Name, &i.Rating, &i.Wins, &i.Losses); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (s *Storage) HasVoted(username, pairKey string) bool {
	var count int
	s.DB.QueryRow("SELECT COUNT(*) FROM votes WHERE username = ? AND pair_key = ?", username, pairKey).Scan(&count)
	return count > 0
}

func (s *Storage) SaveVote(username, pairKey string, winner, loser *models.Item) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec("INSERT INTO votes (username, pair_key) VALUES (?, ?)", username, pairKey); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("UPDATE items SET rating = ?, wins = wins + 1 WHERE id = ?", winner.Rating, winner.ID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("UPDATE items SET rating = ?, losses = losses + 1 WHERE id = ?", loser.Rating, loser.ID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
