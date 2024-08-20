package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// By default this is a pre-configured Gorm DB instance.
// Change this type based on the database package of your likings.
var dbInstance *sql.DB

// Get returns the instantiated DB instance.
func GetDb() *sql.DB {
	return dbInstance
}

func init() {
	driver := os.Getenv("DB_DRIVER")
	appDbFile := os.Getenv("APP_DB_NAME")
	// Create a default *sql.DB exposed by the Playground/db package
	// based on the given configuration.

	db, err := sql.Open(driver, appDbFile)

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// defer db.Close()

	dbInstance = db
}
