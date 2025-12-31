package transcript

import "encoding/json"

// Message represents a single message entry in a Claude transcript JSONL file.
type Message struct {
	Type    string         `json:"type"`
	Message MessageContent `json:"message"`
}

// MessageContent holds the content of a message.
// Content can be either a string (user messages) or an array of ContentBlock (assistant messages).
type MessageContent struct {
	Content []ContentBlock `json:"content"`
}

// UnmarshalJSON implements custom unmarshaling for MessageContent to handle
// both string content (user messages) and array content (assistant messages).
// Returns empty content (nil) if neither format can be parsed.
func (mc *MessageContent) UnmarshalJSON(data []byte) error {
	// Try array format first (assistant messages)
	type rawMessageContent struct {
		Content []ContentBlock `json:"content"`
	}
	var raw rawMessageContent
	if err := json.Unmarshal(data, &raw); err == nil {
		mc.Content = raw.Content
		return nil
	}

	// Fall back to string format (user messages)
	type stringMessageContent struct {
		Content string `json:"content"`
	}
	var strContent stringMessageContent
	if err := json.Unmarshal(data, &strContent); err == nil {
		mc.Content = []ContentBlock{{Type: "text", Text: strContent.Content}}
		return nil
	}

	// If neither works, return empty content
	mc.Content = nil
	return nil
}

// ContentBlock represents a single content block within a message.
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
