# Claude Desktop + Webex MCP: Quick Start

## 5-Minute Setup

### 1️⃣ Build
```bash
git clone https://github.com/raja-aiml/webex-mcp-server-go.git
cd webex-mcp-server-go
make build
pwd  # Note this path!
```

### 2️⃣ Configure

**macOS**: Edit `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "webex": {
      "command": "/YOUR/PATH/HERE/build/webex-mcp-server",
      "env": {
        "WEBEX_PUBLIC_WORKSPACE_API_KEY": "your-webex-token"
      }
    }
  }
}
```

### 3️⃣ Restart Claude Desktop
- Quit Claude completely (Cmd+Q on Mac)
- Start Claude again

### 4️⃣ Test It!

In Claude, try these:

```
"List my Webex rooms"
"Show my recent messages"
"Send a message to [room name]"
```

## Common Fixes

❌ **Not working?**

1. **Path must be absolute**: `/Users/john/webex-mcp-server-go/build/webex-mcp-server` ✅
   Not: `./build/webex-mcp-server` ❌

2. **Check JSON syntax**: 
   ```bash
   # macOS/Linux
   python3 -m json.tool ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

3. **Make executable**:
   ```bash
   chmod +x /path/to/webex-mcp-server
   ```

4. **Test standalone**:
   ```bash
   /path/to/webex-mcp-server
   # Should see: "Starting webex-mcp-server v0.1.0 in stdio mode"
   ```

## What Can You Do?

Once connected, ask Claude to:

- 📝 **Messages**: "Send a Webex message to the Engineering team"
- 🏠 **Rooms**: "Create a new room called 'Project Alpha'"
- 👥 **People**: "Add sarah@company.com to the Project Alpha room"
- 📊 **Reports**: "Summarize today's messages in the Support room"
- 🔍 **Search**: "Find all messages mentioning 'deadline'"

## Get Your Webex Token

1. Go to [developer.webex.com](https://developer.webex.com)
2. Sign in
3. Click your avatar → "My Webex Apps"
4. Personal Access Token → Copy

## Example Conversation

**You**: "List my 5 most recent Webex rooms"

**Claude**: "I'll check your Webex rooms for you."
*[Uses list_rooms tool]*
"Here are your 5 most recent Webex rooms:
1. Project Alpha - Team collaboration space
2. Engineering Standup - Daily meetings
3. Customer Support - Support team channel
4. Random - Water cooler chat
5. Announcements - Company updates"

**You**: "Send a message to Project Alpha saying 'Meeting moved to 3pm'"

**Claude**: "I'll send that message to the Project Alpha room."
*[Uses create_a_message tool]*
"✓ Message sent successfully to Project Alpha: 'Meeting moved to 3pm'"

That's it! You're now controlling Webex through Claude. 🎉