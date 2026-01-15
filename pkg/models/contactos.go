package models

import (
	"gorm.io/gorm"
)

type Contacto struct {
	gorm.Model        // Esto añade ID, CreatedAt, UpdatedAt, DeletedAt automáticamente
	Nombre       string `gorm:"size:255;not null"`
	Email        string `gorm:"uniqueIndex"`
	Telefono     string
	TipoContacto string // 'Abogacía', 'Biomédica', etc.
	Metadata     string `gorm:"type:jsonb"` // Aquí guardaremos los datos extra
}