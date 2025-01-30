package main

import (
	"rango/api/handlers"
	"rango/api/internal"
)

func addRoutes(app *App) {
	authHandler := handlers.AuthHandler{
		Srv: app.authSrv,
	}

	docsHandler := handlers.DocumentsHTTPHandler{
		DocSrv: app.documentSrv,
	}

	searchHandler := handlers.SearchHTTPHandler{
		DB: app.db,
	}

	unauthenticatedGroup := app.router.Group("")
	unauthenticatedGroup.RouteFunc("POST /login", authHandler.Login)
	unauthenticatedGroup.RouteFunc("POST /register", authHandler.Register)

	authenticatedGroup := app.router.Group("")
	authenticatedGroup.Use(internal.AuthMiddleware(app.authSrv))
	authenticatedGroup.RouteFunc("POST /document", docsHandler.Upload)
	authenticatedGroup.RouteFunc("POST /index", searchHandler.CreateIndex)
}
