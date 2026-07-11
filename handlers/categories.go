package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) ListCategories(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := h.Service.ListCategories(ctx)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]any{"categories": result})
}
