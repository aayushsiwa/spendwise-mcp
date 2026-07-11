package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"aayushsiwa/spendwise-mcp/models"
)

type Client interface {
	SearchRecords(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error)
	GetRecord(ctx context.Context, id string) (*models.Record, error)
	GetFinancialSummary(ctx context.Context, params models.SummaryParams) (*models.Summary, error)
	ListCategories(ctx context.Context) ([]models.Category, error)
	ListBudgets(ctx context.Context, month, year int) ([]models.Budget, error)
	GetBudgetProgress(ctx context.Context, month, year int) ([]models.BudgetProgress, error)
	ListGoals(ctx context.Context) ([]models.Goal, error)
	GetGoal(ctx context.Context, id string) (*models.Goal, error)
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *HTTPClient) SearchRecords(ctx context.Context, params models.SearchRecordsParams) (*models.SearchRecordsResult, error) {
	query := url.Values{}
	addString(query, "from", params.FromDate)
	addString(query, "to", params.ToDate)
	addString(query, "category", params.Category)
	addString(query, "type", params.RecordType)
	addString(query, "search", params.Search)
	addString(query, "groupBy", params.GroupBy)
	addFloat(query, "minAmount", params.MinAmount)
	addFloat(query, "maxAmount", params.MaxAmount)
	addInt(query, "page", params.Page)
	addInt(query, "limit", params.Limit)

	var result models.SearchRecordsResult
	if err := c.get(ctx, "/records", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *HTTPClient) GetRecord(ctx context.Context, id string) (*models.Record, error) {
	var result models.Record
	if err := c.get(ctx, "/records/"+url.PathEscape(id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *HTTPClient) GetFinancialSummary(ctx context.Context, params models.SummaryParams) (*models.Summary, error) {
	query := url.Values{}
	addString(query, "from", params.FromDate)
	addString(query, "to", params.ToDate)
	addString(query, "category", params.Category)
	addString(query, "type", params.RecordType)

	var envelope struct {
		Summary models.Summary `json:"summary"`
	}
	if err := c.get(ctx, "/summary", query, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Summary, nil
}

func (c *HTTPClient) ListCategories(ctx context.Context) ([]models.Category, error) {
	var envelope struct {
		Categories []models.Category `json:"categories"`
	}
	if err := c.get(ctx, "/categories", nil, &envelope); err != nil {
		return nil, err
	}
	return envelope.Categories, nil
}

func (c *HTTPClient) ListBudgets(ctx context.Context, month, year int) ([]models.Budget, error) {
	query := url.Values{}
	addInt(query, "month", month)
	addInt(query, "year", year)

	var envelope struct {
		Budgets []models.Budget `json:"budgets"`
	}
	if err := c.get(ctx, "/budgets", query, &envelope); err != nil {
		return nil, err
	}
	return envelope.Budgets, nil
}

func (c *HTTPClient) GetBudgetProgress(ctx context.Context, month, year int) ([]models.BudgetProgress, error) {
	query := url.Values{}
	addInt(query, "month", month)
	addInt(query, "year", year)

	var envelope struct {
		Progress []models.BudgetProgress `json:"progress"`
	}
	if err := c.get(ctx, "/budgets/progress", query, &envelope); err != nil {
		return nil, err
	}
	return envelope.Progress, nil
}

func (c *HTTPClient) ListGoals(ctx context.Context) ([]models.Goal, error) {
	var envelope struct {
		Goals []models.Goal `json:"goals"`
	}
	if err := c.get(ctx, "/goals", nil, &envelope); err != nil {
		return nil, err
	}
	return envelope.Goals, nil
}

func (c *HTTPClient) GetGoal(ctx context.Context, id string) (*models.Goal, error) {
	var envelope struct {
		Goal models.Goal `json:"goal"`
	}
	if err := c.get(ctx, "/goals/"+url.PathEscape(id), nil, &envelope); err != nil {
		return nil, err
	}
	return &envelope.Goal, nil
}

func (c *HTTPClient) get(ctx context.Context, path string, query url.Values, out any) error {
	fullURL := c.baseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Error struct {
				Type    string         `json:"type"`
				Message string         `json:"message"`
				Details map[string]any `json:"details"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Error.Message != "" {
			return fmt.Errorf("backend %s: %s", apiErr.Error.Type, apiErr.Error.Message)
		}
		return fmt.Errorf("backend request failed with status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode backend response: %w", err)
	}

	return nil
}

func addString(query url.Values, key, value string) {
	if value != "" {
		query.Set(key, value)
	}
}

func addInt(query url.Values, key string, value int) {
	if value > 0 {
		query.Set(key, strconv.Itoa(value))
	}
}

func addFloat(query url.Values, key string, value float64) {
	if value > 0 {
		query.Set(key, strconv.FormatFloat(value, 'f', -1, 64))
	}
}
