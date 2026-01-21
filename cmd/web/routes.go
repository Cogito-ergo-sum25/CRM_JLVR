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

	mux.Get("/contactos", handlers.Repo.ListaContactos)
	
	mux.Get("/expediente/{id}", handlers.Repo.DetalleExpediente)
	mux.Get("/expediente/{id}/editar", handlers.Repo.EditarContacto)
	mux.Post("/expediente/{id}/editar", handlers.Repo.PostEditarContacto)
	mux.Get("/expediente/{id}/eliminar", handlers.Repo.EliminarContacto)

	mux.Post("/expediente/{id}/familiar", handlers.Repo.PostNuevoFamiliar)
	mux.Post("/expediente/{id}/familiar/{familiarID}/editar", handlers.Repo.PostEditarFamiliar)
	mux.Get("/expediente/{id}/familiar/{familiarID}/eliminar", handlers.Repo.EliminarFamiliar)

	mux.Post("/expediente/{id}/cobro", handlers.Repo.PostNuevoCobro)
	mux.Post("/expediente/{id}/cobro/{cobroID}/editar", handlers.Repo.PostEditarCobro)
	mux.Get("/expediente/{id}/cobro/{cobroID}/eliminar", handlers.Repo.EliminarCobro)


	mux.Get("/calendario", handlers.Repo.Calendario)
	mux.Get("/api/eventos-calendario", handlers.Repo.EventosCalendario)

	// Servir archivos estáticos (CSS/JS)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}