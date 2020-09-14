package db

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

type Config struct {
	Addr     string
	User     string
	Database string
	CertPath string
}

func TestOpen(t *testing.T) {
	config := &Config{
		Addr:     "127.0.0.1:26257",
		User:     "root",
		Database: "rksbook",
		CertPath: "",
	}

	db, err := Open("postgres", connString(config))
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	query := `INSERT INTO notification (text) VALUES ('qwerty')`

	if _, err := tx.Exec(query); err != nil {
		panic(err)
	}

	tx.Commit()
}

func connString(config *Config) string {
	switch {
	case config.CertPath != "":
		ssl := fmt.Sprintf("&sslmode=%s&sslcert=%s/client.%s.crt&sslkey=%s/client.%s.key&sslrootcert=%s/ca.crt",
			"verify-full", config.CertPath, config.User, config.CertPath, config.User, config.CertPath)

		return fmt.Sprintf("postgres://%s@%s/%s?%s", config.User, config.Addr, config.Database, ssl)
	default:
		return fmt.Sprintf("postgresql://%s@%s/%s?sslmode=disable", config.User, config.Addr, config.Database)
	}
}
