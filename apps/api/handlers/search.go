package handlers

import (
	"encoding/json"
	"net/http"
	"rango/api/internal"
	"rango/api/internal/core"
	"rango/api/internal/db/generated"
	"rango/api/internal/search"
	"rango/api/internal/services"
)

type SearchHTTPHandler struct {
	DB generated.DBTX
}

type CreateIndexRequestBody struct {
	Name   string            `json:"name"`
	Engine search.EngineType `json:"engine"`
}

func (h *SearchHTTPHandler) CreateIndex(w http.ResponseWriter, r *http.Request) {
	var body CreateIndexRequestBody

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	service, err := services.NewSearchService(services.NewSearchServiceParams{
		DB:     h.DB,
		Engine: body.Engine,
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	user, ok := r.Context().Value(internal.USER_CONTEXT_KEY).(core.User)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user"))
		return
	}

	idx, err := service.CreateIndex(r.Context(), services.CreateIndexParams{
		Name:   body.Name,
		Engine: string(body.Engine),
		Org:    user.Org,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(idx)
	return
}
