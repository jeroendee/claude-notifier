package transcript

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTranscript_ReturnsErrorForNonExistentFile(t *testing.T) {
	t.Parallel()

	_, err := ParseTranscript("/nonexistent/path/to/file.jsonl")

	if err == nil {
		t.Error("ParseTranscript() expected error for non-existent file, got nil")
	}
}

func TestParseTranscript_ReturnsErrorForInvalidJSONLine(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.jsonl")

	content := `{"type":"assistant","message":{"content":[{"type":"text","text":"hello"}]}}
not valid json
{"type":"user","message":{"content":[{"type":"text","text":"hi"}]}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := ParseTranscript(filePath)

	if err == nil {
		t.Error("ParseTranscript() expected error for invalid JSON, got nil")
	}
}

func TestParseTranscript_ParsesEmptyFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "empty.jsonl")

	if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("ParseTranscript() messages = %d, want 0", len(messages))
	}
}

func TestParseTranscript_ParsesSingleAssistantMessage(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "single.jsonl")

	content := `{"type":"assistant","message":{"content":[{"type":"text","text":"Hello, I can help you."}]}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("ParseTranscript() messages = %d, want 1", len(messages))
	}

	if messages[0].Type != "assistant" {
		t.Errorf("ParseTranscript() message type = %q, want %q", messages[0].Type, "assistant")
	}
}

func TestParseTranscript_ParsesMultipleMessagesOfDifferentTypes(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "multi.jsonl")

	content := `{"type":"user","message":{"content":[{"type":"text","text":"Help me"}]}}
{"type":"assistant","message":{"content":[{"type":"text","text":"Sure!"}]}}
{"type":"user","message":{"content":[{"type":"text","text":"Thanks"}]}}
{"type":"assistant","message":{"content":[{"type":"text","text":"Welcome"}]}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 4 {
		t.Errorf("ParseTranscript() messages = %d, want 4", len(messages))
	}
}

func TestParseTranscript_FiltersOnlyAssistantMessages(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "mixed.jsonl")

	content := `{"type":"user","message":{"content":[{"type":"text","text":"Help me"}]}}
{"type":"assistant","message":{"content":[{"type":"text","text":"Sure!"}]}}
{"type":"system","message":{"content":[{"type":"text","text":"System msg"}]}}
{"type":"assistant","message":{"content":[{"type":"text","text":"Done"}]}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	// Count assistant messages
	assistantCount := 0
	for _, m := range messages {
		if m.Type == "assistant" {
			assistantCount++
		}
	}

	if assistantCount != 2 {
		t.Errorf("ParseTranscript() assistant messages = %d, want 2", assistantCount)
	}
}

func TestParseTranscript_HandlesNestedContentArrayStructure(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nested.jsonl")

	content := `{"type":"assistant","message":{"content":[{"type":"text","text":"First block"},{"type":"text","text":"Second block"}]}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("ParseTranscript() messages = %d, want 1", len(messages))
	}

	if len(messages[0].Message.Content) != 2 {
		t.Errorf("ParseTranscript() content blocks = %d, want 2", len(messages[0].Message.Content))
	}
}

func TestExtractLastHeading_ReturnsEmptyStringWhenNoHeadingsExist(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "This is plain text without any headings."},
				},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "" {
		t.Errorf("ExtractLastHeading() = %q, want empty string", heading)
	}
}

func TestExtractLastHeading_ExtractsSingleMarkdownHeading(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "## My Heading\n\nSome content here."},
				},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "My Heading" {
		t.Errorf("ExtractLastHeading() = %q, want %q", heading, "My Heading")
	}
}

func TestExtractLastHeading_ExtractsLastHeadingFromMultipleHeadings(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "## First Heading\n\nContent\n\n## Second Heading\n\nMore content\n\n## Last Heading\n\nFinal content"},
				},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "Last Heading" {
		t.Errorf("ExtractLastHeading() = %q, want %q", heading, "Last Heading")
	}
}

func TestExtractLastHeading_StripsHashPrefixFromHeading(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "### Triple Hash Heading"},
				},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "Triple Hash Heading" {
		t.Errorf("ExtractLastHeading() = %q, want %q", heading, "Triple Hash Heading")
	}
}

func TestExtractLastHeading_HandlesDifferentHeadingLevels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "level 2 heading",
			text: "## Level Two",
			want: "Level Two",
		},
		{
			name: "level 3 heading",
			text: "### Level Three",
			want: "Level Three",
		},
		{
			name: "level 4 heading",
			text: "#### Level Four",
			want: "Level Four",
		},
		{
			name: "level 5 heading",
			text: "##### Level Five",
			want: "Level Five",
		},
		{
			name: "level 6 heading",
			text: "###### Level Six",
			want: "Level Six",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messages := []Message{
				{
					Type: "assistant",
					Message: MessageContent{
						Content: []ContentBlock{
							{Type: "text", Text: tt.text},
						},
					},
				},
			}

			heading := ExtractLastHeading(messages)

			if heading != tt.want {
				t.Errorf("ExtractLastHeading() = %q, want %q", heading, tt.want)
			}
		})
	}
}

func TestExtractLastHeading_IgnoresHeadingsInUserMessages(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "user",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "## User Heading Should Be Ignored"},
				},
			},
		},
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "## Assistant Heading"},
				},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "Assistant Heading" {
		t.Errorf("ExtractLastHeading() = %q, want %q", heading, "Assistant Heading")
	}
}

func TestExtractLastHeading_HandlesContentBlocksWithTextTypeOnly(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "tool_use", Text: "## Tool Heading"},
					{Type: "text", Text: "## Text Heading"},
					{Type: "image", Text: "## Image Heading"},
				},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "Text Heading" {
		t.Errorf("ExtractLastHeading() = %q, want %q", heading, "Text Heading")
	}
}

func TestExtractLastHeading_HandlesMultipleAssistantMessages(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "## First Message Heading"},
				},
			},
		},
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "## Second Message Heading"},
				},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "Second Message Heading" {
		t.Errorf("ExtractLastHeading() = %q, want %q", heading, "Second Message Heading")
	}
}

func TestExtractLastHeading_HandlesEmptyMessages(t *testing.T) {
	t.Parallel()

	var messages []Message

	heading := ExtractLastHeading(messages)

	if heading != "" {
		t.Errorf("ExtractLastHeading() = %q, want empty string", heading)
	}
}

func TestExtractLastHeading_HandlesMessagesWithEmptyContent(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{},
			},
		},
	}

	heading := ExtractLastHeading(messages)

	if heading != "" {
		t.Errorf("ExtractLastHeading() = %q, want empty string", heading)
	}
}

func TestParseTranscript_ParsesStringContent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "string_content.jsonl")

	// User messages have content as a string, not an array
	content := `{"type":"user","message":{"content":"Help me with this task"}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("ParseTranscript() messages = %d, want 1", len(messages))
	}

	if messages[0].Type != "user" {
		t.Errorf("ParseTranscript() message type = %q, want %q", messages[0].Type, "user")
	}

	if len(messages[0].Message.Content) != 1 {
		t.Fatalf("ParseTranscript() content blocks = %d, want 1", len(messages[0].Message.Content))
	}

	if messages[0].Message.Content[0].Type != "text" {
		t.Errorf("ParseTranscript() content type = %q, want %q", messages[0].Message.Content[0].Type, "text")
	}

	if messages[0].Message.Content[0].Text != "Help me with this task" {
		t.Errorf("ParseTranscript() content text = %q, want %q", messages[0].Message.Content[0].Text, "Help me with this task")
	}
}

func TestParseTranscript_ParsesMixedStringAndArrayContent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "mixed_content.jsonl")

	// Real-world transcript format: user messages have string content, assistant messages have array content
	content := `{"type":"user","message":{"content":"Help me with this task"}}
{"type":"assistant","message":{"content":[{"type":"text","text":"Sure, I can help you."}]}}
{"type":"user","message":{"content":"Thanks!"}}
{"type":"assistant","message":{"content":[{"type":"text","text":"## Summary\n\nTask completed."}]}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 4 {
		t.Fatalf("ParseTranscript() messages = %d, want 4", len(messages))
	}

	// Verify user message (string content)
	if messages[0].Message.Content[0].Text != "Help me with this task" {
		t.Errorf("ParseTranscript() user content = %q, want %q", messages[0].Message.Content[0].Text, "Help me with this task")
	}

	// Verify assistant message (array content)
	if messages[1].Message.Content[0].Text != "Sure, I can help you." {
		t.Errorf("ParseTranscript() assistant content = %q, want %q", messages[1].Message.Content[0].Text, "Sure, I can help you.")
	}

	// Verify heading extraction still works
	heading := ExtractLastHeading(messages)
	if heading != "Summary" {
		t.Errorf("ExtractLastHeading() = %q, want %q", heading, "Summary")
	}
}

func TestParseTranscript_ParsesEmptyStringContent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "empty_string.jsonl")

	content := `{"type":"user","message":{"content":""}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("ParseTranscript() messages = %d, want 1", len(messages))
	}

	if len(messages[0].Message.Content) != 1 {
		t.Fatalf("ParseTranscript() content blocks = %d, want 1", len(messages[0].Message.Content))
	}

	if messages[0].Message.Content[0].Text != "" {
		t.Errorf("ParseTranscript() content text = %q, want empty string", messages[0].Message.Content[0].Text)
	}
}

func TestParseTranscript_ParsesWhitespaceOnlyContent(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "whitespace.jsonl")

	content := `{"type":"user","message":{"content":"   "}}`

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	messages, err := ParseTranscript(filePath)

	if err != nil {
		t.Errorf("ParseTranscript() unexpected error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("ParseTranscript() messages = %d, want 1", len(messages))
	}

	if len(messages[0].Message.Content) != 1 {
		t.Fatalf("ParseTranscript() content blocks = %d, want 1", len(messages[0].Message.Content))
	}

	if messages[0].Message.Content[0].Text != "   " {
		t.Errorf("ParseTranscript() content text = %q, want %q", messages[0].Message.Content[0].Text, "   ")
	}
}

func TestExtractSessionSummary_ExtractsSingleSessionMarker(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "Task completed.\n\n[SESSION: Fixed transcript parser string handling]"},
				},
			},
		},
	}

	summary := ExtractSessionSummary(messages)

	if summary != "Fixed transcript parser string handling" {
		t.Errorf("ExtractSessionSummary() = %q, want %q", summary, "Fixed transcript parser string handling")
	}
}

func TestExtractSessionSummary_ReturnsLastSessionMarkerFromMultipleMessages(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "[SESSION: First task completed]"},
				},
			},
		},
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "[SESSION: Second task completed]"},
				},
			},
		},
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "[SESSION: Final task completed]"},
				},
			},
		},
	}

	summary := ExtractSessionSummary(messages)

	if summary != "Final task completed" {
		t.Errorf("ExtractSessionSummary() = %q, want %q", summary, "Final task completed")
	}
}

func TestExtractSessionSummary_ReturnsEmptyStringWhenNoSessionMarkerExists(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "This is plain text without any session markers."},
				},
			},
		},
	}

	summary := ExtractSessionSummary(messages)

	if summary != "" {
		t.Errorf("ExtractSessionSummary() = %q, want empty string", summary)
	}
}

func TestExtractSessionSummary_IgnoresSessionMarkerInUserMessages(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "user",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "[SESSION: User session should be ignored]"},
				},
			},
		},
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "[SESSION: Assistant session marker]"},
				},
			},
		},
	}

	summary := ExtractSessionSummary(messages)

	if summary != "Assistant session marker" {
		t.Errorf("ExtractSessionSummary() = %q, want %q", summary, "Assistant session marker")
	}
}

func TestExtractSessionSummary_ReturnsEmptyStringForEmptyMessages(t *testing.T) {
	t.Parallel()

	var messages []Message

	summary := ExtractSessionSummary(messages)

	if summary != "" {
		t.Errorf("ExtractSessionSummary() = %q, want empty string", summary)
	}
}

func TestExtractSessionSummary_TrimsWhitespaceFromSessionValue(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "[SESSION:   Whitespace trimmed   ]"},
				},
			},
		},
	}

	summary := ExtractSessionSummary(messages)

	if summary != "Whitespace trimmed" {
		t.Errorf("ExtractSessionSummary() = %q, want %q", summary, "Whitespace trimmed")
	}
}

func TestExtractSessionSummary_OnlyChecksTextContentBlocks(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "tool_use", Text: "[SESSION: Tool session ignored]"},
					{Type: "text", Text: "[SESSION: Text session found]"},
					{Type: "image", Text: "[SESSION: Image session ignored]"},
				},
			},
		},
	}

	summary := ExtractSessionSummary(messages)

	if summary != "Text session found" {
		t.Errorf("ExtractSessionSummary() = %q, want %q", summary, "Text session found")
	}
}

func TestExtractSessionSummary_ReturnsLastSessionMarkerFromSameTextBlock(t *testing.T) {
	t.Parallel()

	messages := []Message{
		{
			Type: "assistant",
			Message: MessageContent{
				Content: []ContentBlock{
					{Type: "text", Text: "[SESSION: First summary]\n\nSome work done.\n\n[SESSION: Final summary]"},
				},
			},
		},
	}

	summary := ExtractSessionSummary(messages)

	if summary != "Final summary" {
		t.Errorf("ExtractSessionSummary() = %q, want %q", summary, "Final summary")
	}
}
