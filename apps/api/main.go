package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rango/api/internal"
	"rango/auth"
	"rango/platform/db/generated"
	"rango/platform/db/repositories"
	"rango/platform/eventbus"
	"rango/platform/storage"
	"rango/router"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	dbCtx := context.Background()

	conn, err := pgx.Connect(
		dbCtx,
		os.Getenv("DB_CONN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	queries := generated.New(conn)

	authSrv := createAuthService(queries)
	documentSrv := createDocumentSrv(queries)
	router := createRootRouter()
	eventBus := eventbus.New()

	app := &App{
		router:      router,
		authSrv:     authSrv,
		documentSrv: documentSrv,
		eventBus:    eventBus,
	}

	addRoutes(app)

	app.ServeHttp(":3000")
}

// App is a container for all services needed for the http api and background workers
type App struct {
	router      *router.Root
	authSrv     *auth.AuthService
	documentSrv *internal.DocumentService
	eventBus    *eventbus.EventBus
}

func (app *App) ServeHttp(addr string) {
	mux := app.router.BuildMux()

	log.Printf("Listening on %s", addr)

	err := http.ListenAndServe(addr, mux)

	if err != nil {
		log.Fatal(err)
	}
}

type RealClock struct{}

// Now implements core.Clock.
func (r *RealClock) Now() time.Time {
	return time.Now()
}

func createAuthService(queries *generated.Queries) *auth.AuthService {
	return auth.NewAuthService(auth.NewAuthServiceParams{
		Repository: repositories.NewPGAuthRepository(queries),
		Clock:      &RealClock{},
	})
}

func createDocumentSrv(queries *generated.Queries) *internal.DocumentService {
	workingDir, err := os.Getwd()

	if err != nil {
		log.Fatalf(err.Error())
	}

	storage := storage.NewDiscStorage(storage.NewDiscStorageParams{
		BasePath: workingDir,
	})

	return &internal.DocumentService{
		Backend:    storage,
		Repository: repositories.NewPGDocumentRepository(queries),
	}
}

func createRootRouter() *router.Root {
	root := router.NewRootRouter()
	root.Use(internal.LoggingMiddleware)

	return root
}
