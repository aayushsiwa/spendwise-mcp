package session

import (
	"context"

	"github.com/google/uuid"
)

type Context struct {
	RequestID    string
	ActorID      string
	ClientName   string
	BackendToken string
	ToolName     string
}

type contextKey struct{}

func New(actorID, clientName, backendToken, toolName string) *Context {
	return &Context{
		RequestID:    uuid.NewString(),
		ActorID:      actorID,
		ClientName:   clientName,
		BackendToken: backendToken,
		ToolName:     toolName,
	}
}

func WithContext(ctx context.Context, session *Context) context.Context {
	return context.WithValue(ctx, contextKey{}, session)
}

func FromContext(ctx context.Context) *Context {
	value, _ := ctx.Value(contextKey{}).(*Context)
	return value
}
