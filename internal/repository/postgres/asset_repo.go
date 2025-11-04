package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
)

type AssetRepo struct {
	db *DB
}

func NewAssetRepo(db *DB) *AssetRepo {
	return &AssetRepo{db: db}
}

func (r *AssetRepo) Create(asset *domain.Asset) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO assets (name, location, criticality, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.Pool.QueryRow(ctx, query, asset.Name, asset.Location, asset.Criticality).
		Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert asset: %w", err)
	}
	return nil
}

func (r *AssetRepo) FindAll() ([]domain.Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name, COALESCE(location,''), criticality, created_at, updated_at
          FROM assets ORDER BY id;`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query assets: %w", err)
	}
	defer rows.Close()

	var assets []domain.Asset
	for rows.Next() {
		var a domain.Asset
		if err := rows.Scan(&a.ID, &a.Name, &a.Location, &a.Criticality, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan asset: %w", err)
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *AssetRepo) FindByID(id int64) (*domain.Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name, COALESCE(location,''), criticality, created_at, updated_at
          FROM assets WHERE id=$1;`

	var a domain.Asset
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&a.ID, &a.Name, &a.Location, &a.Criticality, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find asset: %w", err)
	}
	return &a, nil
}
