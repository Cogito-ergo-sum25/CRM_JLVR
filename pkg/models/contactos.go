package models

import (
	"time"
	"gorm.io/gorm"
)

// Contacto es la entidad principal
type Contacto struct {
	gorm.Model
	Nombre           string    `gorm:"size:255;not null"`
	Email            string    `gorm:"size:255"`
	Telefono         string    `gorm:"size:50"`
	TipoRelacion     string    `gorm:"size:50"` 
	Expediente       string    `gorm:"size:100"`
	Juzgado          string    `gorm:"size:255"`
	FechaCumpleanios *time.Time `gorm:"type:date"`
	RecomendadoPor   string    `gorm:"size:255"`
	Notas            string    `gorm:"type:text"`
	
	// ESTOS SON LOS CAMPOS QUE TE FALTAN:
	// Nota: Deben empezar con Mayúscula para ser públicos
	Nominas    []Nomina    `gorm:"foreignKey:ContactoID"`
	Familiares []Familiar  `gorm:"foreignKey:ContactoID"`
}

// Nomina para el control de honorarios
type Nomina struct {
	gorm.Model
	ContactoID uint
	Fecha      time.Time `gorm:"type:date"`
	Cantidad   float64   `gorm:"type:decimal(10,2)"`
	Motivo     string    `gorm:"size:255"`
}

// Familiar para red de contactos
type Familiar struct {
	gorm.Model
	ContactoID uint
	Nombre     string `gorm:"size:255"`
	Parentesco string `gorm:"size:100"`
	Telefono   string `gorm:"size:50"`
}