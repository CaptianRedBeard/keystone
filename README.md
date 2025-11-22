# Keystone

**Keystone** is a Go-based orchestration layer that serves as a uniform entry point for AI agents.  
It provides a consistent interface for agents, a CLI for users, and a shared foundation for tracking usage, managing context, and coordinating multi-agent workflows.

## Status

Currently at Phase 3: Workflow Continuity & Ticket System

Completed:
- Agent interface and manager
- CLI for agent registration, listing, and execution
- Configuration system with YAML + CLI merge support
- Unified logging and usage tracking
- Mock and real provider integration (Venice)
- Ticket-based workflow system with per-ticket context, TTL, and hop limits
- CLI commands for ticket creation, inspection, listing, monitoring, and cleanup
- Agent interface updated to support ticket-aware handling
- Workflow continuity tested via multi-step agent CLI runs

## Features

- Register, list, and run AI agents via CLI
- Per-agent configuration for provider, model, and parameters
- Configurable provider backend (Venice and future integrations)
- Ticket-based workflow system with namespaced context
- CLI integration for ticket creation, inspection, and monitoring
- Step incrementing and TTL enforcement for agent runs
- Usage tracking with token accounting
- Unified file + console logging
- Extensible architecture ready for multi-agent orchestration

## Project Structure (Phase 3)

cmd/
    agent.go
    agent_register.go
    config.go
    print.go
    root.go
    test_helpers.go
    ticket.go
    usage.go
    workflow.go

internal/
    agent/         # Agent interfaces, registry, and manager
    config/        # Config loader, saver, and merger
    logger/        # Centralized logging utilities
    providers/     # Provider interfaces and Venice implementation
    tickets/       # Ticket struct and JSON storage backend
    usage/         # In-memory usage tracker
    workflow/      # Workflow engine and store

agents/           # Example agent YAML definitions
configs/          # Configuration templates or sample files
workflows/        # Sample workflow YAML files
main.go           # Application entry point
init_config.sh    # Initialize default configuration
keystone          # Compiled CLI binary

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