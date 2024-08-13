package plugin

type WebhookContent struct {
	Attachments []AdaptiveCard `json:"attachments"`
}

type AttachmentContent interface{}

type AdaptiveCard struct {
	ContentType string              `json:"contentType"`
	Content     AdaptiveCardContent `json:"content"`
}

type AdaptiveCardContent struct {
	Schema  string             `json:"$schema"`
	Type    string             `json:"type"`
	Version string             `json:"version"`
	Body    []AdaptiveCardBody `json:"body"`
}

type AdaptiveCardBody struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Wrap bool   `json:"wrap"`
}
