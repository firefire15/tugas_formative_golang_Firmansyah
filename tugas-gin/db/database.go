package db

import (
	"database/sql"
	"embed"
	"log"
	"net/http" // <-- 1. WAJIB TAMBAHKAN IMPORT INI

	migrate "github.com/rubenv/sql-migrate"
)


var migrationFiles embed.FS

func RunMigrations(db *sql.DB) error {
	log.Println("Memeriksa migrasi struktur tabel...")

	migrations := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(migrationFiles),
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	if n > 0 {
		log.Printf("Sukses! %d file migrasi baru berhasil diterapkan.\n", n)
	} else {
		log.Println("Struktur database sudah sesuai versi terbaru. Tidak ada perubahan.")
	}

	return nil
}