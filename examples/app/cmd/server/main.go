package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/gtantech/pdm/examples/app/internal/view"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", templ.Handler(view.Home()).ServeHTTP)
	http.ListenAndServe(":8080", r)
}
