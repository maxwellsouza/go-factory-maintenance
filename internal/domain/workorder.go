package domain

import "time"

// WorkOrderType classifica o tipo de OS.
type WorkOrderType string

const (
	WOTypeCorrective  WorkOrderType = "corrective"
	WOTypePreventive  WorkOrderType = "preventive"
	WOTypeCondition   WorkOrderType = "condition"
	WOTypeImprovement WorkOrderType = "improvement"
)

// WorkOrderStatus representa o estado do ciclo de vida da OS.
type WorkOrderStatus string

const (
	WOStatusOpen       WorkOrderStatus = "open"
	WOStatusInProgress WorkOrderStatus = "in_progress"
	WOStatusDone       WorkOrderStatus = "done"
	WOStatusCanceled   WorkOrderStatus = "canceled"
)

type WorkOrder struct {
	ID              int64           `json:"id"`
	AssetID         int64           `json:"asset_id"`
	Type            WorkOrderType   `json:"type"`                   // corrective|preventive|condition|improvement
	Status          WorkOrderStatus `json:"status"`                 // open|in_progress|done|canceled
	Title           string          `json:"title"`                  // resumo rápido
	Description     string          `json:"description"`            // detalhes
	BreakdownAt     *time.Time      `json:"breakdown_at,omitempty"` // quando a máquina parou (se aplicável)
	ClosedAt        *time.Time      `json:"closed_at,omitempty"`
	DowntimeMinutes *int64          `json:"downtime_minutes,omitempty"`
	Cause           string          `json:"cause,omitempty"`
	Solution        string          `json:"solution,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (wo *WorkOrder) Normalize() {
	if wo.Status == "" {
		wo.Status = WOStatusOpen
	}
	if wo.Type == "" {
		wo.Type = WOTypeCorrective
	}
}
