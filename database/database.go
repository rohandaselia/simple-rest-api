package database

import (
	"database/sql"
	"log"

	"simple-rest-api/config"
	"simple-rest-api/database/migration"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	DB *sql.DB
}

func ConnectDB(cfg config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		log.Fatal("Gagal membuka koneksi database:", err)
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	DBMigrate(db)

	return db
}
