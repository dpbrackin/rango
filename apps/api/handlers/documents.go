package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"rango/api/internal"
	"rango/api/internal/core"
	"rango/api/internal/eventbus"
)

type DocumentsHTTPHandler struct {
	DocSrv *internal.DocumentService
}

func (h *DocumentsHTTPHandler) Upload(w http.ResponseWriter, r *http.Request) {
	reader, header, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user, ok := r.Context().Value(internal.USER_CONTEXT_KEY).(core.User)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user"))
		return
	}

	_, err = h.DocSrv.CreateDocument(r.Context(), internal.AddDocumentParams{
		User:   user,
		Name:   header.Filename,
		Reader: reader,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

type DocumentEventHandler struct {
	indexSrv    *internal.IndexService
	documentSrv *internal.DocumentService
}

type NewDocumentEventHandlerParams struct {
	IndexSrv    *internal.IndexService
	DocumentSrv *internal.DocumentService
}

func NewDocumentEventHandler(params NewDocumentEventHandlerParams) *DocumentEventHandler {
	return &DocumentEventHandler{
		indexSrv:    params.IndexSrv,
		documentSrv: params.DocumentSrv,
	}
}

func (w *DocumentEventHandler) HandleDocumentCreatedEvent(event eventbus.Event) {
	doc, ok := event.Data.(core.Document)

	if !ok {
		slog.Error("Failed to convert event data to core.Document")
		return
	}

	ctx := context.Background()

	err := w.documentSrv.ExtractContent(ctx, &doc)

	if err != nil {
		slog.Error("Content extraction failed", slog.String("error", err.Error()))
		return
	}

	err = w.indexSrv.IndexDocument(ctx, internal.IndexDocumentParams{
		Document: doc,
	})

	if err != nil {
		slog.Error("Indexing Failed", slog.String("error", err.Error()))
	}

}
