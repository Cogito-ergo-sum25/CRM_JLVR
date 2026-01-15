CREATE TABLE contactos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    telefono VARCHAR(50),
    tipo_contacto VARCHAR(50), -- 'cliente', 'proveedor', 'juzgado', 'hospital'
    metadata JSONB,             -- Aquí vive la magia
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);