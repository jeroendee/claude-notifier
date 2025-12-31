package transcript

// Message represents a single message entry in a Claude transcript JSONL file.
type Message struct {
	Type    string         `json:"type"`
	Message MessageContent `json:"message"`
}

// MessageContent holds the content array of a message.
type MessageContent struct {
	Content []ContentBlock `json:"content"`
}

// ContentBlock represents a single content block within a message.
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
