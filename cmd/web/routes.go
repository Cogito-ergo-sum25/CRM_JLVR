package main

import (
	"net/http"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/config"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/nuevo-contacto", handlers.Repo.NuevoContacto)
	mux.Post("/nuevo-contacto", handlers.Repo.PostNuevoContacto)

	// Servir archivos estáticos (CSS/JS)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}