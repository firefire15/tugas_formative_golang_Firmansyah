package database

import(
	"database/sql"
	"fmt"
	"os"
	"log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


var(
	DB *sql.DB
	err error
)

func ConnectDB(){
	_ = godotenv.Load()

	host     := os.Getenv("PGHOST")
	port     := os.Getenv("PGPORT")
	user     := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbname   := os.Getenv("PGDATABASE")
	
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Gagal membuka koneksi:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon (Ping Gagal):", err)
	}

	fmt.Println("Successfully connected to database!")
}

