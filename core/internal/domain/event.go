package domain

import (
	"github.com/point-unknown/catalyst/pkg/cloudevent"
)

// CloudEvent is a type alias for the standard Catalyst CloudEvent.
type CloudEvent = cloudevent.Event

// NewEvent helper (legacy support for core)
func NewEvent(source, eventType string, data interface{}) (CloudEvent, error) {
	return cloudevent.New(source, eventType, data)
}
