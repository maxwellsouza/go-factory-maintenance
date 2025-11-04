package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
)

type WorkOrderRepo struct {
	db *DB
}

func NewWorkOrderRepo(db *DB) *WorkOrderRepo {
	return &WorkOrderRepo{db: db}
}

func (r *WorkOrderRepo) Create(order *domain.WorkOrder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO work_orders (asset_id, type, status, title, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.Pool.QueryRow(ctx, query,
		order.AssetID,
		order.Type,
		order.Status,
		order.Title,
		order.Description,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert work order: %w", err)
	}
	return nil
}

func (r *WorkOrderRepo) FindAll() ([]domain.WorkOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			SELECT id, asset_id, type, status, title,
					COALESCE(description,'') AS description,
					breakdown_at, closed_at,
					downtime_minutes,
					COALESCE(cause,'')    AS cause,
					COALESCE(solution,'') AS solution,
					created_at, updated_at
			FROM work_orders
			ORDER BY id;
			`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query work_orders: %w", err)
	}
	defer rows.Close()

	var list []domain.WorkOrder
	for rows.Next() {
		var o domain.WorkOrder
		if err := rows.Scan(
			&o.ID, &o.AssetID, &o.Type, &o.Status, &o.Title, &o.Description,
			&o.BreakdownAt, &o.ClosedAt, &o.DowntimeMinutes,
			&o.Cause, &o.Solution, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan work_order: %w", err)
		}
		list = append(list, o)
	}
	return list, nil
}

func (r *WorkOrderRepo) FindByStatus(status domain.WorkOrderStatus) ([]domain.WorkOrder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			SELECT id, asset_id, type, status, title,
					COALESCE(description,'') AS description,
					breakdown_at, closed_at,
					downtime_minutes,
					COALESCE(cause,'')    AS cause,
					COALESCE(solution,'') AS solution,
					created_at, updated_at
			FROM work_orders
			WHERE status=$1
			ORDER BY id;
			`

	rows, err := r.db.Pool.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("query by status: %w", err)
	}
	defer rows.Close()

	var list []domain.WorkOrder
	for rows.Next() {
		var o domain.WorkOrder
		if err := rows.Scan(
			&o.ID, &o.AssetID, &o.Type, &o.Status, &o.Title, &o.Description,
			&o.BreakdownAt, &o.ClosedAt, &o.DowntimeMinutes,
			&o.Cause, &o.Solution, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan work_order: %w", err)
		}
		list = append(list, o)
	}
	return list, nil
}
