package handlers

import (
	"context"

	"aayushsiwa/spendwise-mcp/services"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) ListGoals(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := h.Service.ListGoals(ctx)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]any{"goals": result})
}

func (h *Handler) GetGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := requireString(req, "goal_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.GetGoal(ctx, id)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]any{"goal": result})
}

func (h *Handler) CreateGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := requireString(req, "name")
	if err != nil {
		return errorResult(ctx, err)
	}
	targetAmount, err := requireFloat(req, "target_amount")
	if err != nil {
		return errorResult(ctx, err)
	}
	currentAmount, err := optionalFloat(req, "current_amount")
	if err != nil {
		return errorResult(ctx, err)
	}
	targetDate, err := optionalString(req, "target_date")
	if err != nil {
		return errorResult(ctx, err)
	}
	categoryID, err := optionalString(req, "category_id")
	if err != nil {
		return errorResult(ctx, err)
	}
	status, err := optionalString(req, "status")
	if err != nil {
		return errorResult(ctx, err)
	}
	description, err := optionalString(req, "description")
	if err != nil {
		return errorResult(ctx, err)
	}
	monthlyContribution, err := optionalFloat(req, "monthly_contribution")
	if err != nil {
		return errorResult(ctx, err)
	}

	var catIDPtr *string
	if categoryID != "" {
		catIDPtr = strPtr(categoryID)
	}

	result, err := h.Service.CreateGoal(ctx, services.CreateGoalParams{
		Name:                name,
		TargetAmount:        targetAmount,
		CurrentAmount:       currentAmount,
		TargetDate:          targetDate,
		CategoryID:          catIDPtr,
		Status:              status,
		Description:         description,
		MonthlyContribution: monthlyContribution,
	})
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, result)
}

func (h *Handler) UpdateGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	goalID, err := requireString(req, "goal_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	name, err := optionalString(req, "name")
	if err != nil {
		return errorResult(ctx, err)
	}
	targetAmount, err := optionalFloatPtr(req, "target_amount")
	if err != nil {
		return errorResult(ctx, err)
	}
	currentAmount, err := optionalFloatPtr(req, "current_amount")
	if err != nil {
		return errorResult(ctx, err)
	}
	targetDate, err := optionalString(req, "target_date")
	if err != nil {
		return errorResult(ctx, err)
	}
	category, err := optionalString(req, "category")
	if err != nil {
		return errorResult(ctx, err)
	}
	status, err := optionalString(req, "status")
	if err != nil {
		return errorResult(ctx, err)
	}
	description, err := optionalString(req, "description")
	if err != nil {
		return errorResult(ctx, err)
	}
	monthlyContribution, err := optionalFloatPtr(req, "monthly_contribution")
	if err != nil {
		return errorResult(ctx, err)
	}

	params := services.UpdateGoalParams{GoalID: goalID}
	for _, p := range []struct {
		val string
		set **string
	}{
		{name, &params.Name},
		{targetDate, &params.TargetDate},
		{category, &params.Category},
		{status, &params.Status},
		{description, &params.Description},
	} {
		if p.val != "" {
			v := p.val
			*p.set = &v
		}
	}
	if targetAmount != nil {
		params.TargetAmount = targetAmount
	}
	if currentAmount != nil {
		params.CurrentAmount = currentAmount
	}
	if monthlyContribution != nil {
		params.MonthlyContribution = monthlyContribution
	}

	if err := h.Service.UpdateGoal(ctx, params); err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]string{"message": "goal updated"})
}

func (h *Handler) DeleteGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	goalID, err := requireString(req, "goal_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.DeleteGoal(ctx, goalID)
	if err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, result)
}

func (h *Handler) AddGoalProgress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	goalID, err := requireString(req, "goal_id")
	if err != nil {
		return errorResult(ctx, err)
	}
	amount, err := requireFloat(req, "amount")
	if err != nil {
		return errorResult(ctx, err)
	}

	if err := h.Service.AddGoalProgress(ctx, services.AddGoalProgressParams{
		GoalID: goalID,
		Amount: amount,
	}); err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]string{"message": "goal progress added"})
}
