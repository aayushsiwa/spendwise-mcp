package handlers

import (
	"context"

	"aayushsiwa/spendwise-mcp/services"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) ListCategories(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := h.Service.ListCategories(ctx)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]any{"categories": result})
}

func (h *Handler) CreateCategory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := requireString(req, "name")
	if err != nil {
		return errorResult(ctx, err)
	}
	icon, err := optionalString(req, "icon")
	if err != nil {
		return errorResult(ctx, err)
	}
	color, err := optionalString(req, "color")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.CreateCategory(ctx, services.CreateCategoryParams{
		Name:  name,
		Icon:  icon,
		Color: color,
	})
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, result)
}

func (h *Handler) UpdateCategory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	categoryID, err := requireString(req, "category_id")
	if err != nil {
		return errorResult(ctx, err)
	}
	name, err := requireString(req, "name")
	if err != nil {
		return errorResult(ctx, err)
	}
	icon, err := optionalString(req, "icon")
	if err != nil {
		return errorResult(ctx, err)
	}
	color, err := optionalString(req, "color")
	if err != nil {
		return errorResult(ctx, err)
	}

	if err := h.Service.UpdateCategory(ctx, services.UpdateCategoryParams{
		CategoryID: categoryID,
		Name:       name,
		Icon:       icon,
		Color:      color,
	}); err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]string{"message": "category updated"})
}

func (h *Handler) DeleteCategory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	categoryID, err := requireString(req, "category_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.DeleteCategory(ctx, categoryID)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, result)
}
