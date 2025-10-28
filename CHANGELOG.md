# Changelog

## [v0.2.0] - 2025-10-28
### Added
- Full Agent Framework: `Agent` interface, `AgentManager`, `LifecycleManager`
- Per-agent configuration (provider, model, parameters, prompt templates)
- CLI commands: `agent list`, `agent run <agent>`, `agent register`
- Logging system with per-agent logs (file + console)
- Usage tracking via `Tracker` module
- Unit tests for agent, config, usage modules

### Changed
- Refactored internal package structure for clarity
- Improved YAML config merging and validation
- Venice mock provider for offline testing

### Removed
- Phase 1 prototype agent scaffolding

---

## [v0.1.0] - 2025-10-26
### Added
- Initial CLI scaffold using Cobra
- Basic YAML configuration loader
- Minimal Venice API provider
- In-memory usage tracker
- Sample prototype agent
- Unit tests for core modules

### Changed
- None

### Removed
- None
