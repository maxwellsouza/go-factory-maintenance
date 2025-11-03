package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/handlers"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository/postgres"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})

	ctx := context.Background()
	db, err := postgres.New(ctx)
	if err != nil {
		log.Fatalf("❌ failed to connect to database: %v", err)
	}
	defer db.Pool.Close()

	assetRepo := postgres.NewAssetRepo(db)
	workOrderRepo := postgres.NewWorkOrderRepo(db)

	assetService := service.NewAssetService(assetRepo)
	workOrderService := service.NewWorkOrderService(workOrderRepo)

	assetHandler := handlers.NewAssetHandler(assetService)
	workOrderHandler := handlers.NewWorkOrderHandler(workOrderService)

	assetHandler.RegisterRoutes(r)
	workOrderHandler.RegisterRoutes(r)

	log.Println("🚀 API (Postgres) running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
