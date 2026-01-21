-- 1. Tabla Principal de Contactos (Abogada)
CREATE TABLE IF NOT EXISTS contactos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    telefono VARCHAR(50),
    tipo_relacion VARCHAR(50), 
    expediente VARCHAR(100),
    juzgado VARCHAR(255),
    fecha_cumpleanios DATE,
    recomendado_por VARCHAR(255),
    notas TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 2. Tabla de Nómina / Honorarios (Relacionada a Contactos)
CREATE TABLE IF NOT EXISTS nominas (
    id SERIAL PRIMARY KEY,
    contacto_id INTEGER REFERENCES contactos(id) ON DELETE CASCADE,
    fecha DATE NOT NULL,
    cantidad DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    motivo VARCHAR(255), -- Ej: 'Anticipo para copias', 'Honorarios Mayo'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 3. Tabla de Familiares (Relacionada a Contactos)
CREATE TABLE IF NOT EXISTS familiares (
    id SERIAL PRIMARY KEY,
    contacto_id INTEGER REFERENCES contactos(id) ON DELETE CASCADE,
    nombre VARCHAR(255) NOT NULL,
    parentesco VARCHAR(100), -- Ej: 'Esposo', 'Hijo'
    telefono VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
