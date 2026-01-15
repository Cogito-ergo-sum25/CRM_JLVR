package handlers

import (
	"log"
	"net/http"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/config"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/models"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/render"
	"gorm.io/gorm"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  *gorm.DB // Añadiremos GORM aquí
}

func NewRepo(a *config.AppConfig, db *gorm.DB) *Repository {
	return &Repository{
		App: a,
		DB:  db,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home es el handler para la página de inicio
// Home es el handler para la página de inicio que lista los expedientes
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	var contactos []models.Contacto

	// Consultar todos los contactos de la base de datos usando GORM
	result := m.DB.Find(&contactos)
	if result.Error != nil {
		log.Println("Error al recuperar contactos:", result.Error)
		http.Error(w, "Error al cargar los expedientes", http.StatusInternalServerError)
		return
	}

	// Preparar los datos para el template
	data := make(map[string]interface{})
	data["contactos"] = contactos

	// Renderizar la página de inicio pasando los datos
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// NuevoContacto renderiza la página con Tailwind
func (m *Repository) NuevoContacto(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "nuevo-contacto.page.tmpl", &models.TemplateData{})
}

// PostNuevoContacto procesa el formulario y guarda en Postgres
func (m *Repository) PostNuevoContacto(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    
    contacto := models.Contacto{
        Nombre:       r.Form.Get("nombre"),
        Email:        r.Form.Get("email"),
        Telefono:     r.Form.Get("telefono"),
        TipoRelacion: r.Form.Get("tipo_relacion"),
        Expediente:   r.Form.Get("expediente"),
        Juzgado:      r.Form.Get("juzgado"),
        Notas:        r.Form.Get("notas"),
    }

    m.DB.Create(&contacto)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

