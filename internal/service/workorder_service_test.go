package service_test

import (
	"testing"

	"github.com/maxwellsouza/go-factory-maintenance/internal/domain"
	"github.com/maxwellsouza/go-factory-maintenance/internal/repository/memory"
	"github.com/maxwellsouza/go-factory-maintenance/internal/service"
)

func TestWorkOrderService_CreateAndListByStatus(t *testing.T) {
	repo := memory.NewWorkOrderMemoryRepo()
	svc := service.NewWorkOrderService(repo)

	cases := []struct {
		name  string
		input domain.WorkOrder
	}{
		{
			name:  "create corrective open by default",
			input: domain.WorkOrder{AssetID: 1, Title: "Trocar rolete", Description: "Barulho no rolo"},
		},
		{
			name:  "create preventive explicit open",
			input: domain.WorkOrder{AssetID: 1, Type: domain.WOTypePreventive, Status: domain.WOStatusOpen, Title: "Preventiva mensal"},
		},
		{
			name:  "create in_progress",
			input: domain.WorkOrder{AssetID: 1, Type: domain.WOTypeCorrective, Status: domain.WOStatusInProgress, Title: "Ajuste de correia"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			o := tc.input
			if err := svc.Create(&o); err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			if o.ID == 0 {
				t.Fatalf("expected ID to be set")
			}
		})
	}

	// List sem filtro → todos
	all, err := svc.List("")
	if err != nil {
		t.Fatalf("List(\"\") error = %v", err)
	}
	if len(all) != len(cases) {
		t.Fatalf("expected %d work orders, got %d", len(cases), len(all))
	}

	// Filtro por status open
	open, err := svc.List(string(domain.WOStatusOpen))
	if err != nil {
		t.Fatalf("List(open) error = %v", err)
	}
	if len(open) < 2 {
		t.Fatalf("expected at least 2 open work orders, got %d", len(open))
	}

	// Verifica defaults do primeiro: type=corrective, status=open
	if all[0].Type != domain.WOTypeCorrective {
		t.Fatalf("expected default type corrective, got %s", all[0].Type)
	}
	if all[0].Status != domain.WOStatusOpen {
		t.Fatalf("expected default status open, got %s", all[0].Status)
	}
}
