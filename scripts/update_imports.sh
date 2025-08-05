#!/bin/bash

# Update all import paths from old structure to new internal structure
find . -name "*.go" -type f -exec sed -i '' \
  -e 's|"github.com/raja-aiml/webex-mcp-server-go/app"|"github.com/raja-aiml/webex-mcp-server-go/internal/app"|g' \
  -e 's|"github.com/raja-aiml/webex-mcp-server-go/config"|"github.com/raja-aiml/webex-mcp-server-go/internal/config"|g' \
  -e 's|"github.com/raja-aiml/webex-mcp-server-go/handlers"|"github.com/raja-aiml/webex-mcp-server-go/internal/handlers"|g' \
  -e 's|"github.com/raja-aiml/webex-mcp-server-go/server"|"github.com/raja-aiml/webex-mcp-server-go/internal/server"|g' \
  -e 's|"github.com/raja-aiml/webex-mcp-server-go/tools"|"github.com/raja-aiml/webex-mcp-server-go/internal/tools"|g' \
  -e 's|"github.com/raja-aiml/webex-mcp-server-go/webex"|"github.com/raja-aiml/webex-mcp-server-go/internal/webex"|g' \
  {} \;

echo "Import paths updated successfully"