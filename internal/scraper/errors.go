package scraper

import "errors"

var (
	ErrInvalidResponse   = errors.New("invalid server response")
	ErrMissingServerInfo = errors.New("missing server info")
)
