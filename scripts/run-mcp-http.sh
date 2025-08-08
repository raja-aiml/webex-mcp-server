#!/bin/bash
# Wrapper script for running MCP server in HTTP mode with docker-compose

# Source environment variables
source .env

# Run docker-compose with HTTP profile
docker-compose --profile http up