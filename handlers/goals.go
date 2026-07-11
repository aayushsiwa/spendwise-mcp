package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) ListGoals(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := h.Service.ListGoals(ctx)
	if err != nil {
		return errorResult(err)
	}
	return jsonResult(map[string]any{"goals": result})
}

func (h *Handler) GetGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := requireString(req, "goal_id")
	if err != nil {
		return errorResult(err)
	}

	result, err := h.Service.GetGoal(ctx, id)
	if err != nil {
		return errorResult(err)
	}
	return jsonResult(map[string]any{"goal": result})
}
