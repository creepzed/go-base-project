package models

type ErrorResponse struct {
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

type SentResponse struct {
	MessageId   string `json:"message_id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
