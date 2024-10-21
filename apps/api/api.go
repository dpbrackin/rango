package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rango/api/handlers"
	"rango/auth"
	"rango/db/generated"
	"rango/db/repositories"
	"rango/router"
	"time"

	"github.com/jackc/pgx/v5"
)

type API struct {
	rootRouter  *router.Root
	authService *auth.AuthService
}

type RealClock struct{}

func (r RealClock) Now() time.Time {
	return time.Now()
}

func NewAPI() *API {
	root := router.NewRootRouter()

	root.Use(LoggingMiddleware)

	ctx := context.Background()

	conn, err := pgx.Connect(
		ctx,
		os.Getenv("DB_DSN"),
	)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	queries := generated.New(conn)

	authService := auth.NewAuthService(auth.NewAuthServiceParams{
		Repository: repositories.NewPGAuthRepository(queries),
		Clock:      &RealClock{},
	})

	api := &API{
		rootRouter:  root,
		authService: authService,
	}

	api.BuildRoutes()

	return api
}

func (api *API) BuildRoutes() {

	authHandler := handlers.AuthHandler{
		Srv: api.authService,
	}

	unauthenticatedGroup := api.rootRouter.Group("")
	unauthenticatedGroup.RouteFunc("POST /login", authHandler.Login)
	unauthenticatedGroup.RouteFunc("POST /register", authHandler.Register)
}

func (api *API) Serve(addr string) {
	mux := api.rootRouter.BuildMux()

	log.Printf("Listening on %s", addr)

	err := http.ListenAndServe(addr, mux)

	if err != nil {
		log.Fatal(err)
	}
}
