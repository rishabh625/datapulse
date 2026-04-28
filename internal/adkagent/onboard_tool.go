package adkagent

import (
	"datapulse/internal/model"
	"datapulse/internal/onboarding"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// OnboardToolOutput is returned to the model after each onboard_cloud_account call.
type OnboardToolOutput struct {
	Success   bool                `json:"success"`
	AccountID string              `json:"account_id,omitempty"`
	Message   string              `json:"message"`
	Checks    []model.CheckResult `json:"checks,omitempty"`
	Warnings  []string            `json:"warnings,omitempty"`
}

// NewOnboardAccountTool wraps the onboarding service as an ADK function tool.
func NewOnboardAccountTool(svc *onboarding.Service) (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name: "onboard_cloud_account",
		Description: "Registers a cloud account with Datapulse after you have collected tenant_id, " +
			"cloud (aws|gcp), account_identifier, trigger_mode, pipelines, and either aws CUR settings " +
			"(cur_bucket, cur_format, read_only_iam_role_arn, region, optional cur_prefix) or gcp billing export fields.",
	}, func(ctx tool.Context, in model.OnboardingRequest) (OnboardToolOutput, error) {
		resp, err := svc.Onboard(ctx, in)
		if err != nil {
			return OnboardToolOutput{
				Success: false,
				Message: err.Error(),
			}, nil
		}
		return OnboardToolOutput{
			Success:   true,
			AccountID: resp.AccountID,
			Message:   "Account registered successfully.",
			Checks:    resp.Checks,
			Warnings:  resp.Warnings,
		}, nil
	})
}
