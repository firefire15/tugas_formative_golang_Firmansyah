package main

import (
	"tugas-gin/routers"
	"os"
	"log"
	"github.com/joho/godotenv"
) 


func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal memuat file .env")
	}

	routers.StartBioskopServer().Run(":"+os.Getenv("PORT"))
}