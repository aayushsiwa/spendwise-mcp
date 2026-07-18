package handlers

import (
	"context"

	"aayushsiwa/spendwise-mcp/models"
	"aayushsiwa/spendwise-mcp/services"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handler) SearchRecords(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fromDate, err := optionalString(req, "from_date")
	if err != nil {
		return errorResult(ctx, err)
	}
	toDate, err := optionalString(req, "to_date")
	if err != nil {
		return errorResult(ctx, err)
	}
	category, err := optionalString(req, "category")
	if err != nil {
		return errorResult(ctx, err)
	}
	recordType, err := optionalString(req, "record_type")
	if err != nil {
		return errorResult(ctx, err)
	}
	search, err := optionalString(req, "search")
	if err != nil {
		return errorResult(ctx, err)
	}
	groupBy, err := optionalString(req, "group_by")
	if err != nil {
		return errorResult(ctx, err)
	}
	page, err := optionalInt(req, "page")
	if err != nil {
		return errorResult(ctx, err)
	}
	limit, err := optionalInt(req, "limit")
	if err != nil {
		return errorResult(ctx, err)
	}
	minAmount, err := optionalFloat(req, "min_amount")
	if err != nil {
		return errorResult(ctx, err)
	}
	maxAmount, err := optionalFloat(req, "max_amount")
	if err != nil {
		return errorResult(ctx, err)
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
		return errorResult(ctx, err)
	}

	return successResult(ctx, result)
}

func (h *Handler) CreateRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	date, err := requireString(req, "date")
	if err != nil {
		return errorResult(ctx, err)
	}
	category, err := requireString(req, "category")
	if err != nil {
		return errorResult(ctx, err)
	}
	amount, err := optionalFloat(req, "amount")
	if err != nil {
		return errorResult(ctx, err)
	}
	recordType, err := requireString(req, "record_type")
	if err != nil {
		return errorResult(ctx, err)
	}
	description, err := optionalString(req, "description")
	if err != nil {
		return errorResult(ctx, err)
	}
	note, err := optionalString(req, "note")
	if err != nil {
		return errorResult(ctx, err)
	}
	idempotencyKey, err := optionalString(req, "idempotency_key")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.CreateRecord(ctx, services.CreateRecordParams{
		Date:           date,
		Description:    description,
		Category:       category,
		Amount:         amount,
		RecordType:     recordType,
		Note:           note,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, result)
}

func (h *Handler) GetRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := requireString(req, "record_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.GetRecord(ctx, id)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, result)
}

func (h *Handler) GetFinancialSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fromDate, err := requireString(req, "from_date")
	if err != nil {
		return errorResult(ctx, err)
	}
	toDate, err := requireString(req, "to_date")
	if err != nil {
		return errorResult(ctx, err)
	}
	category, err := optionalString(req, "category")
	if err != nil {
		return errorResult(ctx, err)
	}
	recordType, err := optionalString(req, "record_type")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.GetFinancialSummary(ctx, models.SummaryParams{
		FromDate:   fromDate,
		ToDate:     toDate,
		Category:   category,
		RecordType: recordType,
	})
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, result)
}

func (h *Handler) UpdateRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	recordID, err := requireString(req, "record_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	date, err := optionalString(req, "date")
	if err != nil {
		return errorResult(ctx, err)
	}
	description, err := optionalString(req, "description")
	if err != nil {
		return errorResult(ctx, err)
	}
	category, err := optionalString(req, "category")
	if err != nil {
		return errorResult(ctx, err)
	}
	amount, err := optionalFloatPtr(req, "amount")
	if err != nil {
		return errorResult(ctx, err)
	}
	recordType, err := optionalString(req, "record_type")
	if err != nil {
		return errorResult(ctx, err)
	}
	note, err := optionalString(req, "note")
	if err != nil {
		return errorResult(ctx, err)
	}

	params := services.UpdateRecordParams{RecordID: recordID}
	for _, p := range []struct {
		val string
		set **string
	}{
		{date, &params.Date},
		{description, &params.Description},
		{category, &params.Category},
		{recordType, &params.RecordType},
		{note, &params.Note},
	} {
		if p.val != "" {
			v := p.val
			*p.set = &v
		}
	}
	if amount != nil {
		params.Amount = amount
	}

	if err := h.Service.UpdateRecord(ctx, params); err != nil {
		return errorResult(ctx, err)
	}
	return successResult(ctx, map[string]string{"message": "record updated"})
}

func (h *Handler) DeleteRecord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	recordID, err := requireString(req, "record_id")
	if err != nil {
		return errorResult(ctx, err)
	}

	result, err := h.Service.DeleteRecord(ctx, recordID)
	if err != nil {
		return errorResult(ctx, err)
	}

	return successResult(ctx, result)
}
