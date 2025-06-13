# Getting Started with Aphelion CLI

This guide will help you get started with the Aphelion CLI for interacting with the Aphelion Gateway.

## Prerequisites

1. Auth0 account and application configured
2. Access to an Aphelion Gateway instance

## Setup

### 1. Authenticate (Zero Configuration!)

The CLI automatically discovers all configuration from the Aphelion Gateway. Just login:

```bash
aphelion auth login
```

That's it! The CLI will:
- Auto-discover Auth0 configuration from the gateway
- Open your browser for secure authentication  
- Handle the OAuth2 PKCE flow automatically
- Securely store your tokens

Check your authentication status:

```bash
aphelion auth status
aphelion auth whoami
```

### 2. Optional: Set Custom Endpoint

If you're using a different Aphelion Gateway instance:

```bash
aphelion config set endpoint https://your-custom-gateway.com
aphelion auth login  # Will auto-discover Auth0 config from your gateway
```

## Basic Usage

### Agent Sessions

Create an agent session to execute tools:

```bash
# Create a new session
aphelion agents create

# List your sessions
aphelion agents list

# Get session details
aphelion agents describe SESSION_ID
```

### Tool Execution

Execute built-in tools:

```bash
# Echo tool
aphelion agents execute SESSION_ID --tool=echo --params='{"message":"Hello World"}'

# Get current time
aphelion agents execute SESSION_ID --tool=get_current_time

# Calculate expressions
aphelion agents execute SESSION_ID --tool=calculate --params='{"expression":"2+2*3"}'
```

### Service Management

Register your own API services:

```bash
# List available services
aphelion services list

# Register a new service
aphelion services register --spec=my-api.yaml

# List your registered services
aphelion services list --mine
```

### Memory Operations

Manage session memories:

```bash
# List memories
aphelion memory list

# Search memories
aphelion memory search "calculation"

# Get session memory
aphelion memory sessions SESSION_ID

# Create session summary
aphelion memory summarize SESSION_ID
```

### Tool Discovery

Find and discover tools:

```bash
# Search for tools
aphelion search tools "weather"

# Get popular tools
aphelion search popular

# Get personalized recommendations
aphelion search recommendations
```

### Analytics

View usage analytics:

```bash
# Your usage analytics
aphelion analytics user

# Tool usage statistics
aphelion analytics tools --user-only

# Session analytics
aphelion analytics sessions --user-only
```

## Advanced Usage

### Multiple Profiles

Manage multiple environments:

```bash
# Create profiles for different environments
aphelion config profiles create staging
aphelion config profiles create production

# Configure staging profile
aphelion config profiles switch staging
aphelion config set endpoint https://staging.aphelion.com
aphelion config set auth.domain staging.auth0.com

# Switch between profiles
aphelion config profiles switch production
aphelion config profiles switch staging
```

### Automation & Scripting

Use the CLI in scripts:

```bash
#!/bin/bash

# Check authentication
if ! aphelion auth status > /dev/null 2>&1; then
    echo "Not authenticated. Please run: aphelion auth login"
    exit 1
fi

# Create session and capture ID
SESSION_JSON=$(aphelion agents create --output json)
SESSION_ID=$(echo "$SESSION_JSON" | jq -r '.session_id')

# Execute tool and check result
if aphelion agents execute "$SESSION_ID" --tool=calculate --params='{"expression":"2+2"}' --output json | jq -e '.success'; then
    echo "Calculation successful"
else
    echo "Calculation failed"
    exit 1
fi

# Clean up
aphelion agents delete "$SESSION_ID"
```

### Output Formats

All commands support multiple output formats:

```bash
# Table format (default)
aphelion agents list

# JSON format for scripting
aphelion agents list --output json | jq '.[] | .session_id'

# YAML format
aphelion agents list --output yaml
```

## Troubleshooting

### Authentication Issues

```bash
# Check authentication status
aphelion auth status

# View detailed user info
aphelion auth whoami

# Clear and re-authenticate
aphelion auth logout
aphelion auth login
```

### Configuration Issues

```bash
# View current configuration
aphelion config list

# Reset to defaults
aphelion config reset

# Check different profiles
aphelion config profiles list
```

### Verbose Output

Enable verbose logging for debugging:

```bash
aphelion --verbose agents list
```

## Next Steps

- Explore the [API documentation](https://docs.aphelion.dev)
- Register your own services with OpenAPI specs
- Build automation workflows using the CLI
- Join the [community Discord](https://discord.gg/aphelion) for support