# spendwise-mcp

Separate MCP server for SpendWise.

## Goal

Expose SpendWise business capabilities to AI assistants through a safe, production-oriented MCP server.

## Initial Direction

- separate Git repository and deployment unit
- Go implementation to stay close to the existing backend stack
- thin adapter over the existing SpendWise backend
- read-only tools first, writes later after audit, idempotency, and confirmation infrastructure

## Planned V1

- `search_records`
- `get_record_details`
- `get_financial_summary`
- `list_categories`
- `list_budgets`
- `get_budget_progress`
- `list_goals`
- `get_goal_details`

## Backend Contract

The MCP server is expected to call the existing SpendWise backend over HTTP for now. The default backend base URL should be `http://localhost:8090/api/v1` unless overridden by environment.

## Reference Design

The detailed design blueprint currently lives in `../spendwise-be/MCP.md` and should be treated as the source architecture document until it is copied or moved into this repository.
