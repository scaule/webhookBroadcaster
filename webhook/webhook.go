package webhook

type Webhook struct {
	User  string `json:"user"`
	Data string `json:"data"`
	Type string  `json:"type"`
}
