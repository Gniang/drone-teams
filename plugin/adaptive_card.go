package plugin

type WebhookContent struct {
	Attachments []AdaptiveCard `json:"attachments"`
}

type AdaptiveCard struct {
	ContentType string              `json:"contentType"`
	Content     AdaptiveCardContent `json:"content"`
}

type AdaptiveCardContent struct {
	Schema  string        `json:"$schema"`
	Type    string        `json:"type"`
	Version string        `json:"version"`
	Body    []interface{} `json:"body"`
}
