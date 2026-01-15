package main

import (
	"log"
	"github.com/Cogito-ergo-sum25/CRM_JLVR/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=usuario_crm password=password_crm dbname=crm_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}

	// ¡LA MAGIA! Esto crea o actualiza las tablas automáticamente
	err = db.AutoMigrate(&models.Contacto{})
	if err != nil {
		log.Fatal("Error en migración:", err)
	}

	log.Println("¡Tablas migradas y conexión exitosa!")
}