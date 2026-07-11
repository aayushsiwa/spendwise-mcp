package services

import (
	"context"
	"strings"
	"time"

	"aayushsiwa/spendwise-mcp/backend"
	apperrors "aayushsiwa/spendwise-mcp/errors"
	"aayushsiwa/spendwise-mcp/models"
)

type Service interface {
	SearchRecords(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error)
	GetRecord(ctx context.Context, id string) (*models.Record, error)
	GetFinancialSummary(ctx context.Context, params models.SummaryParams) (*models.Summary, error)
	ListCategories(ctx context.Context) ([]models.Category, error)
	ListBudgets(ctx context.Context, month, year int) ([]models.Budget, error)
	GetBudgetProgress(ctx context.Context, month, year int) ([]models.BudgetProgress, error)
	ListGoals(ctx context.Context) ([]models.Goal, error)
	GetGoal(ctx context.Context, id string) (*models.Goal, error)
}

type MCPService struct {
	client backend.Client
}

func NewService(client backend.Client) *MCPService {
	return &MCPService{client: client}
}

func (s *MCPService) SearchRecords(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 25
	}
	if params.Limit > 100 {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"limit": map[string]any{"message": "limit must be 100 or less", "value": params.Limit}})
	}
	if params.GroupBy != "" && params.GroupBy != "category" && params.GroupBy != "month" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"group_by": map[string]any{"message": "group_by must be category or month", "value": params.GroupBy}})
	}
	if err := validateDateRange(params.FromDate, params.ToDate); err != nil {
		return nil, err
	}
	if err := validateRecordType(params.RecordType, true); err != nil {
		return nil, err
	}
	if params.MinAmount > 0 && params.MaxAmount > 0 && params.MinAmount > params.MaxAmount {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"min_amount": map[string]any{"message": "min_amount must be less than or equal to max_amount", "value": params.MinAmount}, "max_amount": map[string]any{"message": "max_amount must be greater than or equal to min_amount", "value": params.MaxAmount}})
	}
	return s.client.SearchRecords(ctx, params)
}

func (s *MCPService) GetRecord(ctx context.Context, id string) (*models.Record, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"record_id": map[string]any{"message": "record_id is required", "value": id}})
	}
	return s.client.GetRecord(ctx, id)
}

func (s *MCPService) GetFinancialSummary(ctx context.Context, params models.SummaryParams) (*models.Summary, error) {
	if err := requireDateRange(params.FromDate, params.ToDate); err != nil {
		return nil, err
	}
	if err := validateRecordType(params.RecordType, false); err != nil {
		return nil, err
	}
	return s.client.GetFinancialSummary(ctx, params)
}

func (s *MCPService) ListCategories(ctx context.Context) ([]models.Category, error) {
	return s.client.ListCategories(ctx)
}

func (s *MCPService) ListBudgets(ctx context.Context, month, year int) ([]models.Budget, error) {
	if err := validateMonthYear(month, year); err != nil {
		return nil, err
	}
	return s.client.ListBudgets(ctx, month, year)
}

func (s *MCPService) GetBudgetProgress(ctx context.Context, month, year int) ([]models.BudgetProgress, error) {
	if err := validateMonthYear(month, year); err != nil {
		return nil, err
	}
	return s.client.GetBudgetProgress(ctx, month, year)
}

func (s *MCPService) ListGoals(ctx context.Context) ([]models.Goal, error) {
	return s.client.ListGoals(ctx)
}

func (s *MCPService) GetGoal(ctx context.Context, id string) (*models.Goal, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"goal_id": map[string]any{"message": "goal_id is required", "value": id}})
	}
	return s.client.GetGoal(ctx, id)
}

func validateDateRange(fromDate, toDate string) error {
	if fromDate == "" || toDate == "" {
		return nil
	}
	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return apperrors.NewValidation("Validation failed", map[string]any{"from_date": map[string]any{"message": "from_date must be in YYYY-MM-DD format", "value": fromDate}})
	}
	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return apperrors.NewValidation("Validation failed", map[string]any{"to_date": map[string]any{"message": "to_date must be in YYYY-MM-DD format", "value": toDate}})
	}
	if from.After(to) {
		return apperrors.NewValidation("Validation failed", map[string]any{"from_date": map[string]any{"message": "from_date must be before or equal to to_date", "value": fromDate}, "to_date": map[string]any{"message": "to_date must be after or equal to from_date", "value": toDate}})
	}
	return nil
}

func requireDateRange(fromDate, toDate string) error {
	if fromDate == "" || toDate == "" {
		return apperrors.NewValidation("Validation failed", map[string]any{"from_date": map[string]any{"message": "from_date is required", "value": fromDate}, "to_date": map[string]any{"message": "to_date is required", "value": toDate}})
	}
	return validateDateRange(fromDate, toDate)
}

func validateRecordType(recordType string, allowTransfer bool) error {
	if recordType == "" {
		return nil
	}
	allowed := map[string]bool{"income": true, "expense": true}
	if allowTransfer {
		allowed["transfer"] = true
	}
	if !allowed[recordType] {
		if allowTransfer {
			return apperrors.NewValidation("Validation failed", map[string]any{"record_type": map[string]any{"message": "record_type must be income, expense, or transfer", "value": recordType}})
		}
		return apperrors.NewValidation("Validation failed", map[string]any{"record_type": map[string]any{"message": "record_type must be income or expense", "value": recordType}})
	}
	return nil
}

func validateMonthYear(month, year int) error {
	if month < 1 || month > 12 {
		return apperrors.NewValidation("Validation failed", map[string]any{"month": map[string]any{"message": "month must be between 1 and 12", "value": month}})
	}
	if year < 1 {
		return apperrors.NewValidation("Validation failed", map[string]any{"year": map[string]any{"message": "year must be a positive number", "value": year}})
	}
	return nil
}
