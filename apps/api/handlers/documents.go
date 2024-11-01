package handlers

import (
	"net/http"
	"rango/api/internal"
	"rango/core"
)

type DocumentsHandler struct {
	DocSrv *internal.DocumentService
}

func (h *DocumentsHandler) Upload(w http.ResponseWriter, r *http.Request) {
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
