package main

import (
	"bytes"
	"flag"
	"os"
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
