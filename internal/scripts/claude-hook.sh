#!/bin/bash
# Claude Code Stop Hook - sends notification to menubar app

# Fail-safe: don't fail if claude-notifier not running
curl -s -X POST http://localhost:19199/notify \
    -H "Content-Type: application/json" \
    -d '{"message": "Claude Code task completed at '"$(date '+%H:%M:%S')"'"}' \
    --connect-timeout 1 \
    || true
