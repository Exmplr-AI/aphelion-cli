# Aphelion CLI

A command-line interface for interacting with the Aphelion platform, providing tools for authentication, service management, search, analytics, and more.

## Installation

### Using Homebrew (Recommended)

```bash
brew install Exmplr-AI/tap/aphelion
```

### Using Go

```bash
go install github.com/Exmplr-AI/aphelion-cli@latest
```

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/Exmplr-AI/aphelion-cli.git
   cd aphelion-cli
   ```

2. Build the CLI:
   ```bash
   make build
   ```

3. Install the binary:
   ```bash
   make install
   ```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/Exmplr-AI/aphelion-cli/releases).

### Using Make Targets

- `make build` - Build the binary for current platform
- `make build-all` - Build for all supported platforms
- `make test` - Run tests
- `make lint` - Run linter
- `make fmt` - Format code
- `make clean` - Clean build artifacts
- `make install` - Install binary to /usr/local/bin
- `make dev` - Run full development workflow

## Usage

The CLI provides several commands:

- `aphelion auth` - Authentication commands
- `aphelion search` - Search functionality
- `aphelion analytics` - Analytics tools
- `aphelion memory` - Memory management
- `aphelion services` - Service management
- `aphelion agents` - Agent operations
- `aphelion config` - Configuration management
- `aphelion version` - Version information

For detailed help on any command, use:
```bash
aphelion [command] --help
```

## Development

### Prerequisites

- Go 1.21 or later
- Make

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Contributing

1. Make your changes
2. Run `make dev` to ensure code quality
3. Submit a pull request

## Examples

See the `examples/` directory for usage examples and sample configurations.