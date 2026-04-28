package adkagent

import (
	"datapulse/internal/model"
	"datapulse/internal/orchestrator"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type EvaluateTriggerToolOutput struct {
	Success    bool                            `json:"success"`
	Message    string                          `json:"message"`
	Evaluation model.TriggerEvaluationResponse `json:"evaluation,omitempty"`
}

// NewEvaluateTriggerTool exposes orchestrator trigger evaluation to the agent.
func NewEvaluateTriggerTool(svc *orchestrator.Service) (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "evaluate_pipeline_trigger",
		Description: "Evaluates pipeline trigger conditions via MCP signals. " +
			"Requires tenant_id, account_identifier, pipeline_id, and event_status.",
	}, func(ctx tool.Context, in model.TriggerEvaluationRequest) (EvaluateTriggerToolOutput, error) {
		resp, err := svc.EvaluateTrigger(ctx, in)
		if err != nil {
			return EvaluateTriggerToolOutput{
				Success: false,
				Message: err.Error(),
			}, nil
		}
		return EvaluateTriggerToolOutput{
			Success:    true,
			Message:    "Trigger evaluated successfully.",
			Evaluation: resp,
		}, nil
	})
}
