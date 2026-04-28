package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ToolRequest struct {
	Tool      string         `json:"tool"`
	Arguments map[string]any `json:"arguments"`
}

type ToolResponse struct {
	Result map[string]any `json:"result,omitempty"`
	Error  string         `json:"error,omitempty"`
}

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) CallTool(ctx context.Context, tool string, args map[string]any) (map[string]any, error) {
	body, err := json.Marshal(ToolRequest{Tool: tool, Arguments: args})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/mcp", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var output ToolResponse
	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		return nil, err
	}
	if output.Error != "" {
		return nil, fmt.Errorf("tool %s failed: %s", tool, output.Error)
	}
	return output.Result, nil
}
