package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"aayushsiwa/spendwise-mcp/backend"
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
		return nil, fmt.Errorf("limit must be 100 or less")
	}
	if params.GroupBy != "" && params.GroupBy != "category" && params.GroupBy != "month" {
		return nil, fmt.Errorf("group_by must be category or month")
	}
	if err := validateDateRange(params.FromDate, params.ToDate); err != nil {
		return nil, err
	}
	if err := validateRecordType(params.RecordType, true); err != nil {
		return nil, err
	}
	if params.MinAmount > 0 && params.MaxAmount > 0 && params.MinAmount > params.MaxAmount {
		return nil, fmt.Errorf("min_amount must be less than or equal to max_amount")
	}
	return s.client.SearchRecords(ctx, params)
}

func (s *MCPService) GetRecord(ctx context.Context, id string) (*models.Record, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("record_id is required")
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
		return nil, fmt.Errorf("goal_id is required")
	}
	return s.client.GetGoal(ctx, id)
}

func validateDateRange(fromDate, toDate string) error {
	if fromDate == "" || toDate == "" {
		return nil
	}
	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return fmt.Errorf("from_date must be in YYYY-MM-DD format")
	}
	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return fmt.Errorf("to_date must be in YYYY-MM-DD format")
	}
	if from.After(to) {
		return fmt.Errorf("from_date must be before or equal to to_date")
	}
	return nil
}

func requireDateRange(fromDate, toDate string) error {
	if fromDate == "" || toDate == "" {
		return fmt.Errorf("from_date and to_date are required")
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
			return fmt.Errorf("record_type must be income, expense, or transfer")
		}
		return fmt.Errorf("record_type must be income or expense")
	}
	return nil
}

func validateMonthYear(month, year int) error {
	if month < 1 || month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}
	if year < 1 {
		return fmt.Errorf("year must be a positive number")
	}
	return nil
}
