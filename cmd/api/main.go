package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/handlers"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository/memory"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Healthz
	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// Injeção manual (sem banco ainda)
	assetRepo := memory.NewAssetMemoryRepo()
	workOrderRepo := memory.NewWorkOrderMemoryRepo()

	assetService := service.NewAssetService(assetRepo)
	workOrderService := service.NewWorkOrderService(workOrderRepo)

	assetHandler := handlers.NewAssetHandler(assetService)
	workOrderHandler := handlers.NewWorkOrderHandler(workOrderService)

	assetHandler.RegisterRoutes(r)
	workOrderHandler.RegisterRoutes(r)

	log.Println("🚀 API running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
