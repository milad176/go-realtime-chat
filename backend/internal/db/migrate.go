package db

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(db *pgxpool.Pool) {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		log.Printf("Running migration: %s\n", file)

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		_, err = db.Exec(context.Background(), string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}
	}

	log.Println("All migrations ran successfully")
}
