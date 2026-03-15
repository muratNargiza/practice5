package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "postgres")
	password := getenv("DB_PASSWORD", "postgres")
	dbname := getenv("DB_NAME", "practice5")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			name       VARCHAR(100) NOT NULL,
			email      VARCHAR(100) UNIQUE NOT NULL,
			gender     VARCHAR(10)  NOT NULL,
			birth_date DATE         NOT NULL
		);

		CREATE TABLE IF NOT EXISTS user_friends (
			user_id   INTEGER REFERENCES users(id) ON DELETE CASCADE,
			friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			PRIMARY KEY (user_id, friend_id),
			CHECK (user_id <> friend_id)
		);
	`)
	return err
}
