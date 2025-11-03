package service

import (
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository"
)

type WorkOrderService struct {
	repo repository.WorkOrderRepository
}

func NewWorkOrderService(r repository.WorkOrderRepository) *WorkOrderService {
	return &WorkOrderService{repo: r}
}

func (s *WorkOrderService) Create(order *domain.WorkOrder) error {
	order.Normalize()
	return s.repo.Create(order)
}

func (s *WorkOrderService) List(status string) ([]domain.WorkOrder, error) {
	if status == "" {
		return s.repo.FindAll()
	}
	return s.repo.FindByStatus(domain.WorkOrderStatus(status))
}
