package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/response"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

type AssetHandler struct {
	service *service.AssetService
}

func NewAssetHandler(s *service.AssetService) *AssetHandler { return &AssetHandler{service: s} }

func (h *AssetHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/assets")
	g.POST("", h.create)
	g.GET("", h.list)
}

// DTO de entrada com validação (não “suje” o domínio com tags binding)
type createAssetRequest struct {
	Name        string             `json:"name" binding:"required,min=2"`
	Location    string             `json:"location"`
	Criticality domain.Criticality `json:"criticality" binding:"omitempty,oneof=A B C"`
}

func (h *AssetHandler) create(c *gin.Context) {
	var req createAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err) // 422 com detalhes
		return
	}

	a := domain.Asset{
		Name:        req.Name,
		Location:    req.Location,
		Criticality: req.Criticality,
	}

	if err := h.service.Create(&a); err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *AssetHandler) list(c *gin.Context) {
	assets, err := h.service.List()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assets)
}
