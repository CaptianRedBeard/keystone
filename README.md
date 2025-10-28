# Keystone

**Keystone** is a Go-based orchestration layer that serves as a uniform entry point for AI agents.  
It provides a consistent interface for agents, a CLI for users, and a shared foundation for tracking usage, managing context, and coordinating multi-agent workflows.

## Status

Currently at Phase 2: Agent Framework

Completed:
- Agent interface and manager
- CLI for agent registration, listing, and execution
- Configuration system with YAML + CLI merge support
- Unified logging and usage tracking
- Mock and real provider integration (Venice)
- Unit tests for config, usage, and agent subsystems

## Features

- Register, list, and run AI agents via CLI
- Per-agent configuration for provider, model, and parameters
- Configurable provider backend (Venice and future integrations)
- Usage tracking with token accounting
- Unified file + console logging
- Extensible architecture ready for multi-agent orchestration

## Project Structure (Phase 1)

cmd/
    agent.go       # Agent CLI commands
    config.go      # Config CLI commands
    usage.go       # Usage CLI commands
    root.go        # Root Cobra command

internal/
    agent/         # Agent interfaces, registry, and manager
    config/        # Config loader, saver, and merger
    logger/        # Centralized logging utilities
    usage/         # In-memory usage tracker

providers/
    provider.go    # Provider interface
    venice/        # Venice provider implementation (mock + real)

main.go            # Application entry point

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

4. (Optional) View config and usage:
```bash
keystone config view
keystone usage
```

## License

MIT License