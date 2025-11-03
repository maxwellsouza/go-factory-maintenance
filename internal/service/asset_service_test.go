package service_test

import (
	"testing"

	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository/memory"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

func TestAssetService_CreateAndList(t *testing.T) {
	repo := memory.NewAssetMemoryRepo()
	svc := service.NewAssetService(repo)

	tests := []struct {
		name  string
		input domain.Asset
	}{
		{
			name:  "create asset with criticality A",
			input: domain.Asset{Name: "Cortadeira", Location: "Galpão A", Criticality: domain.CriticalityA},
		},
		{
			name:  "create asset with default criticality (should be B)",
			input: domain.Asset{Name: "Rebobinadeira", Location: "Galpão B"}, // sem criticality → Normalize() = B
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := tt.input
			if err := svc.Create(&a); err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			if a.ID == 0 {
				t.Fatalf("expected ID to be set, got %d", a.ID)
			}
			if a.Name == "" {
				t.Fatalf("name should be persisted")
			}
		})
	}

	list, err := svc.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list) != len(tests) {
		t.Fatalf("expected %d assets, got %d", len(tests), len(list))
	}

	// Verifica default da criticidade para o segundo asset
	if list[1].Criticality != domain.CriticalityB {
		t.Fatalf("expected default criticality B, got %s", list[1].Criticality)
	}
}
