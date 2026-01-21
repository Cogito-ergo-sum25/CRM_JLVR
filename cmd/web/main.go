package main

import (
	"log"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/models"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/handlers"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/render"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	dsn := "host=localhost user=usuario_crm password=password_crm dbname=crm_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}

	err = db.AutoMigrate(&models.Contacto{}, &models.Nomina{}, &models.Familiar{}, &models.FechaImportante{})
	if err != nil {
		log.Fatal("Error en migración:", err)
	}

	log.Println("¡Tablas migradas y conexión exitosa!")

	tc, _ := render.CreateTemplateCache()
    app := config.AppConfig{
        TemplateCache: tc,
        UseCache:      false, // Para que Tailwind se refresque al editar
    }

    repo := handlers.NewRepo(&app, db)
    handlers.NewHandlers(repo)
    render.NewTemplates(&app)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: routes(&app),
    }

    log.Println("Servidor corriendo en http://localhost:8080")
    srv.ListenAndServe()
}