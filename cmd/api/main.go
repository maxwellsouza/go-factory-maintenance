package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Modo de execução: "release" em produção, "debug" no dev
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	// Middlewares essenciais
	r.Use(gin.Recovery()) // evita crash do servidor em panic
	r.Use(gin.Logger())   // log simples de requests

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	log.Println("🚀 API (Gin) running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
