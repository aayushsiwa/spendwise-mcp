package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	apperrors "aayushsiwa/spendwise-mcp/errors"
	"aayushsiwa/spendwise-mcp/services"
	"aayushsiwa/spendwise-mcp/session"

	"github.com/mark3labs/mcp-go/mcp"
)

type ToolHandler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)

type Handler struct {
	Service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{Service: service}
}

func successResult(ctx context.Context, v any) (*mcp.CallToolResult, error) {
	requestID := ""
	if sessionCtx := session.FromContext(ctx); sessionCtx != nil {
		requestID = sessionCtx.RequestID
	}
	return marshalToolResult(map[string]any{
		"ok":         true,
		"request_id": requestID,
		"data":       v,
	}, false)
}

func errorResult(ctx context.Context, err error) (*mcp.CallToolResult, error) {
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		appErr = apperrors.NewInternal("An unexpected error occurred", err)
	}

	requestID := ""
	if sessionCtx := session.FromContext(ctx); sessionCtx != nil {
		requestID = sessionCtx.RequestID
	}

	return marshalToolResult(map[string]any{
		"ok":         false,
		"request_id": requestID,
		"error": map[string]any{
			"type":    appErr.Type,
			"message": appErr.Message,
			"details": appErr.Details,
		},
	}, true)
}

func marshalToolResult(v any, isError bool) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	if isError {
		return mcp.NewToolResultError(string(data)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func requireString(req mcp.CallToolRequest, key string) (string, error) {
	value, ok := req.GetArguments()[key]
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}
	str, ok := value.(string)
	if !ok || strings.TrimSpace(str) == "" {
		return "", fmt.Errorf("%s must be a non-empty string", key)
	}
	return str, nil
}

func optionalString(req mcp.CallToolRequest, key string) (string, error) {
	value, ok := req.GetArguments()[key]
	if !ok || value == nil {
		return "", nil
	}
	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}
	return str, nil
}

func optionalInt(req mcp.CallToolRequest, key string) (int, error) {
	value, ok := req.GetArguments()[key]
	if !ok || value == nil {
		return 0, nil
	}
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("%s must be a number", key)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("%s must be a number", key)
	}
}

func optionalFloat(req mcp.CallToolRequest, key string) (float64, error) {
	value, ok := req.GetArguments()[key]
	if !ok || value == nil {
		return 0, nil
	}
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be a number", key)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("%s must be a number", key)
	}
}
