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
	var nominas []models.Nomina

	// 1. Obtener datos básicos
	m.DB.Find(&contactos)
	m.DB.Find(&nominas)

	// 2. Calcular total de dinero cobrado
	var totalDinero float64
	for _, n := range nominas {
		totalDinero += n.Cantidad
	}

	// 3. Lógica para Próximos Cumpleaños (Mes actual)
	var cumpleaneros []models.Contacto
	mesActual := time.Now().Month()
	m.DB.Where("EXTRACT(MONTH FROM fecha_cumpleanios) = ?", mesActual).
		Order("EXTRACT(DAY FROM fecha_cumpleanios) ASC").
		Limit(5).
		Find(&cumpleaneros)

	// Enviamos todo al template
	data := make(map[string]interface{})
	data["totalContactos"] = len(contactos)
	data["totalDinero"] = totalDinero
	data["cumpleaneros"] = cumpleaneros

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// NuevoContacto renderiza la página con Tailwind
func (m *Repository) NuevoContacto(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "nuevo-contacto.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ListaContactos(w http.ResponseWriter, r *http.Request) {
    // 1. Obtener el término de búsqueda de la URL (ej: ?search=perez)
    searchTerm := r.URL.Query().Get("search")
    
    var contactos []models.Contacto
    
    // 2. Lógica de filtrado con GORM
    if searchTerm != "" {
        // Busca coincidencias parciales en Nombre o Expediente
        m.DB.Where("nombre ILIKE ? OR expediente ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%").Find(&contactos)
    } else {
        m.DB.Find(&contactos)
    }

    data := make(map[string]interface{})
    data["contactos"] = contactos
    data["searchTerm"] = searchTerm

    render.RenderTemplate(w, "contactos.page.tmpl", &models.TemplateData{
        Data: data,
    })
}

// PostNuevoContacto procesa el formulario y guarda en Postgres
func (m *Repository) PostNuevoContacto(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        log.Println(err)
        return
    }

    // Manejo de la fecha de cumpleaños
    var fechaCumple *time.Time
    fechaStr := r.Form.Get("fecha_cumpleanios")
    if fechaStr != "" {
        t, _ := time.Parse("2006-01-02", fechaStr)
        fechaCumple = &t
    }

    nuevoContacto := models.Contacto{
        Nombre:           r.Form.Get("nombre"),
        Email:            r.Form.Get("email"),
        Telefono:         r.Form.Get("telefono"),
        TipoRelacion:     r.Form.Get("tipo_relacion"),
        Expediente:       r.Form.Get("expediente"),
        Juzgado:          r.Form.Get("juzgado"),
        Notas:            r.Form.Get("notas"),
        RecomendadoPor:   r.Form.Get("recomendado_por"), // Nuevo
        FechaCumpleanios: fechaCumple,                   // Nuevo
    }

    result := m.DB.Create(&nuevoContacto)
    if result.Error != nil {
        log.Println("Error al crear:", result.Error)
        return
    }

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
    
    // Actualizamos TODOS los campos
    contacto.Nombre = r.Form.Get("nombre")
    contacto.Email = r.Form.Get("email")
    contacto.Telefono = r.Form.Get("telefono")
    contacto.Expediente = r.Form.Get("expediente")
    contacto.Juzgado = r.Form.Get("juzgado")
    contacto.TipoRelacion = r.Form.Get("tipo_relacion")
    contacto.RecomendadoPor = r.Form.Get("recomendado_por")
    contacto.Notas = r.Form.Get("notas")

    // Fecha de cumpleaños
    fechaStr := r.Form.Get("fecha_cumpleanios")
    if fechaStr != "" {
        t, _ := time.Parse("2006-01-02", fechaStr)
        contacto.FechaCumpleanios = &t
    } else {
        contacto.FechaCumpleanios = nil
    }

    m.DB.Save(&contacto)
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

func (m *Repository) PostEditarFamiliar(w http.ResponseWriter, r *http.Request) {
    contactoID := chi.URLParam(r, "id")
    familiarID := chi.URLParam(r, "familiarID")

    r.ParseForm()
    m.DB.Model(&models.Familiar{}).Where("id = ?", familiarID).Updates(models.Familiar{
        Nombre:     r.Form.Get("nombre"),
        Parentesco: r.Form.Get("parentesco"),
        Telefono:   r.Form.Get("telefono"),
    })

    http.Redirect(w, r, "/expediente/"+contactoID, http.StatusSeeOther)
}

func (m *Repository) EliminarFamiliar(w http.ResponseWriter, r *http.Request) {
    contactoID := chi.URLParam(r, "id")
    familiarID := chi.URLParam(r, "familiarID")

    m.DB.Delete(&models.Familiar{}, familiarID)

    http.Redirect(w, r, "/expediente/"+contactoID, http.StatusSeeOther)
}

func (m *Repository) PostNuevoCobro(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    contactoID, _ := strconv.Atoi(idStr)

    err := r.ParseForm()
    if err != nil {
        log.Println(err)
        return
    }

    cantidad, _ := strconv.ParseFloat(r.Form.Get("cantidad"), 64)
    fechaStr := r.Form.Get("fecha")
    fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		// Si hay error, podrías usar la fecha actual como respaldo o avisar al usuario
		fecha = time.Now() 
	}

    cobro := models.Nomina{
        ContactoID: uint(contactoID),
        Cantidad:   cantidad,
        Fecha:      fecha,
        Motivo:     r.Form.Get("motivo"),
    }

    m.DB.Create(&cobro)

    http.Redirect(w, r, "/expediente/"+idStr, http.StatusSeeOther)
}

func (m *Repository) PostEditarCobro(w http.ResponseWriter, r *http.Request) {
    contactoID := chi.URLParam(r, "id")
    cobroID := chi.URLParam(r, "cobroID")

    r.ParseForm()
    cantidad, _ := strconv.ParseFloat(r.Form.Get("cantidad"), 64)
    fecha, _ := time.Parse("2006-01-02", r.Form.Get("fecha"))

    // Actualizar directamente en la DB buscando por ID del cobro
    m.DB.Model(&models.Nomina{}).Where("id = ?", cobroID).Updates(models.Nomina{
        Cantidad: cantidad,
        Fecha:    fecha,
        Motivo:   r.Form.Get("motivo"),
    })

    http.Redirect(w, r, "/expediente/"+contactoID, http.StatusSeeOther)
}

func (m *Repository) EliminarCobro(w http.ResponseWriter, r *http.Request) {
    // Necesitamos el ID del cobro
    cobroIDStr := chi.URLParam(r, "cobroID")
    contactoIDStr := chi.URLParam(r, "id")
    
    id, _ := strconv.Atoi(cobroIDStr)

    // Borrado físico o lógico dependiendo de si usas gorm.Model en Nomina
    m.DB.Delete(&models.Nomina{}, id)

    // Redirigir de vuelta al detalle del expediente
    http.Redirect(w, r, "/expediente/"+contactoIDStr, http.StatusSeeOther)
}

