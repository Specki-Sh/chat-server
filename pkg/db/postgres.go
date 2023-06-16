package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var database *sql.DB

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func initDB(config Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

// StartDbConnection Creates connection to database
func StartDbConnection(config Config) {
	database = initDB(config)
}

// GetDBConn func for getting db conn globally
func GetDBConn() *sql.DB {
	return database
}

func CloseDbConnection() error {
	if err := database.Close(); err != nil {
		return fmt.Errorf("error occurred on database connection closing: %s", err.Error())
	}
	return nil
}
