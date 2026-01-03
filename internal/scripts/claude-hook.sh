#!/bin/bash
# Claude Code Stop Hook - sends notification with transcript summary

# Read hook input from stdin
HOOK_INPUT=$(cat)
TRANSCRIPT_PATH=$(echo "$HOOK_INPUT" | jq -r '.transcript_path // empty')

# Default fallback
PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
PROJECT_NAME="$(basename "$PROJECT_DIR")"
SUMMARY=""

# Find claude-notifier binary (check common locations)
NOTIFIER_BIN=""
for path in "$HOME/.local/bin/claude-notifier" "/usr/local/bin/claude-notifier" "$HOME/go/bin/claude-notifier"; do
    if [ -x "$path" ]; then
        NOTIFIER_BIN="$path"
        break
    fi
done

# Generate summary using Go parser if binary and transcript exist
if [ -n "$NOTIFIER_BIN" ] && [ -n "$TRANSCRIPT_PATH" ] && [ -f "$TRANSCRIPT_PATH" ]; then
    SUMMARY=$("$NOTIFIER_BIN" summary "$TRANSCRIPT_PATH" 2>/dev/null)
fi

# Build message: time + project name + optional summary
TIME=$(date +%H:%M:%S)
if [ -n "$SUMMARY" ]; then
    MESSAGE="${TIME} ${PROJECT_NAME} - ${SUMMARY}"
else
    MESSAGE="${TIME} ${PROJECT_NAME}"
fi

# Send notification (fail-safe)
# Use jq to safely escape JSON to prevent injection from special chars in paths/summaries
JSON_PAYLOAD=$(jq -n --arg msg "$MESSAGE" '{message: $msg}')
curl -s -X POST http://localhost:19199/notify \
    -H "Content-Type: application/json" \
    -d "$JSON_PAYLOAD" \
    --connect-timeout 1 \
    || true
