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
		`CREATE TABLE IF NOT EXISTS users (passcode TEXT PRIMARY KEY);`,
		`CREATE TABLE IF NOT EXISTS items (id TEXT PRIMARY KEY, name TEXT, rating REAL, wins INT, losses INT);`,
		`CREATE TABLE IF NOT EXISTS votes (passcode TEXT, pair_key TEXT, PRIMARY KEY (passcode, pair_key));`,
	}

	for _, q := range queries {
		if _, err := s.DB.Exec(q); err != nil {
			log.Printf("Tablo hatası: %v", err)
		}
	}
}

func (s *Storage) seedData() {
	// sample codes
	passcodes := []string{
		"1234",  // Test
		"admin", // Admin
		"7777",  // Lucky
		"0000",  // Guest
	}
	for _, p := range passcodes {
		s.DB.Exec("INSERT OR IGNORE INTO users (passcode) VALUES (?)", p)
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
	log.Println("DB init complete: Single Passcode Auth Mode.")
}

// --- METHODS ---

// CheckPasscode checks if the passcode is valid
func (s *Storage) CheckPasscode(code string) bool {
	var exists int
	err := s.DB.QueryRow("SELECT 1 FROM users WHERE passcode = ?", code).Scan(&exists)
	return err == nil && exists == 1
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

func (s *Storage) HasVoted(passcode, pairKey string) bool {
	var count int
	s.DB.QueryRow("SELECT COUNT(*) FROM votes WHERE passcode = ? AND pair_key = ?", passcode, pairKey).Scan(&count)
	return count > 0
}

func (s *Storage) SaveVote(passcode, pairKey string, winner, loser *models.Item) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec("INSERT INTO votes (passcode, pair_key) VALUES (?, ?)", passcode, pairKey); err != nil {
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
