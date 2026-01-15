package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/config"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/models"
)

var app *config.AppConfig

// NewTemplates configura el paquete render
func NewTemplates(a *config.AppConfig) {
	app = a
}

// Mapa de funciones para los templates (puedes añadir más luego)
var functions = template.FuncMap{
	// Aquí puedes añadir funciones personalizadas para Tailwind si necesitas
}

// AddDefaultData añade datos comunes a todos los templates (como sesiones, etc)
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renderiza un template usando el cache
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Println("Error creando cache de templates:", err)
			return
		}
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Println("No se pudo encontrar el template en el cache:", tmpl)
		return
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)

	err = t.Execute(buf, td)
	if err != nil {
		log.Println("Error ejecutando template:", err)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error escribiendo template al navegador:", err)
	}
}

// CreateTemplateCache crea un cache de templates buscando en la carpeta templates
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// 1. Buscar todas las páginas (archivos .page.tmpl)
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Crear el set de templates para esta página
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// 2. Buscar si hay layouts (archivos .layout.tmpl)
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}