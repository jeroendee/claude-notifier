// Package transcript provides parsing utilities for Claude Code session
// transcript files in JSONL format.
//
// The package reads JSONL transcript files where each line is a JSON object
// representing a message in the conversation. It extracts assistant messages
// and can scan for markdown headings within the content.
//
// Example usage:
//
//	messages, err := transcript.ParseTranscript("/path/to/session.jsonl")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	heading := transcript.ExtractLastHeading(messages)
//	if heading != "" {
//		fmt.Printf("Last heading: %s\n", heading)
//	}
package transcript
