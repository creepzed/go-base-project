package models

type SendRequest struct {
	OrderID string       `json:"order_id" validate:"required"`
	Channel string       `json:"channel" validate:"required"`
	Event   string       `json:"event" validate:"required"`
	Content EmailContent `json:"content" validate:"required"`
}

type EmailContent struct {
	Subject string `json:"subject"`
	From    struct {
		Address string `json:"address" validate:"required,email"`
		Name    string `json:"name" validate:"required"`
	} `json:"from" validate:"required"`
	To struct {
		Address string `json:"address" validate:"required,email"`
		Name    string `json:"name" validate:"required"`
	} `json:"to" validate:"required"`
	Attachments []struct {
		Type     string `json:"type" validate:"required"`
		Filename string `json:"filename" validate:"required"`
		Content  string `json:"content" validate:"required,base64"`
	} `json:"attachments"`
	Body map[string]interface{} `json:"body" validate:"required"`
}
