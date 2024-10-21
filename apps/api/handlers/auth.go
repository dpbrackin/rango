package handlers

import (
	"encoding/json"
	"net/http"
	"rango/auth"
)

type AuthHandler struct {
	Srv *auth.AuthService
}

type RegisterRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body LoginRequestBody

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	_, err = handler.Srv.AuthenticateWithPassword(ctx, auth.PasswordCredentials{
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body RegisterRequestBody

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	_, err = handler.Srv.Register(ctx, auth.PasswordCredentials{
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
