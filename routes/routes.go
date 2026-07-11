package routes

import (
	"context"

	"aayushsiwa/spendwise-mcp/handlers"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Route struct {
	Name        string
	Tool        mcp.Tool
	HandlerFunc handlers.ToolHandler
}

type Routes []Route

func NewRoutes(h *handlers.Handler) Routes {
	return Routes{
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
		{
			Name: "ListCategories",
			Tool: mcp.NewTool("list_categories",
				mcp.WithDescription("List available categories for classification and planning."),
			),
			HandlerFunc: h.ListCategories,
		},
		{
			Name: "ListBudgets",
			Tool: mcp.NewTool("list_budgets",
				mcp.WithDescription("Retrieve budgets for a specific month and year."),
				mcp.WithNumber("month", mcp.Required(), mcp.Description("Month from 1 to 12.")),
				mcp.WithNumber("year", mcp.Required(), mcp.Description("Four digit year.")),
			),
			HandlerFunc: h.ListBudgets,
		},
		{
			Name: "GetBudgetProgress",
			Tool: mcp.NewTool("get_budget_progress",
				mcp.WithDescription("Return budgets with spent amount and percentage consumed for a given month."),
				mcp.WithNumber("month", mcp.Required(), mcp.Description("Month from 1 to 12.")),
				mcp.WithNumber("year", mcp.Required(), mcp.Description("Four digit year.")),
			),
			HandlerFunc: h.GetBudgetProgress,
		},
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
	}
}

func AttachRoutes(mcpServer *server.MCPServer, routes Routes) {
	for _, route := range routes {
		handlerFunc := route.HandlerFunc
		mcpServer.AddTool(route.Tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handlerFunc(ctx, req)
		})
	}
}
