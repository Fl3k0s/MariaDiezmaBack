-- ==============================================================================
-- Migration: 000004_create_holidays_table.sql
-- Festivos de Madrid definidos año a año para calendario y disponibilidad web.
-- Permite eliminación/actualización en cualquier momento.
-- ==============================================================================

CREATE TABLE IF NOT EXISTS holidays (
    id VARCHAR(36) PRIMARY KEY,
    date VARCHAR(10) NOT NULL UNIQUE, -- 'YYYY-MM-DD'
    year INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    scope VARCHAR(50) DEFAULT 'madrid', -- 'nacional', 'autonomico', 'local'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_holidays_year ON holidays(year);
CREATE INDEX IF NOT EXISTS idx_holidays_date ON holidays(date);

-- ==============================================================================
-- Festivos Oficiales de Madrid para 2025 y 2026
-- ON CONFLICT (date) DO NOTHING asegura idempotencia sin sobreescribir eliminaciones
-- ==============================================================================

INSERT INTO holidays (id, date, year, name, description, scope, created_at, updated_at)
VALUES
    -- Festivos Madrid 2025
    ('50000000-0000-0000-2025-000000000001', '2025-01-01', 2025, 'Año Nuevo', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000002', '2025-01-06', 2025, 'Epifanía del Señor (Reyes Magos)', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000003', '2025-04-17', 2025, 'Jueves Santo', 'Festivo autonómico Comunidad de Madrid', 'autonomico', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000004', '2025-04-18', 2025, 'Viernes Santo', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000005', '2025-05-01', 2025, 'Fiesta del Trabajo', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000006', '2025-05-02', 2025, 'Fiesta de la Comunidad de Madrid', 'Día de la Comunidad de Madrid', 'autonomico', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000007', '2025-05-15', 2025, 'San Isidro Labrador', 'Patrón de Madrid capital', 'local', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000008', '2025-07-25', 2025, 'Santiago Apóstol', 'Festivo autonómico Comunidad de Madrid', 'autonomico', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000009', '2025-08-15', 2025, 'Asunción de la Virgen', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000010', '2025-10-12', 2025, 'Fiesta Nacional de España', 'Día de la Hispanidad', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000011', '2025-11-01', 2025, 'Todos los Santos', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000012', '2025-11-10', 2025, 'Lunes siguiente a La Almudena', 'Patrona de Madrid capital (traslado)', 'local', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000013', '2025-12-06', 2025, 'Día de la Constitución Española', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000014', '2025-12-08', 2025, 'Inmaculada Concepción', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2025-000000000015', '2025-12-25', 2025, 'Natividad del Señor (Navidad)', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

    -- Festivos Madrid 2026
    ('50000000-0000-0000-2026-000000000001', '2026-01-01', 2026, 'Año Nuevo', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000002', '2026-01-06', 2026, 'Epifanía del Señor (Reyes Magos)', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000003', '2026-04-02', 2026, 'Jueves Santo', 'Festivo autonómico Comunidad de Madrid', 'autonomico', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000004', '2026-04-03', 2026, 'Viernes Santo', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000005', '2026-05-01', 2026, 'Fiesta del Trabajo', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000006', '2026-05-02', 2026, 'Fiesta de la Comunidad de Madrid', 'Día de la Comunidad de Madrid', 'autonomico', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000007', '2026-05-15', 2026, 'San Isidro Labrador', 'Patrón de Madrid capital', 'local', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000008', '2026-08-15', 2026, 'Asunción de la Virgen', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000009', '2026-10-12', 2026, 'Fiesta Nacional de España', 'Día de la Hispanidad', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000010', '2026-11-02', 2026, 'Lunes siguiente a Todos los Santos', 'Festivo autonómico Comunidad de Madrid (traslado)', 'autonomico', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000011', '2026-11-09', 2026, 'Nuestra Señora de la Almudena', 'Patrona de Madrid capital', 'local', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000012', '2026-12-07', 2026, 'Lunes siguiente a la Constitución Española', 'Festivo autonómico Comunidad de Madrid (traslado)', 'autonomico', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000013', '2026-12-08', 2026, 'Inmaculada Concepción', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('50000000-0000-0000-2026-000000000014', '2026-12-25', 2026, 'Natividad del Señor (Navidad)', 'Festivo nacional', 'nacional', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (date) DO NOTHING;
