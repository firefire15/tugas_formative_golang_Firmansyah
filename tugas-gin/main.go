package main

import (
	"tugas-gin/routers"
	"os"
	"log"
	"github.com/joho/godotenv"
) 


func main(){
	// Load .env file jika ada (opsional untuk production)
	godotenv.Load()
	
	// Pastikan PORT tersedia, default ke 8080 jika tidak ada
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	routers.StartBioskopServer().Run(":" + port)
}

