#!/bin/bash
# Claude Code Stop Hook - sends notification to menubar app

# Extract project directory name (fallback to pwd if CLAUDE_PROJECT_DIR unset)
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
PROJECT_NAME="$(basename "$PROJECT_DIR")"

# Fail-safe: don't fail if claude-notifier not running
curl -s -X POST http://localhost:19199/notify \
    -H "Content-Type: application/json" \
    -d '{"message": "'"$PROJECT_NAME"' completed at '"$(date '+%H:%M:%S')"'"}' \
    --connect-timeout 1 \
    || true
