# Datapulse

Datapulse is a cloud-agnostic agent scaffold for:
- Pipeline health analysis (SLA breach, retries, runtime anomalies)
- Cost and data analysis (cost spikes, attribution, schema/volume drift)

It is split into:
- `agent.go` (module root): **ADK Go** conversational onboarding + trigger triage agent with the built-in dev Web UI
- `cmd/mcpserver`: unified MCP server for both pipeline intelligence and cost/data tools

## Conversational onboarding (ADK Web UI)

Uses [Agent Development Kit for Go](https://adk.dev/get-started/go/) (`google.golang.org/adk`) and Gemini. Set **`GOOGLE_API_KEY`** or **`GEMINI_API_KEY`** from [Google AI Studio](https://aistudio.google.com/app/apikey). Optional: **`GEMINI_MODEL`** (default `gemini-2.0-flash`).

```bash
export GOOGLE_API_KEY="your-key"
go run agent.go web api webui
```

Open **http://localhost:8080/ui/** (chat). ADK internally hosts its own endpoints under `/api` for the dev UI. If port 8080 is in use, put `-port` on the `web` launcher **before** the `api` / `webui` subcommands, and align CORS addresses:

```bash
go run agent.go web -port 8095 api -webui_address localhost:8095 webui -api_server_address http://localhost:8095/api
```

With Docker Compose, the `onboarding-agent` service listens on **8095**; pass `GEMINI_API_KEY` or `GOOGLE_API_KEY` in your environment when you run `docker compose up`.

## Run locally

1. Start unified MCP server:
   - `go run ./cmd/mcpserver`
2. Start agent (ADK chat):
   - `go run agent.go web api webui`

Agent MCP default:
- `MCP_SERVER_URL=http://localhost:8091`

## Cloud SQL

Set `CLOUDSQL_DSN` to use Cloud SQL-backed storage.
When empty, in-memory store is used for local development.

Apply schema from `db/schema.sql`.
