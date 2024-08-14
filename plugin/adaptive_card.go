package plugin

type WebhookContent struct {
	Attachments []AdaptiveCard `json:"attachments"`
}

type CardContentInterface interface{}

type AdaptiveCard struct {
	ContentType string              `json:"contentType"`
	Content     AdaptiveCardContent `json:"content"`
}

type AdaptiveCardContent struct {
	Schema  string          `json:"$schema"`
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Body    []CardContainer `json:"body"`
}

type AdaptiveCardBody struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Wrap bool   `json:"wrap"`
}

type CardContainer struct {
	Type  string          `json:"type"`
	Items []CardColumnSet `json:"items"`
	Style string          `json:"style"`
}
type CardColumnSet struct {
	Type    string       `json:"type"`
	Columns []CardColumn `json:"columns"`
}
type CardColumn struct {
	Type                     string          `json:"type"`
	Width                    string          `json:"width"`
	Items                    []CardTextBlock `json:"items"`
	VerticalContentAlignment string          `json:"verticalContentAlignment"`
	HorizontalAlignment      string          `json:"horizontalAlignment"`
}

type CardTextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Wrap bool   `json:"wrap"`
}
