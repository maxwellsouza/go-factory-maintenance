package memory

import (
	"sync"
	"time"

	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
)

type WorkOrderMemoryRepo struct {
	data map[int64]*domain.WorkOrder
	mu   sync.RWMutex
	next int64
}

func NewWorkOrderMemoryRepo() *WorkOrderMemoryRepo {
	return &WorkOrderMemoryRepo{
		data: make(map[int64]*domain.WorkOrder),
		next: 1,
	}
}

func (r *WorkOrderMemoryRepo) Create(order *domain.WorkOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = r.next
	r.next++
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt
	r.data[order.ID] = order
	return nil
}

func (r *WorkOrderMemoryRepo) FindAll() ([]domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.WorkOrder, 0, len(r.data))
	for _, o := range r.data {
		result = append(result, *o)
	}
	return result, nil
}

func (r *WorkOrderMemoryRepo) FindByStatus(status domain.WorkOrderStatus) ([]domain.WorkOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := []domain.WorkOrder{}
	for _, o := range r.data {
		if o.Status == status {
			result = append(result, *o)
		}
	}
	return result, nil
}
