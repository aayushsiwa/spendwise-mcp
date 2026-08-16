# SpendWise MCP Server (`spendwise-mcp`)

A [Model Context Protocol](https://modelcontextprotocol.io) server that exposes [SpendWise](https://github.com/aayushsiwa/spendwise-fe) financial capabilities to AI assistants. It is a **thin, stateless adapter** over the [SpendWise backend](https://github.com/aayushsiwa/spendwise-be) REST API — it holds no data of its own and simply translates MCP tool calls into backend HTTP requests.

Related repositories:

- Frontend dashboard: <https://github.com/aayushsiwa/spendwise-fe>
- Backend API & database: <https://github.com/aayushsiwa/spendwise-be>

## Features

- **AI-native access to your finance data** — let an assistant search records, read summaries, and manage categories, budgets, and goals through your own backend, with no third-party data sharing.
- **Full CRUD coverage** — records, categories, budgets, and savings goals can all be queried and mutated via tools.
- **Read-only and write tools** — beyond reading, the server can create/update/delete records, categories, budgets, and goals, and add progress to goals.
- **Rich querying** — `search_records` supports date ranges, category/type filters, amount bounds, free-text search, grouping (by category or month), and pagination.
- **Financial summaries** — `get_financial_summary` returns income, expense, net, balances, and category breakdowns for a date range.
- **Budget progress** — `get_budget_progress` reports spent-vs-budgeted per category for a month.
- **Audit context on every call** — each tool invocation carries `X-MCP-Client`, `X-MCP-Actor`, `X-Request-ID`, and (optionally) a bearer token, so backend activity can be attributed.
- **Transport via stdio** — runs as a local subprocess spawned by the MCP client; no separate network listener.
- **Consistent results** — every tool returns a `{ "ok", "request_id", "data" | "error" }` envelope.

## Architecture

```text
AI assistant (client)
      ↓  stdio (MCP)
spendwise-mcp  ──→  HTTP  ──→  spendwise-be  ──→  SQLite / Postgres / MySQL
```

- Built on [`mark3labs/mcp-go`](https://pkg.go.dev/github.com/mark3labs/mcp-go), served over stdio.
- **Stateless adapter:** the MCP server never stores financial data; it proxies the backend on each call.
- **Session context:** `routes.AttachRoutes` wraps every tool call with a session carrying `ActorID`, `ClientName`, `BackendToken`, and a generated `RequestID`, injected into outgoing backend requests as headers.

## Project structure

- **`main.go`** — entrypoint: loads config, builds the backend HTTP client and service, registers tools, and starts the stdio server.
- **`routes/`** — declarative tool definitions (`NewRoutes`) and `AttachRoutes`, which binds each tool to its handler with session context.
- **`handlers/`** — one file per resource (`records.go`, `categories.go`, `budgets.go`, `goals.go`, `handlers.go`); parse MCP arguments, call the service, and shape results.
- **`services/`** — orchestrates backend calls and maps backend responses to MCP results.
- **`backend/`** — HTTP client to the SpendWise REST API (the only component that talks to the backend).
- **`session/`** — request-scoped session (actor, client, token, request ID) and header injection.
- **`config/`** — environment-based configuration.
- **`errors/`** — typed application errors mapped into the result envelope.
- **`models/`** — shared request/response shapes.

## Tools reference

### Records
- `search_records` — search with filters, pagination, and optional grouping.
- `get_record_details` — fetch one record by ID.
- `create_spending_record` — add an income, expense, or transfer.
- `update_spending_record` — partial update by ID.
- `delete_spending_record` — delete by ID.

### Categories
- `list_categories` — list all categories.
- `create_category` — add a category (name, icon, color).
- `update_category` — update a category.
- `delete_category` — delete by ID.

### Budgets
- `list_budgets` — budgets for a month/year.
- `get_budget_progress` — spent vs. budget per category.
- `create_budget` — set a budget for a category/month/year.
- `update_budget` — update a budget amount.
- `delete_budget` — delete by ID.

### Goals
- `list_goals` — all savings goals.
- `get_goal_details` — one goal by ID.
- `create_goal` — set a new savings target.
- `update_goal` — partial update by ID.
- `delete_goal` — delete by ID.
- `add_goal_progress` — contribute toward a goal.

### Summary
- `get_financial_summary` — income, expense, net, balances, and category breakdown for a date range.

## Getting started

### Prerequisites

- A running [SpendWise backend](https://github.com/aayushsiwa/spendwise-be) (default `http://localhost:8090/api/v1`).
- Go 1.26+ (toolchain per `go.mod`).

### Option 1 — Pre-built binary

Download the latest release for your platform from the [releases page](https://github.com/aayushsiwa/spendwise-mcp/releases), make it executable (`chmod +x spendwise-mcp-*`), and point your MCP client at it.

### Option 2 — Build from source

```sh
cp .env.example .env   # fill in SPENDWISE_BACKEND_BASE_URL (and token for production)
go build -o spendwise-mcp .
./spendwise-mcp
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

If the backend requires a token:

```jsonc
"env": {
  "SPENDWISE_BACKEND_BASE_URL": "http://localhost:8090/api/v1",
  "SPENDWISE_BACKEND_TOKEN": "your-token"
}
```

## Environment

| Variable | Default | Description |
|---|---|---|
| `SPENDWISE_BACKEND_BASE_URL` | `http://localhost:8090/api/v1` | SpendWise REST API base. |
| `SPENDWISE_BACKEND_TOKEN` | _(empty)_ | Bearer token forwarded to the backend (requires backend auth to be meaningful). |
| `SPENDWISE_ACTOR_ID` | `mcp-local` | Actor identifier included in audit headers. |
| `SPENDWISE_CLIENT_NAME` | `spendwise-mcp` | Client name included in audit headers. |
| `MCP_SERVER_NAME` | `SpendWise MCP` | MCP server name. |
| `MCP_SERVER_VERSION` | `0.1.0` | MCP server version. |

## Notes

- **Auth is end-to-end.** The `SPENDWISE_BACKEND_TOKEN` is forwarded as a bearer token; it is only enforced if the backend actually validates it.
