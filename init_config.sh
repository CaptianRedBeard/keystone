#!/usr/bin/env bash
# init_config.sh - Initialize default Keystone config

set -euo pipefail

CONFIG_DIR="$HOME/.keystone"
CONFIG_FILE="$CONFIG_DIR/config.yaml"

echo "Initializing Keystone config..."

mkdir -p "$CONFIG_DIR"

DB_FILE="$CONFIG_DIR/usage.db"
AGENTS_DIR="$CONFIG_DIR/agents"

mkdir -p "$AGENTS_DIR"

cat > "$CONFIG_FILE" <<YAML
# Keystone Configuration

db_path: "$DB_FILE"
agents_dir: "$AGENTS_DIR"

# Whether CLI output should default to JSON.
json_output: false

# API keys or other secrets (empty by default).
secrets: {}
YAML

echo "Default config created at $CONFIG_FILE"
