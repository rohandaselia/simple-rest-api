package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

// 1. Create - POST /bioskop
func (d *DBConfig) addBioskop(c *gin.Context) {
	var b Bioskop
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Validasi dasar
	if b.Nama == "" || b.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan lokasi tidak boleh kosong"})
		return
	}

	query := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3)`
	_, err := d.DB.Exec(query, b.Nama, b.Lokasi, b.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bioskop berhasil ditambahkan"})
}

// 2. Read All - GET /bioskop
func (d *DBConfig) getAllBioskop(c *gin.Context) {
	rows, err := d.DB.Query("SELECT id, nama, lokasi, rating FROM bioskop")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari database"})
		return
	}
	defer rows.Close()

	bioskops := []Bioskop{} // Inisialisasi slice kosong agar tidak return null
	for rows.Next() {
		var b Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data bioskop"})
			return
		}
		bioskops = append(bioskops, b)
	}

	// Cek jika ada error selama proses iterasi
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memproses data"})
		return
	}

	c.JSON(http.StatusOK, bioskops)
}

// 3. Read By ID - GET /bioskop/:id
func (d *DBConfig) getBioskopByID(c *gin.Context) {
	id := c.Param("id")

	var b Bioskop
	query := "SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1"
	err := d.DB.QueryRow(query, id).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari database"})
		}
		return
	}

	c.JSON(http.StatusOK, b)
}

// 4. Update - PUT /bioskop/:id
func (d *DBConfig) updateBioskop(c *gin.Context) {
	id := c.Param("id")

	var b Bioskop
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Validasi input
	if b.Nama == "" || b.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan lokasi tidak boleh kosong"})
		return
	}

	query := "UPDATE bioskop SET nama = $1, lokasi = $2, rating = $3 WHERE id = $4"
	result, err := d.DB.Exec(query, b.Nama, b.Lokasi, b.Rating, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data di database"})
		return
	}

	// Cek apakah ada data yang benar-benar diperbarui
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil diperbarui"})
}

// 5. Delete - DELETE /bioskop/:id
func (d *DBConfig) deleteBioskop(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM bioskop WHERE id = $1"
	result, err := d.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data dari database"})
		return
	}

	// Cek apakah ada data yang terhapus (apakah ID valid)
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil dihapus"})
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

	// Endpoints
	router.POST("/bioskop", appDB.addBioskop)           // Create
	router.GET("/bioskop", appDB.getAllBioskop)         // Read All
	router.GET("/bioskop/:id", appDB.getBioskopByID)    // Read by ID
	router.PUT("/bioskop/:id", appDB.updateBioskop)     // Update
	router.DELETE("/bioskop/:id", appDB.deleteBioskop)  // Delete

	// Jalankan Server
	fmt.Println("Server berjalan di port 8080...")
	router.Run(":8080")
}