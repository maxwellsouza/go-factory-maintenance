package repository

import "github.com/maxwellsouza/go-factory-maintenance/internal/domain"

type AssetRepository interface {
	Create(asset *domain.Asset) error
	FindAll() ([]domain.Asset, error)
	FindByID(id int64) (*domain.Asset, error)
}

type WorkOrderRepository interface {
	Create(order *domain.WorkOrder) error
	FindAll() ([]domain.WorkOrder, error)
	FindByStatus(status domain.WorkOrderStatus) ([]domain.WorkOrder, error)
}
