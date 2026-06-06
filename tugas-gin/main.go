package main

import (
	"tugas-gin/routers"
	"os"
	"github.com/joho/godotenv"
) 


func main(){
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	routers.StartBioskopServer().Run(":"+port)
}