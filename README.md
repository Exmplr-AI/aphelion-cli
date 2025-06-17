# Aphelion CLI

A command-line interface for the Aphelion Gateway platform - a unified AI agent platform with tool discovery and memory capabilities.

## Installation

### Homebrew (macOS/Linux) - Recommended

```bash
# Add the tap and install
brew tap exmplr-ai/aphelion
brew install aphelion
```

### Download Binary

Download the latest release from the [releases page](https://github.com/Exmplr-AI/aphelion-cli/releases).

### Build from Source

```bash
git clone https://github.com/Exmplr-AI/aphelion-cli.git
cd aphelion-cli
go build -o aphelion main.go
```

### Install with Go

```bash
go install github.com/Exmplr-AI/aphelion-cli@latest
```

## Quick Start

1. **Login to your account**:
   ```bash
   aphelion auth login
   ```

2. **Initialize a new agent project**:
   ```bash
   aphelion agent init
   ```

3. **Run an agent with cron scheduling**:
   ```bash
   aphelion agent run ./agent.py --cron "*/10 * * * *"
   ```

4. **List available services**:
   ```bash
   aphelion registry list
   ```

5. **Search your memories**:
   ```bash
   aphelion memory search "calculation"
   ```

## Commands

### Authentication

| Command | Description |
|---------|-------------|
| `aphelion auth login` | Login with username and password |
| `aphelion auth register` | Register a new account |
| `aphelion auth profile` | Show user profile information |
| `aphelion auth logout` | Logout and clear credentials |
| `aphelion auth oauth` | Get Auth0 OAuth configuration |

### Agent Development

| Command | Description |
|---------|-------------|
| `aphelion agent init` | Initialize a new agent project with scaffolding |
| `aphelion agent run [file]` | Run an agent with optional cron scheduling |
| `aphelion agent run --cron "*/10 * * * *"` | Schedule agent with cron expression |
| `aphelion agent run --daemon` | Run agent as daemon process |

### Service Registry

| Command | Description |
|---------|-------------|
| `aphelion registry list` | List all public services |
| `aphelion registry my-services` | List your registered services |
| `aphelion registry create` | Register a new API service |
| `aphelion registry add-openapi --file [spec]` | Register service from OpenAPI specification |
| `aphelion registry get [ID]` | Get service details |
| `aphelion registry delete [ID]` | Delete a service |

### Tool Discovery & Execution

| Command | Description |
|---------|-------------|
| `aphelion tools describe [tool-name]` | Show tool parameters, schema, and examples |
| `aphelion tools try --tool [name] --params '[json]'` | Execute a tool with given parameters |
| `aphelion tools try --dry-run` | Validate parameters without execution |

### Memory Management

| Command | Description |
|---------|-------------|
| `aphelion memory list` | List your memories with pagination |
| `aphelion memory search [QUERY]` | Search memories using semantic similarity |
| `aphelion memory stats` | Show memory usage statistics |
| `aphelion memory clear` | Delete all memories (with confirmation) |
| `aphelion memory clear --session [ID]` | Delete memories for specific session |

### Analytics

| Command | Description |
|---------|-------------|
| `aphelion analytics user` | Show user-specific analytics |
| `aphelion analytics tools` | Show tool usage analytics |
| `aphelion analytics sessions` | Show session analytics |

### Utility Commands

| Command | Description |
|---------|-------------|
| `aphelion version` | Show version information |
| `aphelion completion [shell]` | Generate shell completion script |

## Agent Development

### Initializing an Agent

```bash
# Create a new agent project
aphelion agent init

# This creates:
# - agent.py (main agent script)
# - requirements.txt (Python dependencies)
# - .aphelion/config.yaml (agent configuration)
# - .aphelion/session (session management)
```

### Agent Structure

The generated `agent.py` includes:

- **Gateway Integration**: Built-in API client for Aphelion Gateway
- **Session Management**: Automatic session creation and persistence
- **Memory Checkpointing**: Configurable memory saving intervals
- **Tool Discovery**: Search and execute tools via Gateway API
- **Error Handling**: Robust error handling and logging

### Running Agents

```bash
# Run agent once
aphelion agent run ./agent.py

# Schedule with cron (every 10 minutes)
aphelion agent run ./agent.py --cron "*/10 * * * *"

# Run as daemon
aphelion agent run ./agent.py --daemon

# Verbose output
aphelion agent run ./agent.py --verbose
```

### Agent Configuration

Edit `.aphelion/config.yaml` to customize:

```yaml
# Aphelion Agent Configuration
name: "my-agent"
description: "A sample Aphelion agent"
version: "1.0.0"

# Gateway configuration
gateway:
  api_url: "https://api.aphelion.exmplr.ai"
  
# Agent execution settings
execution:
  memory_checkpoint_interval: "10m"
  max_memory_entries: 1000
  
# Logging configuration
logging:
  level: "info"
  file: "agent.log"
```

## Tool Development

### Describing Tools

```bash
# Get tool details, parameters, and examples
aphelion tools describe exmplr_core.search

# Output includes:
# - Tool description and version
# - Parameter schema with types
# - Required vs optional parameters
# - Usage examples
```

### Testing Tools

```bash
# Execute a tool with parameters
aphelion tools try --tool exmplr_core.search --params '{"q": "Multiple Sclerosis"}'

# Validate parameters without execution
aphelion tools try --tool exmplr_core.search --params '{"q": "test"}' --dry-run

# Verbose output with metadata
aphelion tools try --tool exmplr_core.search --params '{"q": "test"}' --verbose
```

## Service Registration

### From OpenAPI Specification

```bash
# Register service from OpenAPI spec
aphelion registry add-openapi --file ./openapi.json

# Override service details
aphelion registry add-openapi --file ./openapi.json \
  --name "Custom API" \
  --description "My custom API service" \
  --base-url "https://api.example.com"
```

The CLI automatically:
- Parses OpenAPI specification
- Generates STELLA manifest
- Converts endpoints to tools
- Registers with Aphelion Gateway

### Manual Registration

```bash
# Register with OpenAPI spec file
aphelion registry create --name "Weather API" \
  --description "Weather service" \
  --spec-file openapi.json

# Register without spec
aphelion registry create --name "My API" \
  --description "Custom API service"
```

## Configuration

The CLI stores configuration in `~/.aphelion/config.yaml`. This includes:

- API endpoint URL
- Authentication tokens
- User preferences
- Output format settings

### Environment Variables

You can override configuration with environment variables:

- `APHELION_API_URL`: API base URL (default: https://api.aphelion.exmplr.ai)
- `APHELION_OUTPUT`: Output format (json, yaml, table)
- `APHELION_VERBOSE`: Enable verbose output

## Global Flags

- `--config`: Specify config file path
- `--api-url`: Override API base URL
- `--output, -o`: Output format (json, yaml, table)
- `--verbose, -v`: Enable verbose output

## Examples

### Agent Workflow

```bash
# 1. Initialize agent project
aphelion agent init

# 2. Customize agent.py for your use case
# 3. Install dependencies
pip install -r requirements.txt

# 4. Test run
aphelion agent run ./agent.py --verbose

# 5. Schedule for production
aphelion agent run ./agent.py --cron "*/10 * * * *"
```

### Tool Discovery & Testing

```bash
# Find tools related to research
aphelion tools describe exmplr_core.meta_analysis

# Test tool execution
aphelion tools try --tool exmplr_core.search \
  --params '{"q": "Multiple Sclerosis", "limit": 10}'

# Validate parameters only
aphelion tools try --tool exmplr_core.search \
  --params '{"q": "test"}' --dry-run
```

### Service Registration

```bash
# Register from OpenAPI spec
aphelion registry add-openapi --file ./medical-api.json

# List your services
aphelion registry my-services

# Get service details and manifest
aphelion registry get service-123
```

### Memory Operations

```bash
# Search memories with threshold
aphelion memory search "Multiple Sclerosis research" --threshold 0.8

# Clear memories for specific session
aphelion memory clear --session abc123

# Get memory usage statistics
aphelion memory stats
```

### Output Formats

```bash
# Table format (default, human-readable)
aphelion registry list

# JSON format (for scripting)
aphelion registry list --output json

# YAML format (for configuration)
aphelion registry list --output yaml
```

## Advanced Usage

### Cron Expressions

The agent runner supports standard cron expressions:

```bash
# Every 10 minutes
aphelion agent run agent.py --cron "*/10 * * * *"

# Every hour at minute 0
aphelion agent run agent.py --cron "0 * * * *"

# Every day at 2:30 AM
aphelion agent run agent.py --cron "30 2 * * *"

# Every Monday at 9 AM
aphelion agent run agent.py --cron "0 9 * * MON"
```

### Multi-Language Agent Support

The CLI supports agents in multiple languages:

```bash
# Python agents
aphelion agent run agent.py

# Node.js agents
aphelion agent run agent.js

# Go agents
aphelion agent run agent.go

# Any executable
aphelion agent run ./custom-agent
```

### Output Piping & Scripting

```bash
# Get service IDs for scripting
aphelion registry list --output json | jq '.[] | .id'

# Count total memories
aphelion memory stats --output json | jq '.total_memories'

# Export memories to file
aphelion memory list --output yaml > memories.yaml
```

## Shell Completion

Generate completion scripts for your shell:

```bash
# Bash
aphelion completion bash > /etc/bash_completion.d/aphelion

# Zsh
aphelion completion zsh > "${fpath[1]}/_aphelion"

# Fish
aphelion completion fish > ~/.config/fish/completions/aphelion.fish

# PowerShell
aphelion completion powershell > aphelion.ps1
```

## Development

### Prerequisites

- Go 1.21 or later
- Access to Aphelion Gateway API

### Building

```bash
# Build for current platform
go build -o aphelion main.go

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o aphelion-linux main.go
GOOS=windows GOARCH=amd64 go build -o aphelion-windows.exe main.go
GOOS=darwin GOARCH=amd64 go build -o aphelion-macos main.go
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Tidy dependencies
go mod tidy
```

### Project Structure

```
aphelion-cli/
├── cmd/                    # Command implementations
│   ├── agent/             # Agent management commands
│   ├── analytics/         # Analytics commands
│   ├── auth/              # Authentication commands
│   ├── memory/            # Memory management commands
│   ├── registry/          # Service registry commands
│   ├── tools/             # Tool discovery commands
│   └── root.go            # Root command and CLI setup
├── internal/
│   └── utils/             # Shared utilities
├── pkg/
│   ├── api/               # API client and types
│   ├── auth/              # Authentication logic
│   └── config/            # Configuration management
├── main.go                # Entry point
├── go.mod                 # Go module definition
└── README.md              # This file
```

## API Documentation

The CLI interacts with the Aphelion Gateway API. For full API documentation, visit:
- API Base URL: https://api.aphelion.exmplr.ai
- Documentation: https://api.aphelion.exmplr.ai/docs

### Key Endpoints

- `POST /auth/login` - User authentication
- `GET /services` - List services
- `GET /search/tools` - Tool discovery
- `POST /tools/{name}/execute` - Tool execution
- `POST /memory` - Save memories
- `GET /memory/search` - Search memories
- `POST /sessions` - Create agent sessions

## Troubleshooting

### Common Issues

**Authentication Errors**
```bash
# Check authentication status
aphelion auth profile

# Re-authenticate
aphelion auth logout
aphelion auth login
```

**Agent Execution Issues**
```bash
# Run with verbose output
aphelion agent run agent.py --verbose

# Check agent configuration
cat .aphelion/config.yaml
```

**Tool Execution Failures**
```bash
# Validate parameters first
aphelion tools try --tool tool-name --params '{}' --dry-run

# Get tool documentation
aphelion tools describe tool-name
```

### Debug Mode

Enable verbose output for troubleshooting:

```bash
# Global verbose flag
aphelion --verbose [command]

# Command-specific verbose
aphelion agent run agent.py --verbose
aphelion tools try --tool test --params '{}' --verbose
```

## Support

- **Issues**: Report bugs and feature requests on [GitHub Issues](https://github.com/exmplr/aphelion-cli/issues)
- **Documentation**: Full documentation at [docs.aphelion.exmplr.ai](https://docs.aphelion.exmplr.ai)
- **Community**: Join discussions on [GitHub Discussions](https://github.com/exmplr/aphelion-cli/discussions)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### Latest Version

#### New Features
- **Agent Development**: Complete agent lifecycle management with `aphelion agent init` and `aphelion agent run`
- **Tool Discovery**: Comprehensive tool exploration with `aphelion tools describe` and `aphelion tools try`
- **OpenAPI Integration**: Automatic service registration from OpenAPI specs with `aphelion registry add-openapi`
- **Enhanced Memory Management**: Session-specific memory operations with `--session` flag
- **Cron Scheduling**: Built-in cron scheduler for automated agent execution
- **Multi-Language Support**: Support for Python, Node.js, Go, and custom executable agents
- **Dry-Run Mode**: Parameter validation without execution for safe testing

#### Improvements
- Enhanced output formatting with consistent table, JSON, and YAML support
- Better error handling and user feedback
- Comprehensive help documentation
- Shell completion support
- Verbose mode for debugging and troubleshooting