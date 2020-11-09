package model

type Event struct {
	Type   string `json:"type"`
	Secret string `json:"secret"`
	Data   string `json:"data"`
}
