package handlers

import (
	"log"
	"net/http"
	"time"

	"strconv"

	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/config"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/models"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/render"
	"github.com/go-chi/chi/v5"
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

// EditarContacto renderiza el formulario para editar un contacto existente
func (m *Repository) EditarContacto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var contacto models.Contacto
	m.DB.First(&contacto, id)

	data := make(map[string]interface{})
	data["contacto"] = contacto

	render.RenderTemplate(w, "editar-contacto.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// PostEditarContacto procesa el formulario de edición y actualiza el contacto
func (m *Repository) PostEditarContacto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var contacto models.Contacto
	m.DB.First(&contacto, id)

	r.ParseForm()
	contacto.Nombre = r.Form.Get("nombre")
	contacto.Expediente = r.Form.Get("expediente")
	contacto.Telefono = r.Form.Get("telefono")
	contacto.RecomendadoPor = r.Form.Get("recomendado_por")
	contacto.Notas = r.Form.Get("notas")
    
    // Manejo de fecha opcional
    fechaStr := r.Form.Get("fecha_cumpleanios")
    if fechaStr != "" {
        t, _ := time.Parse("2006-01-02", fechaStr)
        contacto.FechaCumpleanios = &t
    }

	m.DB.Save(&contacto) // Guarda los cambios
	http.Redirect(w, r, "/expediente/"+idStr, http.StatusSeeOther)
}

// EliminarContacto elimina un contacto por su ID
func (m *Repository) EliminarContacto(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, _ := strconv.Atoi(idStr)

    // GORM detecta que tiene gorm.Model y hace Soft Delete
    result := m.DB.Delete(&models.Contacto{}, id)
    
    if result.Error != nil {
        log.Println("Error al eliminar:", result.Error)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) DetalleExpediente(w http.ResponseWriter, r *http.Request) {
    // Obtenemos el ID de la URL (ej: /expediente/5)
    idStr := chi.URLParam(r, "id")
    id, _ := strconv.Atoi(idStr)

    var contacto models.Contacto
    
    // Preload carga las relaciones definidas en el struct
    result := m.DB.Preload("Nominas").Preload("Familiares").First(&contacto, id)
    if result.Error != nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Calculamos el total de honorarios cobrados
    var totalCobrado float64
    for _, pago := range contacto.Nominas {
        totalCobrado += pago.Cantidad
    }

    data := make(map[string]interface{})
    data["contacto"] = contacto
    data["totalCobrado"] = totalCobrado

    render.RenderTemplate(w, "detalle-expediente.page.tmpl", &models.TemplateData{
        Data: data,
    })
}

func (m *Repository) PostNuevoFamiliar(w http.ResponseWriter, r *http.Request) {
    // 1. Obtener el ID del contacto desde la URL
    idStr := chi.URLParam(r, "id")
    contactoID, _ := strconv.Atoi(idStr)

    // 2. Parsear el formulario
    err := r.ParseForm()
    if err != nil {
        log.Println("Error parseando familiar:", err)
        return
    }

    // 3. Crear el objeto con los datos
    familiar := models.Familiar{
        ContactoID: uint(contactoID),
        Nombre:     r.Form.Get("nombre"),
        Parentesco: r.Form.Get("parentesco"),
        Telefono:   r.Form.Get("telefono"),
    }

    // 4. Guardar en la base de datos
    result := m.DB.Create(&familiar)
    if result.Error != nil {
        log.Println("Error guardando familiar:", result.Error)
        return
    }

    // 5. Redirigir de vuelta a la misma página de detalle
    http.Redirect(w, r, "/expediente/"+idStr, http.StatusSeeOther)
}

