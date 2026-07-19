# spendwise-mcp

MCP server that exposes [SpendWise](https://github.com/aayushsiwa/spendwise-be) financial capabilities to AI assistants.

## Prerequisites

- [SpendWise backend](https://github.com/aayushsiwa/spendwise-be) running and accessible
- Go 1.26.0 or later (required to build from source)

## Setup

### Option 1 — Pre-built binary

1. Download the latest release for your platform from [releases](https://github.com/aayushsiwa/spendwise-mcp/releases).
2. Make it executable: `chmod +x spendwise-mcp-*`
3. Configure your MCP client (see below).

### Option 2 — Build from source

```bash
git clone https://github.com/aayushsiwa/spendwise-mcp
cd spendwise-mcp
go build -o spendwise-mcp .
```

## Configure your MCP client

Add to `~/.config/opencode/opencode.json`, `~/.cursor/mcp.json`, or `claude_desktop_config.json`:

```jsonc
{
  "mcpServers": {
    "spendwise": {
      "command": "/path/to/spendwise-mcp",
      "env": {
        "SPENDWISE_BACKEND_BASE_URL": "http://localhost:8090/api/v1"
      }
    }
  }
}
```

If the backend requires a token, replace the `env` block with:

```jsonc
"env": {
  "SPENDWISE_BACKEND_BASE_URL": "http://localhost:8090/api/v1",
  "SPENDWISE_BACKEND_TOKEN": "your-token"
}
```

## Tools

### Records
- `search_records` — search with filters, pagination, grouping
- `get_record_details` — get one record by ID
- `create_spending_record` — add income, expense, or transfer
- `update_spending_record` — partial update by ID
- `delete_spending_record` — delete by ID

### Categories
- `list_categories` — list all categories
- `create_category` — add a new category
- `update_category` — update an existing category
- `delete_category` — delete by ID

### Budgets
- `list_budgets` — budgets for a month/year
- `get_budget_progress` — spent vs budget per category
- `create_budget` — set a budget for a category/month/year
- `update_budget` — update budget amount
- `delete_budget` — delete by ID

### Goals
- `list_goals` — all savings goals
- `get_goal_details` — one goal by ID
- `create_goal` — set a new savings target
- `update_goal` — partial update by ID
- `delete_goal` — delete by ID
- `add_goal_progress` — contribute toward a goal

### Summary
- `get_financial_summary` — income, expense, net, balances, and category breakdown for a date range

## Environment

| Variable | Default | Description |
|---|---|---|
| `SPENDWISE_BACKEND_BASE_URL` | `http://localhost:8090/api/v1` | SpendWise REST API base |
| `SPENDWISE_BACKEND_TOKEN` | — | Bearer token for backend auth |
| `SPENDWISE_ACTOR_ID` | `mcp-local` | Actor identifier for audit |
| `SPENDWISE_CLIENT_NAME` | `spendwise-mcp` | Client name for audit |
| `MCP_SERVER_NAME` | `SpendWise MCP` | MCP server name |
| `MCP_SERVER_VERSION` | `0.1.0` | MCP server version |
