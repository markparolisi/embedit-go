package media

import (
	"time"
)

type Model struct {
	Name         string    `json:"name"`
	Service      string    `json:"service"`
	Source       string    `json:"source"`
	Type         string    `json:"type"`
	Created      time.Time `json:"createdAt"`
	ThumbnailURL string    `json:"thumbnailURL"`
	MediaURL     string    `json:"mediaURL"`
	Credit       string    `json:"credit"`
}
