package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

type AssetHandler struct {
	service *service.AssetService
}

func NewAssetHandler(s *service.AssetService) *AssetHandler {
	return &AssetHandler{service: s}
}

func (h *AssetHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/assets")
	group.POST("", h.create)
	group.GET("", h.list)
}

func (h *AssetHandler) create(c *gin.Context) {
	var a domain.Asset
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *AssetHandler) list(c *gin.Context) {
	assets, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}
