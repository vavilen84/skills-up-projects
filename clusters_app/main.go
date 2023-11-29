package main

import (
	"clusters_app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/servers/stat", handlers.GetServerStatHandler)
	http.ListenAndServe(":3000", r)
}
