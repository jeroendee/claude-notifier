package transcript

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// headingRegex matches markdown headings (## through ######).
var headingRegex = regexp.MustCompile(`^#{2,6}\s+(.+)$`)

// ParseTranscript reads a JSONL transcript file and returns all messages.
func ParseTranscript(path string) ([]Message, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open transcript: %w", err)
	}
	defer file.Close()

	var messages []Message
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return nil, fmt.Errorf("parse line %d: %w", lineNum, err)
		}

		messages = append(messages, msg)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read transcript: %w", err)
	}

	return messages, nil
}

// ExtractLastHeading scans assistant messages for markdown headings and
// returns the last heading found (without the # prefix).
func ExtractLastHeading(messages []Message) string {
	var lastHeading string

	for _, msg := range messages {
		// Only process assistant messages
		if msg.Type != "assistant" {
			continue
		}

		for _, block := range msg.Message.Content {
			// Only process text content blocks
			if block.Type != "text" {
				continue
			}

			// Scan each line for headings
			lines := strings.Split(block.Text, "\n")
			for _, line := range lines {
				matches := headingRegex.FindStringSubmatch(line)
				if len(matches) >= 2 {
					lastHeading = strings.TrimSpace(matches[1])
				}
			}
		}
	}

	return lastHeading
}
