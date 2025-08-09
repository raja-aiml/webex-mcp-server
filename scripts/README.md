# MCP Server Scripts

## Single Unified Script: `mcp.sh`

Following KISS, DRY, and YAGNI principles, all Docker operations are handled by one script.

### Usage

```bash
./scripts/mcp.sh [command]
```

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `stdio` | Run in stdio mode | `./scripts/mcp.sh stdio` |
| `inspect` | Run with MCP Inspector | `./scripts/mcp.sh inspect` |
| `http` | Start HTTP/SSE server | `./scripts/mcp.sh http` |
| `stop` | Stop HTTP server | `./scripts/mcp.sh stop` |
| `build` | Build Docker image | `./scripts/mcp.sh build` |
| `help` | Show help | `./scripts/mcp.sh help` |

### Quick Start

```bash
# Test with MCP Inspector
./scripts/mcp.sh inspect

# Start HTTP server (background)
./scripts/mcp.sh http

# Start HTTP server (foreground)
./scripts/mcp.sh http -f

# Stop HTTP server
./scripts/mcp.sh stop
```

## Other Scripts

- `test.sh` - Test runner (existing)
- `release.sh` - Release automation (existing)

## Design Principles

- **KISS**: One script, clear commands
- **DRY**: Shared functions, no duplication
- **YAGNI**: Only essential features
- **SOLID**: Single responsibility per command