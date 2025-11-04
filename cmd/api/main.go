package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/handlers"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/middleware"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository/postgres"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
	logrus "github.com/sirupsen/logrus"
)

func main() {
	// Configura logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.LoggerMiddleware())

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

	logrus.Info("🚀 API (Postgres) running on :8080")
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("server error: %v", err)
	}
}
