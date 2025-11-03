package memory

import (
	"sync"
	"time"

	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
)

type AssetMemoryRepo struct {
	data map[int64]*domain.Asset
	mu   sync.RWMutex
	next int64
}

func NewAssetMemoryRepo() *AssetMemoryRepo {
	return &AssetMemoryRepo{
		data: make(map[int64]*domain.Asset),
		next: 1,
	}
}

func (r *AssetMemoryRepo) Create(asset *domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	asset.ID = r.next
	r.next++
	asset.CreatedAt = time.Now()
	asset.UpdatedAt = asset.CreatedAt
	r.data[asset.ID] = asset
	return nil
}

func (r *AssetMemoryRepo) FindAll() ([]domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Asset, 0, len(r.data))
	for _, a := range r.data {
		result = append(result, *a)
	}
	return result, nil
}

func (r *AssetMemoryRepo) FindByID(id int64) (*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.data[id]; ok {
		return a, nil
	}
	return nil, domain.ErrNotFound
}
