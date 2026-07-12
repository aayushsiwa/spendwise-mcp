package handlers

import (
	"context"

	"aayushsiwa/spendwise-mcp/services"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) ListBudgets(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	month, err := optionalInt(req, "month")
	if err != nil {
		return errorResult(ctx, err)
	}
	year, err := optionalInt(req, "year")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.ListBudgets(ctx, month, year)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]any{"budgets": result})
}

func (h *Handler) GetBudgetProgress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	month, err := optionalInt(req, "month")
	if err != nil {
		return errorResult(ctx, err)
	}
	year, err := optionalInt(req, "year")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.GetBudgetProgress(ctx, month, year)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]any{"progress": result})
}

func (h *Handler) CreateBudget(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	categoryID, err := requireString(req, "category_id")
	if err != nil {
		return errorResult(ctx, err)
	}
	month, err := requireInt(req, "month")
	if err != nil {
		return errorResult(ctx, err)
	}
	year, err := requireInt(req, "year")
	if err != nil {
		return errorResult(ctx, err)
	}
	amount, err := requireFloat(req, "amount")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.CreateBudget(ctx, services.CreateBudgetParams{
		CategoryID: categoryID,
		Month:      month,
		Year:       year,
		Amount:     amount,
	})
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, result)
}

func (h *Handler) UpdateBudget(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	budgetID, err := requireString(req, "budget_id")
	if err != nil {
		return errorResult(ctx, err)
	}
	amount, err := requireFloat(req, "amount")
	if err != nil {
		return errorResult(ctx, err)
	}

	if err := h.Service.UpdateBudget(ctx, services.UpdateBudgetParams{
		BudgetID: budgetID,
		Amount:   amount,
	}); err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]string{"message": "budget updated"})
}

func (h *Handler) DeleteBudget(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	budgetID, err := requireString(req, "budget_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.DeleteBudget(ctx, budgetID)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, result)
}
