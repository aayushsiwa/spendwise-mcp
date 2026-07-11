package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) ListBudgets(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	month, err := optionalInt(req, "month")
	if err != nil {
		return errorResult(err)
	}
	year, err := optionalInt(req, "year")
	if err != nil {
		return errorResult(err)
	}

	result, err := h.Service.ListBudgets(ctx, month, year)
	if err != nil {
		return errorResult(err)
	}
	return jsonResult(map[string]any{"budgets": result})
}

func (h *Handler) GetBudgetProgress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	month, err := optionalInt(req, "month")
	if err != nil {
		return errorResult(err)
	}
	year, err := optionalInt(req, "year")
	if err != nil {
		return errorResult(err)
	}

	result, err := h.Service.GetBudgetProgress(ctx, month, year)
	if err != nil {
		return errorResult(err)
	}
	return jsonResult(map[string]any{"progress": result})
}
