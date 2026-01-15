DROP TABLE IF EXISTS contactos;
CREATE TABLE contactos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    telefono VARCHAR(50),
    tipo_relacion VARCHAR(50), -- 'Cliente', 'Contraparte', etc.
    expediente VARCHAR(100),   -- Número de caso
    juzgado VARCHAR(255),      -- Ubicación del proceso
    notas TEXT,                -- Historial detallado
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP       -- Para borrado lógico de GORM
);

