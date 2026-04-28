// Datapulse onboarding agent (Google ADK for Go).
//
// Dev web UI (ADK built-in chat):
//
//	export GOOGLE_API_KEY=...   # or GEMINI_API_KEY
//	go run agent.go web api webui
//
// Chat UI: http://localhost:8080/ui/ — if port 8080 is taken, use e.g.:
//
//	go run agent.go web -port 8095 api -webui_address localhost:8095 webui -api_server_address http://localhost:8095/api
package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"datapulse/internal/adkagent"
	"datapulse/internal/mcp"
	"datapulse/internal/onboarding"
	"datapulse/internal/orchestrator"
	"datapulse/internal/store"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	adklauncher "google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	adkgemini "google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		log.Fatal("set GOOGLE_API_KEY or GEMINI_API_KEY to your Gemini API key (https://aistudio.google.com/app/apikey)")
	}

	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}

	m, err := adkgemini.NewModel(ctx, modelName, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("gemini model: %v", err)
	}

	accountStore, cleanup, err := initStore()
	if err != nil {
		log.Fatalf("init store: %v", err)
	}
	defer cleanup()

	onb := onboarding.NewService(accountStore)
	mcpURL := os.Getenv("MCP_SERVER_URL")
	if mcpURL == "" {
		mcpURL = "http://localhost:8091"
	}
	orch := orchestrator.NewService(mcp.NewClient(mcpURL))

	onboardTool, err := adkagent.NewOnboardAccountTool(onb)
	if err != nil {
		log.Fatalf("onboard tool: %v", err)
	}
	evaluateTool, err := adkagent.NewEvaluateTriggerTool(orch)
	if err != nil {
		log.Fatalf("evaluate tool: %v", err)
	}

	ag, err := llmagent.New(llmagent.Config{
		Name:        "datapulse_onboarding",
		Model:       m,
		Description: "Conversational onboarding plus MCP-backed trigger triage for Datapulse.",
		Instruction: adkagent.OnboardingInstruction,
		Tools:       []tool.Tool{onboardTool, evaluateTool},
	})
	if err != nil {
		log.Fatalf("agent: %v", err)
	}

	cfg := &adklauncher.Config{
		AgentLoader: agent.NewSingleLoader(ag),
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, cfg, os.Args[1:]); err != nil {
		log.Fatalf("run: %v\n\n%s", err, l.CommandLineSyntax())
	}
}

func initStore() (store.AccountStore, func(), error) {
	dsn := os.Getenv("CLOUDSQL_DSN")
	if dsn == "" {
		return store.NewMemoryStore(), func() {}, nil
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, err
	}
	return store.NewCloudSQLStore(db), func() { _ = db.Close() }, nil
}
