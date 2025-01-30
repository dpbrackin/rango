package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rango/api/internal"
	"rango/api/internal/auth"
	"rango/api/internal/db/generated"
	"rango/api/internal/db/repositories"
	"rango/api/internal/eventbus"
	"rango/api/internal/search"
	"rango/api/internal/storage"
	"rango/router"
	"syscall"
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
	eb := eventbus.New()

	authSrv := createAuthService(queries, conn)
	documentSrv := createDocumentSrv(queries, eb)
	router := createRootRouter()
	indexSrv := createIndexSrv()

	app := &App{
		router:      router,
		authSrv:     authSrv,
		documentSrv: documentSrv,
		indexSrv:    indexSrv,
		eventBus:    eb,
		db:          conn,
	}

	addRoutes(app)
	addEventHandlers(app)

	done := make(chan bool, 1)

	go app.gracefulShutdown(done)

	app.ServeHttp(":3000")

	<-done
	log.Println("Gracefully shutdown")
}

// App is a container for all services needed for the http api and background workers
type App struct {
	router      *router.Root
	authSrv     *auth.AuthService
	documentSrv *internal.DocumentService
	indexSrv    *internal.IndexService
	eventBus    *eventbus.EventBus
	server      *http.Server
	db          generated.DBTX
}

func (app *App) ServeHttp(addr string) {
	mux := app.router.BuildMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	app.server = server

	log.Printf("Listening on %s", addr)

	err := server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (app *App) gracefulShutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shudown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

type RealClock struct{}

// Now implements core.Clock.
func (r *RealClock) Now() time.Time {
	return time.Now()
}

func createAuthService(queries *generated.Queries, conn *pgx.Conn) *auth.AuthService {
	return auth.NewAuthService(auth.NewAuthServiceParams{
		AuthRepository: repositories.NewPGAuthRepository(conn),
		OrgRepository:  repositories.NewPGOrganizationRepository(queries),
		Clock:          &RealClock{},
	})
}

func createDocumentSrv(queries *generated.Queries, eb *eventbus.EventBus) *internal.DocumentService {
	workingDir, err := os.Getwd()

	if err != nil {
		log.Fatalf(err.Error())
	}

	storage := storage.NewDiscStorage(storage.NewDiscStorageParams{
		BasePath: workingDir,
	})

	return &internal.DocumentService{
		Storage:    storage,
		Repository: repositories.NewPGDocumentRepository(queries),
		EventBus:   eb,
	}
}

func createRootRouter() *router.Root {
	root := router.NewRootRouter()
	root.Use(internal.LoggingMiddleware)

	return root
}

func createIndexSrv() *internal.IndexService {
	return &internal.IndexService{
		Embeder: search.NewOpenAIEmbedder(),
	}
}
