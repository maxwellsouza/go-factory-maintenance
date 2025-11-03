package domain

import "time"

// Criticality representa a criticidade operacional do ativo.
type Criticality string

const (
	CriticalityA Criticality = "A"
	CriticalityB Criticality = "B"
	CriticalityC Criticality = "C"
)

type Asset struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Location    string      `json:"location,omitempty"`
	Criticality Criticality `json:"criticality,omitempty"` // A, B, C
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (a *Asset) Normalize() {
	if a.Criticality == "" {
		a.Criticality = CriticalityB
	}
}
