package models

// TrackingPost ...
type TrackingPost struct {
	CreateAt int64  `json:"create_at"`
	UserID   string `json:"uid"`
	Payload  `json:"payload"`
}

// Payload ...
type Payload struct {
	ContentID string `json:"content_id"`
}
