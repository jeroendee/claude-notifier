package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleVersionFlag_WithVersionFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--version"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var buf bytes.Buffer
	got := handleVersionFlag(&buf)

	if !got {
		t.Error("handleVersionFlag() = false, want true")
	}

	output := buf.String()
	if !strings.Contains(output, "claude-notifier") {
		t.Errorf("output = %q, want to contain 'claude-notifier'", output)
	}
}

func TestHandleVersionFlag_WithShortFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-v"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var buf bytes.Buffer
	got := handleVersionFlag(&buf)

	if !got {
		t.Error("handleVersionFlag() = false, want true")
	}

	output := buf.String()
	if !strings.Contains(output, "claude-notifier") {
		t.Errorf("output = %q, want to contain 'claude-notifier'", output)
	}
}

func TestHandleVersionFlag_NoFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var buf bytes.Buffer
	got := handleVersionFlag(&buf)

	if got {
		t.Error("handleVersionFlag() = true, want false")
	}

	if buf.Len() != 0 {
		t.Errorf("output = %q, want empty", buf.String())
	}
}

func TestHandleSummarySubcommand_WithValidTranscript(t *testing.T) {
	tmpDir := t.TempDir()
	transcriptPath := filepath.Join(tmpDir, "test.jsonl")

	content := `{"type":"assistant","message":{"content":[{"type":"text","text":"## Task Complete\n\nDone."}]}}`
	if err := os.WriteFile(transcriptPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var buf bytes.Buffer
	handled, err := handleSummarySubcommand([]string{"summary", transcriptPath}, &buf)

	if err != nil {
		t.Errorf("handleSummarySubcommand() error = %v, want nil", err)
	}
	if !handled {
		t.Error("handleSummarySubcommand() handled = false, want true")
	}
	if got := strings.TrimSpace(buf.String()); got != "Task Complete" {
		t.Errorf("handleSummarySubcommand() output = %q, want %q", got, "Task Complete")
	}
}

func TestHandleSummarySubcommand_WithNoHeadings(t *testing.T) {
	tmpDir := t.TempDir()
	transcriptPath := filepath.Join(tmpDir, "test.jsonl")

	content := `{"type":"assistant","message":{"content":[{"type":"text","text":"Just plain text without headings."}]}}`
	if err := os.WriteFile(transcriptPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var buf bytes.Buffer
	handled, err := handleSummarySubcommand([]string{"summary", transcriptPath}, &buf)

	if err != nil {
		t.Errorf("handleSummarySubcommand() error = %v, want nil", err)
	}
	if !handled {
		t.Error("handleSummarySubcommand() handled = false, want true")
	}
	if got := buf.String(); got != "" {
		t.Errorf("handleSummarySubcommand() output = %q, want empty string", got)
	}
}

func TestHandleSummarySubcommand_WithNonexistentFile(t *testing.T) {
	var buf bytes.Buffer
	handled, err := handleSummarySubcommand([]string{"summary", "/nonexistent/path.jsonl"}, &buf)

	if err == nil {
		t.Error("handleSummarySubcommand() error = nil, want error")
	}
	if !handled {
		t.Error("handleSummarySubcommand() handled = false, want true")
	}
}

func TestHandleSummarySubcommand_WithMissingPathArgument(t *testing.T) {
	var buf bytes.Buffer
	handled, err := handleSummarySubcommand([]string{"summary"}, &buf)

	if err == nil {
		t.Error("handleSummarySubcommand() error = nil, want error")
	}
	if !handled {
		t.Error("handleSummarySubcommand() handled = false, want true")
	}
}

func TestHandleSummarySubcommand_WhenNotSummary(t *testing.T) {
	var buf bytes.Buffer
	handled, err := handleSummarySubcommand([]string{"other"}, &buf)

	if err != nil {
		t.Errorf("handleSummarySubcommand() error = %v, want nil", err)
	}
	if handled {
		t.Error("handleSummarySubcommand() handled = true, want false")
	}
}

func TestHandleSummarySubcommand_WithEmptyArgs(t *testing.T) {
	var buf bytes.Buffer
	handled, err := handleSummarySubcommand([]string{}, &buf)

	if err != nil {
		t.Errorf("handleSummarySubcommand() error = %v, want nil", err)
	}
	if handled {
		t.Error("handleSummarySubcommand() handled = true, want false")
	}
}
