package service

import (
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository"
)

type AssetService struct {
	repo repository.AssetRepository
}

func NewAssetService(r repository.AssetRepository) *AssetService {
	return &AssetService{repo: r}
}

func (s *AssetService) Create(asset *domain.Asset) error {
	asset.Normalize()
	return s.repo.Create(asset)
}

func (s *AssetService) List() ([]domain.Asset, error) {
	return s.repo.FindAll()
}
