# Claude Development Instructions

## Important Git Commit Guidelines

**NEVER** include the following in git commit messages:
- 🤖 Generated with [Claude Code](https://claude.ai/code)
- Co-Authored-By: Claude <noreply@anthropic.com>
- Any other Claude attribution or emojis

Keep commit messages clean and professional.

## Environment Setup
**ALWAYS** source the `.env` file before running any commands:
```bash
source .env
```

## Build and Development Lifecycle

**ALWAYS** use `make` commands for all build and development tasks. Never run `go` commands directly.

### Common Make Commands

#### Building
- `make build` - Build the binary for current platform
- `make build all` - Build binaries for all platforms

#### Running
- `make run` - Run the application in stdio mode
- `make run http` - Run the application in HTTP/SSE mode
- `make dev` - Run in development mode with hot reload

#### Testing
- `make test` - Run all tests
- `make test-verbose` - Run tests with verbose output
- `make test-coverage` - Run tests with coverage report

#### Code Quality
- `make fmt` - Format code using gofmt
- `make lint` - Run linting checks

#### Docker
- `make docker build` - Build Docker image
- `make docker run` - Run application in Docker
- `make docker run-dev` - Run in Docker development mode
- `make docker stop` - Stop Docker containers
- `make docker clean` - Clean Docker resources

#### Other Commands
- `make clean` - Clean build artifacts
- `make install` - Build and install the binary
- `make deps` - Download and tidy dependencies
- `make health` - Check service health (when running)

## Important Notes

1. **Environment Variables**: The `.env` file contains the `WEBEX_PUBLIC_WORKSPACE_API_KEY` which is required for the application to function. Always ensure it's sourced.

2. **Testing**: When running tests, the environment variable is automatically included in make commands, but if you need to run specific tests, use:
   ```bash
   source .env && make test
   ```

3. **Development Workflow**:
   - Always start with `source .env`
   - Use `make dev` for development with hot reload
   - Use `make fmt` before committing code
   - Use `make lint` to check for issues
   - Use `make test` to ensure all tests pass

4. **Binary Location**: Built binaries are placed in the `build/` directory

5. **Cross-Platform Builds**: Use `make build all` to build for all supported platforms (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64)

## Manual Testing with MCP Inspector

For manual testing and debugging MCP tools, use the MCP Inspector:

```bash
source .env && npx @modelcontextprotocol/inspector go run .
```

This will:
- Launch the MCP Inspector web interface
- Allow you to manually test individual tools
- See real-time request/response payloads
- Debug tool interactions

The inspector will be available at the URL shown in the terminal output (typically http://localhost:3000).

## Quick Start for Development

```bash
# 1. Source environment variables
source .env

# 2. Install dependencies
make deps

# 3. Run in development mode
make dev

# 4. In another terminal, run tests
source .env && make test

# 5. Format and lint before committing
make fmt
make lint

# 6. For manual testing with MCP Inspector
source .env && npx @modelcontextprotocol/inspector go run .
```

## Debugging

If you encounter issues:
1. Ensure `.env` file exists and is sourced
2. Check that `WEBEX_PUBLIC_WORKSPACE_API_KEY` is set: `echo $WEBEX_PUBLIC_WORKSPACE_API_KEY`
3. Use `make clean` to clean build artifacts and try again
4. For verbose test output: `make test-verbose`

## Never Do

- ❌ Never run `go build` directly - use `make build`
- ❌ Never run `go test` directly - use `make test`
- ❌ Never forget to source `.env` file
- ❌ Never commit without running `make fmt` and `make lint`