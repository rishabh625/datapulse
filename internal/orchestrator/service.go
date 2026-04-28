package orchestrator

import (
	"context"
	"fmt"

	"datapulse/internal/mcp"
	"datapulse/internal/model"
)

type Service struct {
	mcpClient *mcp.Client
}

func NewService(mcpClient *mcp.Client) *Service {
	return &Service{
		mcpClient: mcpClient,
	}
}

func (s *Service) EvaluateTrigger(ctx context.Context, req model.TriggerEvaluationRequest) (model.TriggerEvaluationResponse, error) {
	runs, err := s.mcpClient.CallTool(ctx, "get_pipeline_runs", map[string]any{
		"tenant_id":          req.TenantID,
		"account_identifier": req.AccountIdentifier,
		"pipeline_id":        req.PipelineID,
	})
	if err != nil {
		return model.TriggerEvaluationResponse{}, err
	}

	baseline, err := s.mcpClient.CallTool(ctx, "get_historic_baseline", map[string]any{
		"tenant_id":          req.TenantID,
		"account_identifier": req.AccountIdentifier,
		"pipeline_id":        req.PipelineID,
		"lookback_days":      14,
	})
	if err != nil {
		return model.TriggerEvaluationResponse{}, err
	}

	changes, err := s.mcpClient.CallTool(ctx, "get_recent_changes", map[string]any{
		"tenant_id":          req.TenantID,
		"account_identifier": req.AccountIdentifier,
		"pipeline_id":        req.PipelineID,
	})
	if err != nil {
		return model.TriggerEvaluationResponse{}, err
	}

	spikes, err := s.mcpClient.CallTool(ctx, "detect_cost_spikes", map[string]any{
		"tenant_id":          req.TenantID,
		"account_identifier": req.AccountIdentifier,
		"pipeline_id":        req.PipelineID,
	})
	if err != nil {
		return model.TriggerEvaluationResponse{}, err
	}

	attribution, err := s.mcpClient.CallTool(ctx, "attribute_cost_to_pipeline", map[string]any{
		"tenant_id":          req.TenantID,
		"account_identifier": req.AccountIdentifier,
		"pipeline_id":        req.PipelineID,
	})
	if err != nil {
		return model.TriggerEvaluationResponse{}, err
	}

	dataDrift, err := s.mcpClient.CallTool(ctx, "analyze_data_drift", map[string]any{
		"tenant_id":          req.TenantID,
		"account_identifier": req.AccountIdentifier,
		"pipeline_id":        req.PipelineID,
	})
	if err != nil {
		return model.TriggerEvaluationResponse{}, err
	}

	explanation := fmt.Sprintf(
		"Pipeline %s (%s) shows runtime_anomaly=%v and cost_change=%v%%. Possible reason: %v. Recent changes: %v",
		req.PipelineID,
		req.EventStatus,
		baseline["runtime_anomaly_detected"],
		spikes["cost_change_percent"],
		attribution["reason"],
		changes["changes"],
	)

	return model.TriggerEvaluationResponse{
		PipelineSummary: map[string]any{
			"runs":     runs,
			"baseline": baseline,
			"changes":  changes,
		},
		CostSummary: map[string]any{
			"spikes":      spikes,
			"attribution": attribution,
			"data_drift":  dataDrift,
		},
		Explanation: explanation,
	}, nil
}
