package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

type WorkOrderHandler struct {
	service *service.WorkOrderService
}

func NewWorkOrderHandler(s *service.WorkOrderService) *WorkOrderHandler {
	return &WorkOrderHandler{service: s}
}

func (h *WorkOrderHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/work-orders")
	group.POST("", h.create)
	group.GET("", h.list)
}

func (h *WorkOrderHandler) create(c *gin.Context) {
	var o domain.WorkOrder
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, o)
}

func (h *WorkOrderHandler) list(c *gin.Context) {
	status := c.Query("status")
	orders, err := h.service.List(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
