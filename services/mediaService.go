package services

import "embedit/media"

// Interface to collect all of the Services
type MediaService interface {
	GetMedia(string) ([]media.Model, error)
}
