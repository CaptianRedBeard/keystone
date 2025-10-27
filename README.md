# Keystone

**Keystone** is a Go-based orchestration layer that serves as a uniform entry point for AI agents.  
It provides a consistent interface for agents, a CLI for users, and a shared foundation for tracking usage, managing context, and coordinating multi-agent workflows.

## Status

Work in progress – currently at **Phase 1: Core Infrastructure**.  

Completed in Phase 1:
- CLI scaffold using Cobra
- Configuration management (YAML + environment overrides)
- Minimal Venice API provider
- In-memory usage tracker
- Prototype sample agent
- Unit tests for core modules

## Features

- Run and manage agents via CLI (`keystone agent run <agent>`)
- View usage statistics (`keystone usage`)
- View and edit configuration (`keystone config view`)
- Fully testable offline with mocked provider responses

## Project Structure (Phase 1)

cmd/
agent.go # CLI agent commands
config.go # CLI config commands
usage.go # CLI usage commands
root.go # Root Cobra command

internal/
agents/
sample_agent.go # Prototype agent
agent.go # Agent struct and manager
config/
config.go # Config loader, saver, printer
providers/
provider.go # AIProvider interface
venice/
venice.go # VeniceProvider implementation
usage/
usage.go # In-memory usage tracker
main.go # Entry point

## Getting Started

1. Install Go (1.24+ recommended)
2. Clone the repository:

```bash
git clone https://github.com/CaptianRedBeard/keystone.git
cd keystone
```

3. Run the CLI:
```bash
go run . --help
keystone agent list
keystone agent run sample_agent "Hello world"
```

## License

MIT License