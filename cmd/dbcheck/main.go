package main

import (
	"context"
	"log"

	pg "github.com/maxwellsouza/go-factory-maintenance/internal/repository/postgres"
)

func main() {
	ctx := context.Background()
	db, err := pg.New(ctx)
	if err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}
	defer db.Pool.Close()

	log.Println("✅ DB connection OK (pgxpool)")
}
