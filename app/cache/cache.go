package cache

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nalgeon/redka"
)

var cacheInstance *redka.DB

// Get returns the instantiated DB instance.
func GetCache() *redka.DB {
	return cacheInstance
}

func init() {
	cacheDbFile := os.Getenv("CACHE_DB_NAME")
	// Open or create the data.db file.
	db, err := redka.Open(cacheDbFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Always close the database when you are finished.
	// defer db.Close()
	cacheInstance = db
	// Set some string keys.
	// err = db.Str().Set("name", "alice")
	// slog.Info("set", "err", err)
	// err = db.Str().Set("age", 25)
	// slog.Info("set", "err", err)

	// // Check if the keys exist.
	// count, err := db.Key().Count("name", "age", "city")
	// slog.Info("count", "count", count, "err", err)

	// // Get a key.
	// name, err := db.Str().Get("name")
	// slog.Info("get", "name", name, "err", err)
}
