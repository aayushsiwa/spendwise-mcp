package backend

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"aayushsiwa/spendwise-mcp/models"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *HTTPClient) {
	t.Helper()
	srv := httptest.NewServer(handler)
	client := NewHTTPClient(srv.URL, "")
	return srv, client
}

func TestSearchRecords_Success(t *testing.T) {
	srv, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/records" {
			t.Errorf("expected /records, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.SearchRecordsResult{
			Records: []models.Record{
				{ID: "rec1", Date: "2024-01-15", Description: "test", Amount: 100, Type: "expense"},
			},
			PaginationMetadata: models.PaginationMetadata{
				Page: 1, Limit: 25, TotalCount: 1, TotalPages: 1,
			},
		})
	})
	defer srv.Close()

	result, err := client.SearchRecords(context.Background(), models.SearchRecordsParams{
		FromDate: "2024-01-01",
		ToDate:   "2024-12-31",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(result.Records))
	}
	if result.Records[0].ID != "rec1" {
		t.Errorf("expected rec1, got %s", result.Records[0].ID)
	}
	if result.TotalCount != 1 {
		t.Errorf("expected TotalCount=1, got %d", result.TotalCount)
	}
}

func TestListCategories_Success(t *testing.T) {
	srv, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"categories": []models.Category{
				{ID: "cat1", Name: "food", Icon: "🍕", Color: "#FF0000"},
			},
		})
	})
	defer srv.Close()

	categories, err := client.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(categories))
	}
	if categories[0].Name != "food" {
		t.Errorf("expected food, got %s", categories[0].Name)
	}
}

func TestBackendError_MapsToAppError(t *testing.T) {
	srv, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"type":    "not_found",
				"message": "Record not found",
			},
		})
	})
	defer srv.Close()

	_, err := client.GetRecord(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "Record not found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestAuthHeaderPropagation(t *testing.T) {
	var authHeader string
	srv, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"categories": []models.Category{},
		})
	})
	defer srv.Close()

	client.token = "test-token-123"
	_, err := client.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if authHeader != "Bearer test-token-123" {
		t.Errorf("expected Bearer test-token-123, got %s", authHeader)
	}
}
