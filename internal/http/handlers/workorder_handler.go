package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/http/response"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

type WorkOrderHandler struct {
	service *service.WorkOrderService
}

func NewWorkOrderHandler(s *service.WorkOrderService) *WorkOrderHandler {
	return &WorkOrderHandler{service: s}
}

func (h *WorkOrderHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/work-orders")
	g.POST("", h.create)
	g.GET("", h.list)
}

type createWorkOrderRequest struct {
	AssetID     int64                  `json:"asset_id" binding:"required,gt=0"`
	Type        domain.WorkOrderType   `json:"type" binding:"omitempty,oneof=corrective preventive condition improvement"`
	Status      domain.WorkOrderStatus `json:"status" binding:"omitempty,oneof=open in_progress done canceled"`
	Title       string                 `json:"title" binding:"required,min=3"`
	Description string                 `json:"description"`
}

func (h *WorkOrderHandler) create(c *gin.Context) {
	var req createWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	o := domain.WorkOrder{
		AssetID:     req.AssetID,
		Type:        req.Type,
		Status:      req.Status,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.service.Create(&o); err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, o)
}

func (h *WorkOrderHandler) list(c *gin.Context) {
	status := c.Query("status")
	orders, err := h.service.List(status)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}
