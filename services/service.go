package services

import (
	"context"
	"math"
	"strings"
	"time"

	"aayushsiwa/spendwise-mcp/backend"
	apperrors "aayushsiwa/spendwise-mcp/errors"
	"aayushsiwa/spendwise-mcp/models"
)

type CreateRecordParams struct {
	Date           string
	Description    string
	Category       string
	Amount         float64
	RecordType     string
	Note           string
	IdempotencyKey string
}

type UpdateRecordParams struct {
	RecordID    string
	Date        *string
	Description *string
	Category    *string
	Amount      *float64
	RecordType  *string
	Note        *string
}

type CreateBudgetParams struct {
	CategoryID string
	Month      int
	Year       int
	Amount     float64
}

type UpdateBudgetParams struct {
	BudgetID string
	Amount   float64
}

type CreateGoalParams struct {
	Name                string
	TargetAmount        float64
	CurrentAmount       float64
	TargetDate          string
	CategoryID          *string
	Status              string
	Description         string
	MonthlyContribution float64
}

type UpdateGoalParams struct {
	GoalID              string
	Name                *string
	TargetAmount        *float64
	CurrentAmount       *float64
	TargetDate          *string
	Category            *string
	Status              *string
	Description         *string
	MonthlyContribution *float64
}

type AddGoalProgressParams struct {
	GoalID string
	Amount float64
}

type CreateCategoryParams struct {
	Name  string
	Icon  string
	Color string
}

type UpdateCategoryParams struct {
	CategoryID string
	Name       string
	Icon       string
	Color      string
}

type Service interface {
	SearchRecords(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error)
	GetRecord(ctx context.Context, id string) (*models.Record, error)
	GetFinancialSummary(ctx context.Context, params models.SummaryParams) (*models.Summary, error)
	ListCategories(ctx context.Context) ([]models.Category, error)
	ListBudgets(ctx context.Context, month, year int) ([]models.Budget, error)
	GetBudgetProgress(ctx context.Context, month, year int) ([]models.BudgetProgress, error)
	ListGoals(ctx context.Context) ([]models.Goal, error)
	GetGoal(ctx context.Context, id string) (*models.Goal, error)
	CreateRecord(ctx context.Context, params CreateRecordParams) (*backend.CreateRecordOutput, error)
	UpdateRecord(ctx context.Context, params UpdateRecordParams) error
	DeleteRecord(ctx context.Context, id string) (*backend.DeleteRecordOutput, error)
	CreateBudget(ctx context.Context, params CreateBudgetParams) (*backend.CreateBudgetOutput, error)
	UpdateBudget(ctx context.Context, params UpdateBudgetParams) error
	DeleteBudget(ctx context.Context, id string) (*backend.DeleteBudgetOutput, error)
	CreateGoal(ctx context.Context, params CreateGoalParams) (*backend.CreateGoalOutput, error)
	UpdateGoal(ctx context.Context, params UpdateGoalParams) error
	DeleteGoal(ctx context.Context, id string) (*backend.DeleteBudgetOutput, error)
	AddGoalProgress(ctx context.Context, params AddGoalProgressParams) error
	CreateCategory(ctx context.Context, params CreateCategoryParams) (*backend.CreateCategoryOutput, error)
	UpdateCategory(ctx context.Context, params UpdateCategoryParams) error
	DeleteCategory(ctx context.Context, id string) (*backend.DeleteCategoryOutput, error)
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

func (s *MCPService) CreateRecord(ctx context.Context, params CreateRecordParams) (*backend.CreateRecordOutput, error) {
	if strings.TrimSpace(params.Date) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"date": map[string]any{"message": "date is required"}})
	}
	if strings.TrimSpace(params.Category) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"category": map[string]any{"message": "category is required"}})
	}
	if params.Amount <= 0 {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"amount": map[string]any{"message": "amount must be greater than 0", "value": params.Amount}})
	}
	if err := validateRecordType(params.RecordType, true); err != nil {
		return nil, err
	}
	if _, err := time.Parse("2006-01-02", params.Date); err != nil {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"date": map[string]any{"message": "date must be in YYYY-MM-DD format", "value": params.Date}})
	}
	if params.Description != "" && len(params.Description) > 255 {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"description": map[string]any{"message": "description must be 255 characters or less"}})
	}
	return s.client.CreateRecord(ctx, backend.CreateRecordInput{
		Date:        params.Date,
		Description: params.Description,
		Category:    params.Category,
		Amount:      params.Amount,
		Type:        params.RecordType,
		Note:        params.Note,
	})
}

func (s *MCPService) UpdateRecord(ctx context.Context, params UpdateRecordParams) error {
	if strings.TrimSpace(params.RecordID) == "" {
		return apperrors.NewValidation("Validation failed", map[string]any{"record_id": map[string]any{"message": "record_id is required"}})
	}

	input := backend.UpdateRecordInput{}

	if params.Date != nil {
		if strings.TrimSpace(*params.Date) == "" {
			return apperrors.NewValidation("Validation failed", map[string]any{"date": map[string]any{"message": "date must not be empty"}})
		}
		if _, err := time.Parse("2006-01-02", *params.Date); err != nil {
			return apperrors.NewValidation("Validation failed", map[string]any{"date": map[string]any{"message": "date must be in YYYY-MM-DD format", "value": *params.Date}})
		}
		input.Date = params.Date
	}
	if params.Description != nil {
		if len(*params.Description) > 255 {
			return apperrors.NewValidation("Validation failed", map[string]any{"description": map[string]any{"message": "description must be 255 characters or less"}})
		}
		input.Description = params.Description
	}
	if params.Category != nil {
		if strings.TrimSpace(*params.Category) == "" {
			return apperrors.NewValidation("Validation failed", map[string]any{"category": map[string]any{"message": "category must not be empty"}})
		}
		input.Category = params.Category
	}
	if params.Amount != nil {
		if *params.Amount <= 0 {
			return apperrors.NewValidation("Validation failed", map[string]any{"amount": map[string]any{"message": "amount must be greater than 0", "value": *params.Amount}})
		}
		input.Amount = params.Amount
	}
	if params.RecordType != nil {
		if err := validateRecordType(*params.RecordType, true); err != nil {
			return err
		}
		input.Type = params.RecordType
	}
	if params.Note != nil {
		input.Note = params.Note
	}

	return s.client.UpdateRecord(ctx, params.RecordID, input)
}

func (s *MCPService) DeleteRecord(ctx context.Context, id string) (*backend.DeleteRecordOutput, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"record_id": map[string]any{"message": "record_id is required", "value": id}})
	}
	return s.client.DeleteRecord(ctx, id)
}

func (s *MCPService) CreateBudget(ctx context.Context, params CreateBudgetParams) (*backend.CreateBudgetOutput, error) {
	if err := validateMonthYear(params.Month, params.Year); err != nil {
		return nil, err
	}
	if strings.TrimSpace(params.CategoryID) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"categoryID": map[string]any{"message": "categoryID is required"}})
	}
	if params.Amount <= 0 || math.IsInf(params.Amount, 0) || math.IsNaN(params.Amount) {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"amount": map[string]any{"message": "amount must be a positive number", "value": params.Amount}})
	}
	return s.client.CreateBudget(ctx, backend.CreateBudgetInput{
		CategoryID: params.CategoryID,
		Month:      params.Month,
		Year:       params.Year,
		Amount:     params.Amount,
	})
}

func (s *MCPService) UpdateBudget(ctx context.Context, params UpdateBudgetParams) error {
	if strings.TrimSpace(params.BudgetID) == "" {
		return apperrors.NewValidation("Validation failed", map[string]any{"budget_id": map[string]any{"message": "budget_id is required"}})
	}
	if params.Amount <= 0 || math.IsInf(params.Amount, 0) || math.IsNaN(params.Amount) {
		return apperrors.NewValidation("Validation failed", map[string]any{"amount": map[string]any{"message": "amount must be a positive number", "value": params.Amount}})
	}
	return s.client.UpdateBudget(ctx, params.BudgetID, backend.UpdateBudgetInput{Amount: params.Amount})
}

func (s *MCPService) DeleteBudget(ctx context.Context, id string) (*backend.DeleteBudgetOutput, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"budget_id": map[string]any{"message": "budget_id is required", "value": id}})
	}
	return s.client.DeleteBudget(ctx, id)
}

func (s *MCPService) CreateGoal(ctx context.Context, params CreateGoalParams) (*backend.CreateGoalOutput, error) {
	if strings.TrimSpace(params.Name) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"name": map[string]any{"message": "name is required"}})
	}
	if params.TargetAmount <= 0 || math.IsInf(params.TargetAmount, 0) || math.IsNaN(params.TargetAmount) {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"targetAmount": map[string]any{"message": "targetAmount must be a positive number", "value": params.TargetAmount}})
	}
	if params.CurrentAmount < 0 {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"currentAmount": map[string]any{"message": "currentAmount must not be negative", "value": params.CurrentAmount}})
	}
	if params.TargetDate != "" {
		if _, err := time.Parse("2006-01-02", params.TargetDate); err != nil {
			return nil, apperrors.NewValidation("Validation failed", map[string]any{"targetDate": map[string]any{"message": "targetDate must be in YYYY-MM-DD format", "value": params.TargetDate}})
		}
	}
	if params.Status != "" && params.Status != "active" && params.Status != "achieved" && params.Status != "abandoned" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"status": map[string]any{"message": "status must be active, achieved, or abandoned", "value": params.Status}})
	}
	if params.Description != "" && len(params.Description) > 1000 {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"description": map[string]any{"message": "description must be 1000 characters or less"}})
	}
	return s.client.CreateGoal(ctx, backend.CreateGoalInput{
		Name:                params.Name,
		TargetAmount:        params.TargetAmount,
		CurrentAmount:       params.CurrentAmount,
		TargetDate:          params.TargetDate,
		CategoryID:          params.CategoryID,
		Status:              params.Status,
		Description:         params.Description,
		MonthlyContribution: &params.MonthlyContribution,
	})
}

func (s *MCPService) UpdateGoal(ctx context.Context, params UpdateGoalParams) error {
	if strings.TrimSpace(params.GoalID) == "" {
		return apperrors.NewValidation("Validation failed", map[string]any{"goal_id": map[string]any{"message": "goal_id is required"}})
	}

	input := backend.UpdateGoalInput{}

	if params.Name != nil {
		if strings.TrimSpace(*params.Name) == "" {
			return apperrors.NewValidation("Validation failed", map[string]any{"name": map[string]any{"message": "name must not be empty"}})
		}
		input.Name = params.Name
	}
	if params.TargetAmount != nil {
		if *params.TargetAmount <= 0 || math.IsInf(*params.TargetAmount, 0) || math.IsNaN(*params.TargetAmount) {
			return apperrors.NewValidation("Validation failed", map[string]any{"targetAmount": map[string]any{"message": "targetAmount must be a positive number", "value": *params.TargetAmount}})
		}
		input.TargetAmount = params.TargetAmount
	}
	if params.CurrentAmount != nil {
		if *params.CurrentAmount < 0 {
			return apperrors.NewValidation("Validation failed", map[string]any{"currentAmount": map[string]any{"message": "currentAmount must not be negative", "value": *params.CurrentAmount}})
		}
		input.CurrentAmount = params.CurrentAmount
	}
	if params.TargetDate != nil {
		if *params.TargetDate != "" {
			if _, err := time.Parse("2006-01-02", *params.TargetDate); err != nil {
				return apperrors.NewValidation("Validation failed", map[string]any{"targetDate": map[string]any{"message": "targetDate must be in YYYY-MM-DD format", "value": *params.TargetDate}})
			}
		}
		input.TargetDate = params.TargetDate
	}
	if params.Category != nil {
		input.Category = params.Category
	}
	if params.Status != nil {
		if *params.Status != "active" && *params.Status != "achieved" && *params.Status != "abandoned" {
			return apperrors.NewValidation("Validation failed", map[string]any{"status": map[string]any{"message": "status must be active, achieved, or abandoned", "value": *params.Status}})
		}
		input.Status = params.Status
	}
	if params.Description != nil {
		input.Description = params.Description
	}
	if params.MonthlyContribution != nil {
		if *params.MonthlyContribution < 0 {
			return apperrors.NewValidation("Validation failed", map[string]any{"monthlyContribution": map[string]any{"message": "monthlyContribution must not be negative", "value": *params.MonthlyContribution}})
		}
		input.MonthlyContribution = params.MonthlyContribution
	}

	return s.client.UpdateGoal(ctx, params.GoalID, input)
}

func (s *MCPService) DeleteGoal(ctx context.Context, id string) (*backend.DeleteBudgetOutput, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"goal_id": map[string]any{"message": "goal_id is required", "value": id}})
	}
	return s.client.DeleteGoal(ctx, id)
}

func (s *MCPService) CreateCategory(ctx context.Context, params CreateCategoryParams) (*backend.CreateCategoryOutput, error) {
	if strings.TrimSpace(params.Name) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"name": map[string]any{"message": "name is required"}})
	}
	return s.client.CreateCategory(ctx, backend.CreateCategoryInput{
		Name:  params.Name,
		Icon:  params.Icon,
		Color: params.Color,
	})
}

func (s *MCPService) UpdateCategory(ctx context.Context, params UpdateCategoryParams) error {
	if strings.TrimSpace(params.CategoryID) == "" {
		return apperrors.NewValidation("Validation failed", map[string]any{"category_id": map[string]any{"message": "category_id is required"}})
	}
	if strings.TrimSpace(params.Name) == "" {
		return apperrors.NewValidation("Validation failed", map[string]any{"name": map[string]any{"message": "name is required"}})
	}
	return s.client.UpdateCategory(ctx, params.CategoryID, backend.UpdateCategoryInput{
		Name:  params.Name,
		Icon:  params.Icon,
		Color: params.Color,
	})
}

func (s *MCPService) DeleteCategory(ctx context.Context, id string) (*backend.DeleteCategoryOutput, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperrors.NewValidation("Validation failed", map[string]any{"category_id": map[string]any{"message": "category_id is required", "value": id}})
	}
	return s.client.DeleteCategory(ctx, id)
}

func (s *MCPService) AddGoalProgress(ctx context.Context, params AddGoalProgressParams) error {
	if strings.TrimSpace(params.GoalID) == "" {
		return apperrors.NewValidation("Validation failed", map[string]any{"goal_id": map[string]any{"message": "goal_id is required"}})
	}
	if params.Amount <= 0 || math.IsInf(params.Amount, 0) || math.IsNaN(params.Amount) {
		return apperrors.NewValidation("Validation failed", map[string]any{"amount": map[string]any{"message": "amount must be a positive number", "value": params.Amount}})
	}
	return s.client.AddGoalProgress(ctx, params.GoalID, backend.AddGoalProgressInput{Amount: params.Amount})
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
