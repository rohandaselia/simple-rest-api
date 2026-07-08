package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type DBConfig struct {
	DB *sql.DB
}

type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama"`
	Lokasi string  `json:"lokasi"`
	Rating float64 `json:"rating"`
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}

func (d *DBConfig) addBioskop(c *gin.Context) {
	var b Bioskop
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nama dan lokasi tidak boleh kosong"})
		return
	}

	query := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3)`
	_, err := d.DB.Exec(query, b.Nama, b.Lokasi, b.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil ditambahkan"})
}

func main() {
	// Setup Konfigurasi
	dbConf := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password", // Sesuaikan dengan password database Anda
		DBName:   "bioskop_db",
	}

	// Membuka koneksi database
	db, err := sql.Open("postgres", dbConf.ConnectionString())
	if err != nil {
		log.Fatal("Gagal membuka koneksi database:", err)
	}
	defer db.Close()

	// Inisialisasi struct DBConfig
	appDB := &DBConfig{DB: db}

	// Setup Router Gin
	router := gin.Default()

	// Endpoint
	router.POST("/bioskop", appDB.addBioskop)

	// Jalankan Server
	fmt.Println("Server berjalan di port 8080...")
	router.Run(":8080")
}
