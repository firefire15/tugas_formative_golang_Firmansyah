package main

import (
	"tugas-gin/routers"
	"tugas-gin/db"
	"os"
	"log"
	"github.com/joho/godotenv"
) 

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main(){
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 1. Ambil koneksi database dari connection.go
	database := db.ConnectDB() 
	defer database.Close()

	// 2. Jalankan fungsi migrasi yang baru saja kita ubah di database.go
	// Kita menangkap 1 return value berupa error
	err := db.RunMigrations(database) 
	if err != nil {
		log.Fatalf("Aplikasi gagal berjalan karena migrasi error: %v", err)
	}

	log.Printf("Server Bioskop berjalan di port :%s...\n", port)
	err = routers.StartAPIServer().Run(":" + port)
	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}