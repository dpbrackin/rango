package internal

import (
	"rango/core"
	"rango/platform/eventbus"
)

const (
	DOCUMENT_CREATED_TOPIC = "document.created"
)

func NewDocumentCreatedEvent(doc core.Document) eventbus.Event {
	return eventbus.Event{
		Topic: DOCUMENT_CREATED_TOPIC,
		Data:  doc,
	}
}
