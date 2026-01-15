package models

import (
	"gorm.io/gorm"
)

type Contacto struct {
	gorm.Model
	Nombre       string `gorm:"size:255;not null"`
	Email        string `gorm:"size:255"`
	Telefono     string `gorm:"size:50"`
	TipoRelacion string `gorm:"size:50"` // Cliente, Contraparte, etc.
	Expediente   string `gorm:"size:100"`
	Juzgado      string `gorm:"size:255"`
	Notas        string `gorm:"type:text"`
}