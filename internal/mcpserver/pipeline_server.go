package mcpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"datapulse/internal/mcp"
)

// NewServer serves all MCP tools from one endpoint.
func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/mcp", mcpHandler)
	return mux
}

// Backward compatible constructors.
func NewPipelineServer() http.Handler { return NewServer() }
func NewCostServer() http.Handler     { return NewServer() }

func mcpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req mcp.ToolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{Error: err.Error()})
		return
	}

	now := time.Now().UTC()
	switch req.Tool {
	case "get_pipeline_runs":
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{
			Result: map[string]any{
				"pipeline_id": req.Arguments["pipeline_id"],
				"last_run": map[string]any{
					"status":       "FAILED",
					"retries":      2,
					"duration_sec": 1800,
					"sla_breach":   true,
					"finished_at":  now.Format(time.RFC3339),
				},
			},
		})
	case "get_historic_baseline":
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{
			Result: map[string]any{
				"pipeline_id":              req.Arguments["pipeline_id"],
				"avg_duration_sec":         620,
				"p95_duration_sec":         940,
				"retry_rate_percent":       4.7,
				"runtime_anomaly_detected": true,
			},
		})
	case "get_recent_changes":
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{
			Result: map[string]any{
				"pipeline_id": req.Arguments["pipeline_id"],
				"changes": []map[string]any{
					{"at": now.Add(-3 * time.Hour).Format(time.RFC3339), "change": "partition filter removed"},
					{"at": now.Add(-2 * time.Hour).Format(time.RFC3339), "change": "worker count raised to 30"},
				},
			},
		})
	case "detect_cost_spikes":
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{
			Result: map[string]any{
				"pipeline_id":          req.Arguments["pipeline_id"],
				"cost_change_percent":  42.0,
				"detected_at":          now.Format(time.RFC3339),
				"window":               "24h",
				"spike_confidence":     "high",
				"cur_source_validated": true,
			},
		})
	case "attribute_cost_to_pipeline":
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{
			Result: map[string]any{
				"pipeline_id":          req.Arguments["pipeline_id"],
				"reason":               "Glue job reprocessed full history",
				"dpu_hours_multiplier": 2.8,
				"input_volume_x":       3.1,
			},
		})
	case "analyze_data_drift":
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{
			Result: map[string]any{
				"pipeline_id":         req.Arguments["pipeline_id"],
				"schema_drift":        true,
				"new_columns":         []string{"event_source"},
				"volume_change_x":     2.9,
				"late_arriving_ratio": 0.12,
			},
		})
	case "explain_cost_change":
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{
			Result: map[string]any{
				"pipeline_id": req.Arguments["pipeline_id"],
				"explanation": "Cost increased by 42% as Glue job reprocessed full history after partition filter removal.",
			},
		})
	default:
		_ = json.NewEncoder(w).Encode(mcp.ToolResponse{Error: "unsupported tool"})
	}
}
