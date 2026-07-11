package handlers

import (
	"context"

	"aayushsiwa/spendwise-mcp/models"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) SearchRecords(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fromDate, err := optionalString(req, "from_date")
	if err != nil {
		return errorResult(err)
	}
	toDate, err := optionalString(req, "to_date")
	if err != nil {
		return errorResult(err)
	}
	category, err := optionalString(req, "category")
	if err != nil {
		return errorResult(err)
	}
	recordType, err := optionalString(req, "record_type")
	if err != nil {
		return errorResult(err)
	}
	search, err := optionalString(req, "search")
	if err != nil {
		return errorResult(err)
	}
	groupBy, err := optionalString(req, "group_by")
	if err != nil {
		return errorResult(err)
	}
	page, err := optionalInt(req, "page")
	if err != nil {
		return errorResult(err)
	}
	limit, err := optionalInt(req, "limit")
	if err != nil {
		return errorResult(err)
	}
	minAmount, err := optionalFloat(req, "min_amount")
	if err != nil {
		return errorResult(err)
	}
	maxAmount, err := optionalFloat(req, "max_amount")
	if err != nil {
		return errorResult(err)
	}

	result, err := h.Service.SearchRecords(ctx, models.SearchRecordsParams{
		FromDate:   fromDate,
		ToDate:     toDate,
		Category:   category,
		RecordType: recordType,
		MinAmount:  minAmount,
		MaxAmount:  maxAmount,
		Search:     search,
		GroupBy:    groupBy,
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		return errorResult(err)
	}

	return jsonResult(result)
}

func (h *Handler) GetRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := requireString(req, "record_id")
	if err != nil {
		return errorResult(err)
	}

	result, err := h.Service.GetRecord(ctx, id)
	if err != nil {
		return errorResult(err)
	}

	return jsonResult(result)
}

func (h *Handler) GetFinancialSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fromDate, err := requireString(req, "from_date")
	if err != nil {
		return errorResult(err)
	}
	toDate, err := requireString(req, "to_date")
	if err != nil {
		return errorResult(err)
	}
	category, err := optionalString(req, "category")
	if err != nil {
		return errorResult(err)
	}
	recordType, err := optionalString(req, "record_type")
	if err != nil {
		return errorResult(err)
	}

	result, err := h.Service.GetFinancialSummary(ctx, models.SummaryParams{
		FromDate:   fromDate,
		ToDate:     toDate,
		Category:   category,
		RecordType: recordType,
	})
	if err != nil {
		return errorResult(err)
	}

	return jsonResult(result)
}
