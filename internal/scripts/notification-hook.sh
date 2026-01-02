#!/bin/bash
# Claude Code Notification/Permission Hook
# Fires for idle_prompt (60s waiting) and permission requests

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
PROJECT_NAME="$(basename "$PROJECT_DIR")"
TIME=$(date +%H:%M:%S)
MESSAGE="${TIME} ${PROJECT_NAME} - Waiting for input"

curl -s -X POST http://localhost:19199/notify \
    -H "Content-Type: application/json" \
    -d '{"message": "'"$MESSAGE"'"}' \
    --connect-timeout 1 \
    || true
