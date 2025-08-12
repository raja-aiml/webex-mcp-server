# Configuration Guide

## Environment Variables

### Required
- `WEBEX_PUBLIC_WORKSPACE_API_KEY` - Your Webex API access token

### Optional
- `MCP_SERVER_PORT` - Port for HTTP/SSE mode (default: 3000)
- `LOG_LEVEL` - Logging level: debug, info, warn, error (default: info)

## Configuration File

The server can be configured using a `config.json` file:

```json
{
  "server": {
    "mode": "stdio",
    "port": 3000
  },
  "webex": {
    "api_base_url": "https://webexapis.com/v1"
  },
  "logging": {
    "level": "info",
    "format": "json"
  }
}
```

## Running Modes

### STDIO Mode (Default)
```bash
source .env && make run
```

### HTTP/SSE Mode
```bash
source .env && make run http
```

### Development Mode
```bash
source .env && make dev
```

## Docker Configuration

### Environment File
Create a `.env` file with your Webex API key:
```bash
WEBEX_PUBLIC_WORKSPACE_API_KEY=your_token_here
```

### Docker Compose
```yaml
version: '3.8'
services:
  webex-mcp:
    image: webex-mcp-server:latest
    environment:
      - WEBEX_PUBLIC_WORKSPACE_API_KEY=${WEBEX_PUBLIC_WORKSPACE_API_KEY}
    ports:
      - "3000:3000"
```

## MCP Client Configuration

### Claude Desktop
Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "webex": {
      "command": "/path/to/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your_token_here"
      }
    }
  }
}
```

### Claude Code CLI
Configure Claude Code to use the Webex MCP server:

1. **Global Configuration** (`~/.config/claude-code/settings.json`):
```json
{
  "mcpServers": {
    "webex": {
      "command": "/usr/local/bin/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
      }
    }
  }
}
```

2. **Project-specific Configuration** (`.claude-code/settings.json`):
```json
{
  "mcpServers": {
    "webex": {
      "command": "./build/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
      }
    }
  }
}
```

3. **Environment Variables**:
```bash
export CLAUDE_CODE_MCP_SERVERS='{"webex":{"command":"webex-mcp-server","env":{"WEBEX_PUBLIC_WORKSPACE_API_KEY":"your_token"}}}'
```

### VS Code Extensions

#### 1. VS Code MCP Extension
Configure in `.mcp.json`:
```json
{
  "mcpServers": {
    "webex": {
      "command": "go",
      "args": ["run", "."],
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
      }
    }
  }
}
```

#### 2. Claude for VS Code Extension
Add to VS Code settings (`settings.json`):
```json
{
  "claude.mcpServers": {
    "webex": {
      "command": "${workspaceFolder}/build/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
      }
    }
  },
  "claude.apiKey": "your-claude-api-key",
  "claude.model": "claude-3-opus-20240229"
}
```

#### 3. Continue.dev Extension
Configure in `.continue/config.json`:
```json
{
  "models": [{
    "title": "Claude with Webex MCP",
    "provider": "anthropic",
    "model": "claude-3-opus-20240229",
    "apiKey": "your-claude-api-key",
    "contextProviders": [{
      "name": "mcp",
      "params": {
        "servers": {
          "webex": {
            "command": "webex-mcp-server",
            "env": {
              "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
            }
          }
        }
      }
    }]
  }]
}
```

### GitHub Copilot Configuration

While GitHub Copilot doesn't directly support MCP servers, you can integrate it with Claude Code:

#### 1. Copilot Chat with External Tools
Configure `.github/copilot-instructions.md`:
```markdown
When working with Webex APIs, use the following MCP server tools:
- list_rooms: Get all Webex rooms
- list_messages: Get messages from a room
- create_a_message: Send a message to a room
- get_my_own_details: Get current user details
- list_webhooks: List all webhooks
- create_a_webhook: Create a new webhook

The Webex MCP server is available at ./build/webex-mcp-server
```

#### 2. Custom Copilot Aliases
Add to your shell configuration:
```bash
# ~/.zshrc or ~/.bashrc
alias copilot-webex='gh copilot suggest "using webex-mcp-server at ./build/webex-mcp-server"'
alias copilot-mcp='gh copilot explain "how to use MCP server with $(cat .mcp.json)"'
```

#### 3. Copilot Workspace Integration
Create `.github/copilot/config.yml`:
```yaml
version: 1
tools:
  webex-mcp:
    description: "Webex MCP Server for API interactions"
    command: "./build/webex-mcp-server"
    environment:
      WEBEX_PUBLIC_WORKSPACE_API_KEY: "${secrets.WEBEX_API_KEY}"
```

### IDE Integration Best Practices

#### Development Workflow
1. **VS Code + Claude Code CLI**:
   ```bash
   # Start Claude Code in VS Code terminal
   claude-code --mcp-server webex:./build/webex-mcp-server
   ```

2. **Multi-IDE Setup**:
   - Use `.mcp.json` for VS Code MCP extension
   - Use `CLAUDE.md` for Claude Code instructions
   - Use `.github/copilot-instructions.md` for Copilot context

3. **Debugging Configuration**:
   ```json
   {
     "version": "0.2.0",
     "configurations": [{
       "name": "Debug Webex MCP",
       "type": "go",
       "request": "launch",
       "mode": "debug",
       "program": "${workspaceFolder}",
       "env": {
         "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}",
         "LOG_LEVEL": "debug"
       }
     }]
   }
   ```

### Tool Synchronization

To ensure all AI tools work together:

1. **Shared Environment File** (`.env`):
```bash
WEBEX_PUBLIC_WORKSPACE_API_KEY=your_token_here
CLAUDE_API_KEY=your_claude_key
OPENAI_API_KEY=your_openai_key
```

2. **Unified MCP Configuration** (`.mcp-unified.json`):
```json
{
  "version": "1.0",
  "servers": {
    "webex": {
      "command": "./build/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "${env:WEBEX_PUBLIC_WORKSPACE_API_KEY}"
      },
      "capabilities": ["tools", "completions", "logging"]
    }
  },
  "clients": {
    "claude-code": true,
    "vscode-mcp": true,
    "continue-dev": true
  }
}
```

## Security Configuration

### API Key Management
- Never commit API keys to version control
- Use environment variables or secure vaults
- Rotate keys regularly
- Use read-only tokens when possible

### Network Security
- Use HTTPS for all API communications
- Implement rate limiting
- Enable request logging for audit trails

## Troubleshooting

### Debug Mode
Enable debug logging:
```bash
LOG_LEVEL=debug make run
```

### Common Issues
1. **Missing API Key**: Ensure `.env` file is sourced
2. **Connection Errors**: Check network connectivity to Webex APIs
3. **Permission Errors**: Verify API key has required scopes
4. **Build Failures**: Run `make clean` and rebuild

## Performance Tuning

### Caching
- Response caching is enabled by default
- Cache TTL: 5 minutes for room lists, 1 minute for messages

### Connection Pooling
- Max idle connections: 10
- Max connections per host: 2
- Connection timeout: 30 seconds