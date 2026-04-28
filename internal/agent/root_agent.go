package agent

import (
	"context"

	"datapulse/internal/model"
	"datapulse/internal/onboarding"
	"datapulse/internal/orchestrator"
)

// RootAgent is a thin programmatic boundary for onboarding and triggers.
// For conversational onboarding with the ADK dev Web UI, run the root `agent.go` (see README).
// Keep this interface stable while MCP tools and cloud adapters evolve.
type RootAgent struct {
	onboarding   *onboarding.Service
	orchestrator *orchestrator.Service
}

func NewRootAgent(onboarding *onboarding.Service, orchestrator *orchestrator.Service) *RootAgent {
	return &RootAgent{
		onboarding:   onboarding,
		orchestrator: orchestrator,
	}
}

func (a *RootAgent) HandleOnboarding(ctx context.Context, req model.OnboardingRequest) (model.OnboardingResponse, error) {
	return a.onboarding.Onboard(ctx, req)
}

func (a *RootAgent) HandleTrigger(ctx context.Context, req model.TriggerEvaluationRequest) (model.TriggerEvaluationResponse, error) {
	return a.orchestrator.EvaluateTrigger(ctx, req)
}
