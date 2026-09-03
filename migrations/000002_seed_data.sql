-- ==============================================================================
-- MariaDiezmaBack - Seed Data for PostgreSQL
-- Datos de prueba para entorno de desarrollo y pruebas
-- NOTA: Todos los UUIDs tienen exactamente 36 caracteres para compatibilidad estricta con VARCHAR(36)
-- ==============================================================================

-- ==========================================================
-- 1. USERS (Contraseña para todos: AdminPass123!)
-- Hash bcrypt: $2a$10$jE4fNuXtd2MOAx6PiPPKsuLCCWwEqfJzvGF9MD2ChFBpmlNTdWfr6
-- ==========================================================
INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
VALUES 
    ('10000000-0000-0000-0000-000000000001', 'admin@mariadiezma.com', '$2a$10$jE4fNuXtd2MOAx6PiPPKsuLCCWwEqfJzvGF9MD2ChFBpmlNTdWfr6', 'Administrador Principal', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('10000000-0000-0000-0000-000000000002', 'manager@mariadiezma.com', '$2a$10$jE4fNuXtd2MOAx6PiPPKsuLCCWwEqfJzvGF9MD2ChFBpmlNTdWfr6', 'Gestor de Taller', 'manager', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('10000000-0000-0000-0000-000000000003', 'operador@mariadiezma.com', '$2a$10$jE4fNuXtd2MOAx6PiPPKsuLCCWwEqfJzvGF9MD2ChFBpmlNTdWfr6', 'Operador de Atención', 'viewer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- ==========================================================
-- 2. COLLECTIONS (2 Colecciones de ejemplo)
-- ==========================================================
INSERT INTO collections (id, name, image_path, description, created_at, updated_at)
VALUES
    (
        '20000000-0000-0000-0000-000000000001',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_076.jpg',
        'Diseños inspirados en la delicadeza botánica y tonalidades primaverales.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '20000000-0000-0000-0000-000000000002',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_091.jpg',
        'Colección cálida con texturas fluidas y tonos terracota y dorados.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (id) DO NOTHING;

-- ==========================================================
-- 3. DRESSES (2 colecciones con 5 vestidos cada una = 10 vestidos)
-- Cada con sus 3 imágenes consecutivas (MARIA_DIEZMA_xxx.jpg)
-- ==========================================================
INSERT INTO dresses (id, name, collection, image_path, image1_path, image2_path, image3_path, description, created_at, updated_at)
VALUES
    -- -------------------------------------------------------------
    -- Colección 1: Romance (assets/images/romance/)
    -- -------------------------------------------------------------
    (
        '30000000-0000-0000-0000-000000000001',
        'Magnolia',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_001.jpg',
        'assets/images/romance/MARIA_DIEZMA_001.jpg',
        'assets/images/romance/MARIA_DIEZMA_002.jpg',
        'assets/images/romance/MARIA_DIEZMA_003.jpg',
        'de corte sirena con bordados florales artesanales en tul y escote corazón.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000002',
        'Jazmín',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_004.jpg',
        'assets/images/romance/MARIA_DIEZMA_004.jpg',
        'assets/images/romance/MARIA_DIEZMA_005.jpg',
        'assets/images/romance/MARIA_DIEZMA_006.jpg',
        'Diseño evasé con mangas abullonadas y aplicaciones de pétalos en relieve.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000003',
        'Dalia',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_007.jpg',
        'assets/images/romance/MARIA_DIEZMA_007.jpg',
        'assets/images/romance/MARIA_DIEZMA_008.jpg',
        'assets/images/romance/MARIA_DIEZMA_009.jpg',
        'Silueta en A con cuerpo drapeado en gasa de seda y estampado floral sutil.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000004',
        'Camelia',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_010.jpg',
        'assets/images/romance/MARIA_DIEZMA_010.jpg',
        'assets/images/romance/MARIA_DIEZMA_011.jpg',
        'assets/images/romance/MARIA_DIEZMA_012.jpg',
        'midi con falda de vuelo y cinturón joya bordado a mano.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000005',
        'Azahar',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_013.jpg',
        'assets/images/romance/MARIA_DIEZMA_013.jpg',
        'assets/images/romance/MARIA_DIEZMA_014.jpg',
        'assets/images/romance/MARIA_DIEZMA_015.jpg',
        'Elegante confección en encaje chantilly y espalda baja con botonadura forrada.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),

    -- -------------------------------------------------------------
    -- Colección 2: Nayade de Gala (assets/images/nayade/)
    -- -------------------------------------------------------------
    (
        '30000000-0000-0000-0000-000000000006',
        'Siena',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_016.jpg',
        'assets/images/nayade/MARIA_DIEZMA_016.jpg',
        'assets/images/nayade/MARIA_DIEZMA_017.jpg',
        'assets/images/nayade/MARIA_DIEZMA_018.jpg',
        'fluido de gasa de seda con espalda descubierta en tono cálido terracota.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000007',
        'Aurora',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_019.jpg',
        'assets/images/nayade/MARIA_DIEZMA_019.jpg',
        'assets/images/nayade/MARIA_DIEZMA_020.jpg',
        'assets/images/nayade/MARIA_DIEZMA_021.jpg',
        'Confeccionado en crepé satinado con drapeado asimétrico que evoca la luz dorada.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000008',
        'Coral',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_022.jpg',
        'assets/images/nayade/MARIA_DIEZMA_022.jpg',
        'assets/images/nayade/MARIA_DIEZMA_023.jpg',
        'assets/images/nayade/MARIA_DIEZMA_024.jpg',
        'midi con falda plisada soleil y escote halter cruzado.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000009',
        'Terracota',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_025.jpg',
        'assets/images/nayade/MARIA_DIEZMA_025.jpg',
        'assets/images/nayade/MARIA_DIEZMA_026.jpg',
        'assets/images/nayade/MARIA_DIEZMA_027.jpg',
        'Diseño minimalista en lino sedoso con abertura lateral y tirantes espagueti.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000010',
        'Sol Poniente',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_028.jpg',
        'assets/images/nayade/MARIA_DIEZMA_028.jpg',
        'assets/images/nayade/MARIA_DIEZMA_029.jpg',
        'assets/images/nayade/MARIA_DIEZMA_030.jpg',
        'de fiesta en mikado tornasolado con manga capa desestructurada.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (id) DO NOTHING;

-- ==========================================================
-- 4. REQUESTS (Peticiones, Consultas y Citas de prueba)
-- ==========================================================
INSERT INTO requests (
    id, type, status, priority, sender_name, sender_email, sender_phone,
    subject, message, metadata, internal_note, assigned_to, created_at, updated_at
)
VALUES
    (
        '40000000-0000-0000-0000-000000000001',
        'appointment',
        'pending',
        'high',
        'Lucía Ferrero',
        'lucia.ferrero@example.com',
        '+34 611 222 333',
        'Solicitud de Cita: Primera Consulta Novia - 2026-10-15',
        'Nueva cita solicitada para el día 2026-10-15 en el tramo horario 16:00 - 17:00.\nTipo de cita: Primera Consulta Novia.\nContacto: Lucía Ferrero (+34 611 222 333).',
        '{"fecha": "2026-10-15", "tipo_cita": "Primera Consulta Novia", "tramo_horario": "16:00 - 17:00"}'::jsonb,
        '',
        NULL,
        CURRENT_TIMESTAMP - INTERVAL '1 day',
        CURRENT_TIMESTAMP - INTERVAL '1 day'
    ),
    (
        '40000000-0000-0000-0000-000000000002',
        'appointment',
        'in_progress',
        'urgent',
        'Carmen Navarro',
        'carmen.navarro@example.com',
        '+34 622 333 444',
        'Solicitud de Cita: Prueba de - 2026-10-18',
        'Nueva cita solicitada para el día 2026-10-18 en el tramo horario 10:00 - 11:30.\nTipo de cita: Prueba de Vestido.\nContacto: Carmen Navarro (+34 622 333 444).',
        '{"fecha": "2026-10-18", "tipo_cita": "Prueba de Vestido", "tramo_horario": "10:00 - 11:30"}'::jsonb,
        'Llamada telefónica realizada. Clienta interesada en el Magnolia. Cita confirmada.',
        '10000000-0000-0000-0000-000000000002',
        CURRENT_TIMESTAMP - INTERVAL '2 days',
        CURRENT_TIMESTAMP - INTERVAL '1 day'
    ),
    (
        '40000000-0000-0000-0000-000000000003',
        'general',
        'resolved',
        'medium',
        'Beatriz Morales',
        'beatriz.morales@example.com',
        '+34 633 444 555',
        'Consulta sobre disponibilidad de catálogo',
        'Hola María, me gustaría saber si la colección Romance estará disponible en tiendas de Madrid o solo bajo encargo online. ¡Gracias!',
        '{"source": "web_contact"}'::jsonb,
        'Respondido por correo informando sobre el sistema de atelier y cita previa.',
        '10000000-0000-0000-0000-000000000001',
        CURRENT_TIMESTAMP - INTERVAL '5 days',
        CURRENT_TIMESTAMP - INTERVAL '3 days'
    ),
    (
        '40000000-0000-0000-0000-000000000004',
        'budget',
        'pending',
        'medium',
        'Elena Santamaría',
        'elena.santamaria@example.com',
        '+34 644 555 666',
        'Presupuesto a Medida Madrina',
        'Quisiera consultar presupuesto aproximado para confeccionar un a medida para madrina de boda en seda salvaje.',
        '{"presupuesto_estimado": 1800, "evento_fecha": "2026-12-05"}'::jsonb,
        '',
        NULL,
        CURRENT_TIMESTAMP - INTERVAL '3 hours',
        CURRENT_TIMESTAMP - INTERVAL '3 hours'
    )
ON CONFLICT (id) DO NOTHING;
