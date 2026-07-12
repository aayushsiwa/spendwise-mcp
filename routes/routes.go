package routes

import (
	"context"

	"aayushsiwa/spendwise-mcp/handlers"
	"aayushsiwa/spendwise-mcp/session"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Route struct {
	Name        string
	Tool        mcp.Tool
	HandlerFunc handlers.ToolHandler
}

type Routes []Route

type AttachOptions struct {
	ActorID      string
	ClientName   string
	BackendToken string
}

func NewRoutes(h *handlers.Handler) Routes {
	return Routes{
		// ── Records ──────────────────────────────────────────────
		{
			Name: "SearchRecords",
			Tool: mcp.NewTool("search_records",
				mcp.WithDescription("Search financial records with filtering, pagination, and optional grouping."),
				mcp.WithString("from_date", mcp.Description("Start date in YYYY-MM-DD format.")),
				mcp.WithString("to_date", mcp.Description("End date in YYYY-MM-DD format.")),
				mcp.WithString("category", mcp.Description("Category name filter.")),
				mcp.WithString("record_type", mcp.Description("Record type filter."), mcp.Enum("income", "expense", "transfer")),
				mcp.WithNumber("min_amount", mcp.Description("Minimum amount filter.")),
				mcp.WithNumber("max_amount", mcp.Description("Maximum amount filter.")),
				mcp.WithString("search", mcp.Description("Case-insensitive description search.")),
				mcp.WithString("group_by", mcp.Description("Optional grouping."), mcp.Enum("category", "month")),
				mcp.WithNumber("page", mcp.Description("Page number, defaults to 1.")),
				mcp.WithNumber("limit", mcp.Description("Page size, defaults to 25 and max 100.")),
			),
			HandlerFunc: h.SearchRecords,
		},
		{
			Name: "GetRecordDetails",
			Tool: mcp.NewTool("get_record_details",
				mcp.WithDescription("Retrieve one financial record by ID."),
				mcp.WithString("record_id", mcp.Required(), mcp.Description("The record ID.")),
			),
			HandlerFunc: h.GetRecord,
		},
		{
			Name: "CreateRecord",
			Tool: mcp.NewTool("create_spending_record",
				mcp.WithDescription("Create a new financial record (income, expense, or transfer)."),
				mcp.WithString("date", mcp.Required(), mcp.Description("Transaction date in YYYY-MM-DD format.")),
				mcp.WithString("category", mcp.Required(), mcp.Description("Category name (must exist).")),
				mcp.WithNumber("amount", mcp.Required(), mcp.Description("Positive transaction amount.")),
				mcp.WithString("record_type", mcp.Required(), mcp.Description("Record type."), mcp.Enum("income", "expense", "transfer")),
				mcp.WithString("description", mcp.Description("Optional description (max 255 chars).")),
				mcp.WithString("note", mcp.Description("Optional note (max 1000 chars).")),
				mcp.WithString("idempotency_key", mcp.Description("Optional key to prevent duplicate records.")),
			),
			HandlerFunc: h.CreateRecord,
		},
		{
			Name: "UpdateRecord",
			Tool: mcp.NewTool("update_spending_record",
				mcp.WithDescription("Partially update a financial record by ID."),
				mcp.WithString("record_id", mcp.Required(), mcp.Description("The record ID to update.")),
				mcp.WithString("date", mcp.Description("Transaction date in YYYY-MM-DD format.")),
				mcp.WithString("category", mcp.Description("Category name.")),
				mcp.WithNumber("amount", mcp.Description("Positive transaction amount.")),
				mcp.WithString("record_type", mcp.Description("Record type."), mcp.Enum("income", "expense", "transfer")),
				mcp.WithString("description", mcp.Description("Description (max 255 chars).")),
				mcp.WithString("note", mcp.Description("Note.")),
			),
			HandlerFunc: h.UpdateRecord,
		},
		{
			Name: "DeleteRecord",
			Tool: mcp.NewTool("delete_spending_record",
				mcp.WithDescription("Delete a financial record by ID."),
				mcp.WithString("record_id", mcp.Required(), mcp.Description("The record ID to delete.")),
			),
			HandlerFunc: h.DeleteRecord,
		},

		// ── Summary ─────────────────────────────────────────────
		{
			Name: "GetFinancialSummary",
			Tool: mcp.NewTool("get_financial_summary",
				mcp.WithDescription("Return income, expense, net, balances, and category breakdown for a time range."),
				mcp.WithString("from_date", mcp.Required(), mcp.Description("Start date in YYYY-MM-DD format.")),
				mcp.WithString("to_date", mcp.Required(), mcp.Description("End date in YYYY-MM-DD format.")),
				mcp.WithString("category", mcp.Description("Optional category filter.")),
				mcp.WithString("record_type", mcp.Description("Optional summary type filter."), mcp.Enum("income", "expense")),
			),
			HandlerFunc: h.GetFinancialSummary,
		},

		// ── Categories ───────────────────────────────────────────
		{
			Name: "ListCategories",
			Tool: mcp.NewTool("list_categories",
				mcp.WithDescription("List available categories for classification and planning."),
			),
			HandlerFunc: h.ListCategories,
		},

		// ── Budgets ──────────────────────────────────────────────
		{
			Name: "ListBudgets",
			Tool: mcp.NewTool("list_budgets",
				mcp.WithDescription("Retrieve budgets for a specific month and year."),
				mcp.WithNumber("month", mcp.Description("Month from 1 to 12 (defaults to current).")),
				mcp.WithNumber("year", mcp.Description("Four digit year (defaults to current).")),
			),
			HandlerFunc: h.ListBudgets,
		},
		{
			Name: "GetBudgetProgress",
			Tool: mcp.NewTool("get_budget_progress",
				mcp.WithDescription("Return budgets with spent amount and percentage consumed for a given month."),
				mcp.WithNumber("month", mcp.Description("Month from 1 to 12 (defaults to current).")),
				mcp.WithNumber("year", mcp.Description("Four digit year (defaults to current).")),
			),
			HandlerFunc: h.GetBudgetProgress,
		},
		{
			Name: "CreateBudget",
			Tool: mcp.NewTool("create_budget",
				mcp.WithDescription("Create a new budget for a category in a specific month and year."),
				mcp.WithString("category_id", mcp.Required(), mcp.Description("The category ID.")),
				mcp.WithNumber("month", mcp.Required(), mcp.Description("Month (1-12).")),
				mcp.WithNumber("year", mcp.Required(), mcp.Description("Four digit year.")),
				mcp.WithNumber("amount", mcp.Required(), mcp.Description("Budget amount.")),
			),
			HandlerFunc: h.CreateBudget,
		},
		{
			Name: "UpdateBudget",
			Tool: mcp.NewTool("update_budget",
				mcp.WithDescription("Update the amount of an existing budget."),
				mcp.WithString("budget_id", mcp.Required(), mcp.Description("The budget ID.")),
				mcp.WithNumber("amount", mcp.Required(), mcp.Description("New budget amount.")),
			),
			HandlerFunc: h.UpdateBudget,
		},
		{
			Name: "DeleteBudget",
			Tool: mcp.NewTool("delete_budget",
				mcp.WithDescription("Delete a budget by ID."),
				mcp.WithString("budget_id", mcp.Required(), mcp.Description("The budget ID to delete.")),
			),
			HandlerFunc: h.DeleteBudget,
		},

		// ── Goals ────────────────────────────────────────────────
		{
			Name: "ListGoals",
			Tool: mcp.NewTool("list_goals",
				mcp.WithDescription("List savings and financial goals with status and progress."),
			),
			HandlerFunc: h.ListGoals,
		},
		{
			Name: "GetGoalDetails",
			Tool: mcp.NewTool("get_goal_details",
				mcp.WithDescription("Retrieve one goal by ID."),
				mcp.WithString("goal_id", mcp.Required(), mcp.Description("The goal ID.")),
			),
			HandlerFunc: h.GetGoal,
		},
		{
			Name: "CreateGoal",
			Tool: mcp.NewTool("create_goal",
				mcp.WithDescription("Create a new financial goal."),
				mcp.WithString("name", mcp.Required(), mcp.Description("Goal name.")),
				mcp.WithNumber("target_amount", mcp.Required(), mcp.Description("Target amount to save.")),
				mcp.WithNumber("current_amount", mcp.Description("Current saved amount (defaults to 0).")),
				mcp.WithString("target_date", mcp.Description("Target completion date in YYYY-MM-DD format.")),
				mcp.WithString("category_id", mcp.Description("Optional category ID.")),
				mcp.WithString("status", mcp.Description("Goal status."), mcp.Enum("active", "achieved", "abandoned")),
				mcp.WithString("description", mcp.Description("Optional description.")),
				mcp.WithNumber("monthly_contribution", mcp.Description("Optional monthly contribution amount.")),
			),
			HandlerFunc: h.CreateGoal,
		},
		{
			Name: "UpdateGoal",
			Tool: mcp.NewTool("update_goal",
				mcp.WithDescription("Partially update a financial goal."),
				mcp.WithString("goal_id", mcp.Required(), mcp.Description("The goal ID to update.")),
				mcp.WithString("name", mcp.Description("Goal name.")),
				mcp.WithNumber("target_amount", mcp.Description("Target amount to save.")),
				mcp.WithNumber("current_amount", mcp.Description("Current saved amount.")),
				mcp.WithString("target_date", mcp.Description("Target completion date in YYYY-MM-DD format.")),
				mcp.WithString("category", mcp.Description("Category name.")),
				mcp.WithString("status", mcp.Description("Goal status."), mcp.Enum("active", "achieved", "abandoned")),
				mcp.WithString("description", mcp.Description("Optional description.")),
				mcp.WithNumber("monthly_contribution", mcp.Description("Monthly contribution amount.")),
			),
			HandlerFunc: h.UpdateGoal,
		},
		{
			Name: "DeleteGoal",
			Tool: mcp.NewTool("delete_goal",
				mcp.WithDescription("Delete a goal by ID."),
				mcp.WithString("goal_id", mcp.Required(), mcp.Description("The goal ID to delete.")),
			),
			HandlerFunc: h.DeleteGoal,
		},
		{
			Name: "AddGoalProgress",
			Tool: mcp.NewTool("add_goal_progress",
				mcp.WithDescription("Add progress (contribution) toward a financial goal."),
				mcp.WithString("goal_id", mcp.Required(), mcp.Description("The goal ID.")),
				mcp.WithNumber("amount", mcp.Required(), mcp.Description("Amount to add to current progress.")),
			),
			HandlerFunc: h.AddGoalProgress,
		},
	}
}

func AttachRoutes(mcpServer *server.MCPServer, routes Routes, opts AttachOptions) {
	for _, route := range routes {
		toolName := route.Tool.Name
		handlerFunc := route.HandlerFunc
		mcpServer.AddTool(route.Tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx = session.WithContext(ctx, session.New(opts.ActorID, opts.ClientName, opts.BackendToken, toolName))
			return handlerFunc(ctx, req)
		})
	}
}
