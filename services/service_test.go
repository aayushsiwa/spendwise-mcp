package services_test

import (
	"context"
	"testing"

	"aayushsiwa/spendwise-mcp/backend"
	"aayushsiwa/spendwise-mcp/models"
	"aayushsiwa/spendwise-mcp/services"
)

type mockClient struct {
	backend.Client
	searchRecordsFn func(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error)
}

func (m *mockClient) SearchRecords(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error) {
	return m.searchRecordsFn(ctx, params)
}

func TestSearchRecords_DefaultsApplied(t *testing.T) {
	var captured models.SearchRecordsParams
	client := &mockClient{
		searchRecordsFn: func(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error) {
			captured = params
			return &models.SearchRecordsResult{}, nil
		},
	}
	svc := services.NewService(client)

	result, err := svc.SearchRecords(context.Background(), models.SearchRecordsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if captured.Page != 1 {
		t.Errorf("expected page=1, got %d", captured.Page)
	}
	if captured.Limit != 25 {
		t.Errorf("expected limit=25, got %d", captured.Limit)
	}
}

func TestSearchRecords_Validation(t *testing.T) {
	svc := services.NewService(&mockClient{})
	tests := []struct {
		name    string
		params  models.SearchRecordsParams
		wantErr string
	}{
		{"limit too high", models.SearchRecordsParams{Limit: 200}, "Validation failed"},
		{"invalid group_by", models.SearchRecordsParams{Limit: 10, GroupBy: "invalid"}, "Validation failed"},
		{"bad from_date format", models.SearchRecordsParams{Limit: 10, FromDate: "01-01-2024", ToDate: "2024-12-31"}, "Validation failed"},
		{"bad to_date format", models.SearchRecordsParams{Limit: 10, FromDate: "2024-01-01", ToDate: "31-12-2024"}, "Validation failed"},
		{"from after to", models.SearchRecordsParams{Limit: 10, FromDate: "2024-12-01", ToDate: "2024-01-01"}, "Validation failed"},
		{"min > max", models.SearchRecordsParams{Limit: 10, MinAmount: 100, MaxAmount: 50}, "Validation failed"},
		{"invalid record_type", models.SearchRecordsParams{Limit: 10, RecordType: "savings"}, "Validation failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.SearchRecords(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestGetFinancialSummary_RequiresDates(t *testing.T) {
	svc := services.NewService(&mockClient{})
	_, err := svc.GetFinancialSummary(context.Background(), models.SummaryParams{})
	if err == nil {
		t.Fatal("expected error for missing dates")
	}
}

func TestListBudgets_Validates(t *testing.T) {
	svc := services.NewService(&mockClient{})
	tests := []struct {
		name    string
		month   int
		year    int
		wantErr string
	}{
		{"bad month", 13, 2024, "month must be between 1 and 12"},
		{"bad year", 6, 0, "year must be a positive number"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ListBudgets(context.Background(), tt.month, tt.year)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
