package router_test

import (
	"lib/go/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("TEST"))
}

func TestRouteFunc(t *testing.T) {
	router := router.NewRootRouter()
	router.RouteFunc("GET /", testHandler)
	mux := router.BuildMux()

	req := httptest.NewRequest("GET", "/", nil)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
}

func TestRouteGroup(t *testing.T) {
	router := router.NewRootRouter()

	g := router.Group("/group")
	g.RouteFunc("/test", testHandler)

	mux := router.BuildMux()

	req := httptest.NewRequest("GET", "/group/test", nil)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
}
