package dto

// AsynqEmailDeliveryPayload is the payload for email delivery task.
type AsynqEmailDeliveryPayload struct {
	UserID     int
	TemplateID string
}

// AsynqImageResizePayload is the payload for image resize task.
type AsynqImageResizePayload struct {
	SourceURL string
}
