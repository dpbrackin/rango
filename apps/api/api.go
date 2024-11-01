package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rango/api/handlers"
	"rango/api/internal"
	"rango/auth"
	"rango/db/generated"
	"rango/db/repositories"
	"rango/router"
	"rango/storage"
	"time"

	"github.com/jackc/pgx/v5"
)

type API struct {
	rootRouter  *router.Root
	authService *auth.AuthService
	docService  *internal.DocumentService
}

type RealClock struct{}

func (r RealClock) Now() time.Time {
	return time.Now()
}

func NewAPI() *API {
	root := router.NewRootRouter()

	root.Use(internal.LoggingMiddleware)

	ctx := context.Background()

	conn, err := pgx.Connect(
		ctx,
		os.Getenv("DB_CONN"),
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

	workingDir, err := os.Getwd()

	if err != nil {
		log.Fatalf(err.Error())
	}

	docService := &internal.DocumentService{
		Backend: storage.NewDiscStorage(storage.NewDiscStorageParams{
			BasePath: workingDir,
		}),
		Repository: repositories.NewPGDocumentRepository(queries),
	}

	api := &API{
		rootRouter:  root,
		authService: authService,
		docService:  docService,
	}

	api.BuildRoutes()

	return api
}

func (api *API) BuildRoutes() {
	authHandler := handlers.AuthHandler{
		Srv: api.authService,
	}

	docsHandler := handlers.DocumentsHandler{
		DocSrv: api.docService,
	}

	unauthenticatedGroup := api.rootRouter.Group("")
	unauthenticatedGroup.RouteFunc("POST /login", authHandler.Login)
	unauthenticatedGroup.RouteFunc("POST /register", authHandler.Register)

	authenticatedGroup := api.rootRouter.Group("")
	authenticatedGroup.Use(internal.AuthMiddleware(api.authService))
	authenticatedGroup.RouteFunc("POST /document", docsHandler.Upload)
}

func (api *API) Serve(addr string) {
	mux := api.rootRouter.BuildMux()

	log.Printf("Listening on %s", addr)

	err := http.ListenAndServe(addr, mux)

	if err != nil {
		log.Fatal(err)
	}
}
