package handlers

import (
	"encoding/json"
	"net/http"
	"rango/api/internal/auth"
	"time"
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

	user, err := handler.Srv.AuthenticateWithPassword(ctx, auth.PasswordCredentials{
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	session, err := handler.Srv.CreateSession(ctx, user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(err.Error()))
		return
	}

	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    session.ID,
		Quoted:   false,
		Expires:  session.ExpiresAt,
		MaxAge:   int(session.ExpiresAt.Sub(time.Now()).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
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
