package main

import (
	"fmt"
	"log"
	"os"

	"simple-rest-api/config"
	"simple-rest-api/database"
	"simple-rest-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	dbConf := config.Config{
		Host:     os.Getenv("PGHOST"),
		Port:     os.Getenv("PGPORT"),
		User:     os.Getenv("PGUSER"),
		Password: os.Getenv("PGPASSWORD"), // Sesuaikan dengan password database Anda
		DBName:   os.Getenv("PGDATABASE"),
	}

	db := database.ConnectDB(dbConf)
	defer db.Close()

	router := gin.Default()

	bioskopHandler := handlers.NewBioskopHandler(db)
	router.POST("/bioskop", bioskopHandler.AddBioskop)
	router.GET("/bioskop", bioskopHandler.GetAll)
	router.GET("/bioskop/:id", bioskopHandler.GetByID)
	router.PUT("/bioskop/:id", bioskopHandler.UpdateBioskop)
	router.DELETE("/bioskop/:id", bioskopHandler.DeleteBioskop)

	fmt.Println("Server berjalan di port 3000...")
	if err := router.Run(":3000"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
