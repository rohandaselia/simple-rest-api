package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"simple-rest-api/models"
)

type BioskopHandler struct {
	DB *sql.DB
}

func NewBioskopHandler(db *sql.DB) *BioskopHandler {
	return &BioskopHandler{DB: db}
}

func (h *BioskopHandler) AddBioskop(c *gin.Context) {
	var b models.Bioskop
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if b.Nama == "" && b.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
		return
	}

	query := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3)`
	_, err := h.DB.Exec(query, b.Nama, b.Lokasi, b.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil ditambahkan"})
}

func (h *BioskopHandler) GetAll(c *gin.Context) {
	rows, err := h.DB.Query("SELECT id, nama, lokasi, rating FROM bioskop")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data: " + err.Error()})
		return
	}
	defer rows.Close()

	bioskops := []models.Bioskop{}
	for rows.Next() {
		var b models.Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data: " + err.Error()})
			return
		}
		bioskops = append(bioskops, b)
	}

	c.JSON(http.StatusOK, bioskops)
}

func (h *BioskopHandler) GetByID(c *gin.Context) {
	var b models.Bioskop
	err := h.DB.QueryRow("SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1", c.Param("id")).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BioskopHandler) UpdateBioskop(c *gin.Context) {
	var b models.Bioskop
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if b.Nama == "" && b.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
		return
	}

	res, err := h.DB.Exec("UPDATE bioskop SET nama = $1, lokasi = $2, rating = $3 WHERE id = $4", b.Nama, b.Lokasi, b.Rating, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data: " + err.Error()})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil diperbarui"})
}

func (h *BioskopHandler) DeleteBioskop(c *gin.Context) {
	res, err := h.DB.Exec("DELETE FROM bioskop WHERE id = $1", c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data: " + err.Error()})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil dihapus"})
}
