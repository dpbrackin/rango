package main

import (
	"rango/api/handlers"
	"rango/api/internal"
)

func addEventHandlers(app *App) {
	documentEventHandler := handlers.NewDocumentEventHandler(handlers.NewDocumentEventHandlerParams{
		IndexSrv:    app.indexSrv,
		DocumentSrv: app.documentSrv,
	})

	docCreatedChan := app.eventBus.Subscribe(internal.DOCUMENT_CREATED_TOPIC)

	go func() {
		for event := range docCreatedChan {
			documentEventHandler.HandleDocumentCreatedEvent(event)
		}
	}()

}
