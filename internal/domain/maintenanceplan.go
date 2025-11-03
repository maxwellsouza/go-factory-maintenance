package domain

import "time"

// PlanRuleType define como a preventiva é disparada.
type PlanRuleType string

const (
	PlanRuleTime      PlanRuleType = "time"      // por período (ex: a cada 30 dias)
	PlanRuleMeter     PlanRuleType = "meter"     // por uso (ex: a cada 10.000 m)
	PlanRuleCondition PlanRuleType = "condition" // por condição/medição
)

type MaintenancePlan struct {
	ID            int64        `json:"id"`
	AssetID       int64        `json:"asset_id"`
	RuleType      PlanRuleType `json:"rule_type"`                // time|meter|condition
	FrequencyDays *int64       `json:"frequency_days,omitempty"` // para "time"
	MeterTarget   *int64       `json:"meter_target,omitempty"`   // para "meter" (se for usar)
	LastExecution *time.Time   `json:"last_execution,omitempty"`
	Active        bool         `json:"active"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (p *MaintenancePlan) Normalize() {
	if !p.Active {
		p.Active = true
	}
}
