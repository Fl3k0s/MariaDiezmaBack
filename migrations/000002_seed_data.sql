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
-- 3. DRESSES (Romance: 25 vestidos, Nayade de Gala: 29 vestidos = 54 vestidos)
-- Cada uno con sus 3 imágenes consecutivas (MARIA_DIEZMA_xxx.jpg)
-- ==========================================================
INSERT INTO dresses (id, name, collection, image_path, image1_path, image2_path, image3_path, description, created_at, updated_at)
VALUES
    -- -------------------------------------------------------------
    -- Colección 1: Romance (assets/images/romance/)
    -- -------------------------------------------------------------
    (
        '30000000-0000-0000-0000-000000000001',
        'Aura',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_001.jpg',
        'assets/images/romance/MARIA_DIEZMA_001.jpg',
        'assets/images/romance/MARIA_DIEZMA_002.jpg',
        'assets/images/romance/MARIA_DIEZMA_003.jpg',
        'Diseño etéreo en tul plumeti con cuerpo bordado y sutil escote en V.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000002',
        'Brisa',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_004.jpg',
        'assets/images/romance/MARIA_DIEZMA_004.jpg',
        'assets/images/romance/MARIA_DIEZMA_005.jpg',
        'assets/images/romance/MARIA_DIEZMA_006.jpg',
        'Silueta fluida de gasa de seda con caída vaporosa y tirantes finos trenzados.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000003',
        'Verso',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_007.jpg',
        'assets/images/romance/MARIA_DIEZMA_007.jpg',
        'assets/images/romance/MARIA_DIEZMA_008.jpg',
        'assets/images/romance/MARIA_DIEZMA_009.jpg',
        'Elegante corte recto en crepé de seda con botonadura forrada y cola capilla.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000004',
        'Lírica',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_010.jpg',
        'assets/images/romance/MARIA_DIEZMA_010.jpg',
        'assets/images/romance/MARIA_DIEZMA_011.jpg',
        'assets/images/romance/MARIA_DIEZMA_012.jpg',
        'Cuerpo drapeado artesanal en gasa con escote cruzado y falda evasé con movimiento.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000005',
        'Etérea',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_013.jpg',
        'assets/images/romance/MARIA_DIEZMA_013.jpg',
        'assets/images/romance/MARIA_DIEZMA_014.jpg',
        'assets/images/romance/MARIA_DIEZMA_015.jpg',
        'Vestido vaporoso con mangas acampanadas en tul ilusión y delicados detalles botánicos.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000006',
        'Poema',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_016.jpg',
        'assets/images/romance/MARIA_DIEZMA_016.jpg',
        'assets/images/romance/MARIA_DIEZMA_017.jpg',
        'assets/images/romance/MARIA_DIEZMA_018.jpg',
        'Romántica silueta en A con aplicaciones florales en relieve y espalda descubierta.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000007',
        'Rima',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_019.jpg',
        'assets/images/romance/MARIA_DIEZMA_019.jpg',
        'assets/images/romance/MARIA_DIEZMA_020.jpg',
        'assets/images/romance/MARIA_DIEZMA_021.jpg',
        'Diseño minimalista con escote barco en satén duquesa y cintura marcada con lazo sutil.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000008',
        'Esencia',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_022.jpg',
        'assets/images/romance/MARIA_DIEZMA_022.jpg',
        'assets/images/romance/MARIA_DIEZMA_023.jpg',
        'assets/images/romance/MARIA_DIEZMA_024.jpg',
        'Confección pura y estilizada en crepé georgette con delicadas mangas francesas.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000009',
        'Sutil',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_025.jpg',
        'assets/images/romance/MARIA_DIEZMA_025.jpg',
        'assets/images/romance/MARIA_DIEZMA_026.jpg',
        'assets/images/romance/MARIA_DIEZMA_027.jpg',
        'Vestido midi con falda de suave vuelo y escote corazón ribeteado en ondas.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000010',
        'Eterna',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_028.jpg',
        'assets/images/romance/MARIA_DIEZMA_028.jpg',
        'assets/images/romance/MARIA_DIEZMA_029.jpg',
        'assets/images/romance/MARIA_DIEZMA_030.jpg',
        'Corte clásico atemporal con majestuosa cola desmontable y cuerpo de jacquard floral.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000011',
        'Flora',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_031.jpg',
        'assets/images/romance/MARIA_DIEZMA_031.jpg',
        'assets/images/romance/MARIA_DIEZMA_032.jpg',
        'assets/images/romance/MARIA_DIEZMA_033.jpg',
        'Bordados de inspiración botánica hechos a mano sobre base de chantilly y tul sedoso.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000012',
        'Seda',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_034.jpg',
        'assets/images/romance/MARIA_DIEZMA_034.jpg',
        'assets/images/romance/MARIA_DIEZMA_035.jpg',
        'assets/images/romance/MARIA_DIEZMA_036.jpg',
        'Diseño lencero de satén de seda cortado al bies con caída fluida y espalda cruzada.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000013',
        'Lino',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_037.jpg',
        'assets/images/romance/MARIA_DIEZMA_037.jpg',
        'assets/images/romance/MARIA_DIEZMA_038.jpg',
        'assets/images/romance/MARIA_DIEZMA_039.jpg',
        'Elegancia contemporánea en lino rústico sedoso con botonadura artesanal de nácar.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000014',
        'Nácar',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_040.jpg',
        'assets/images/romance/MARIA_DIEZMA_040.jpg',
        'assets/images/romance/MARIA_DIEZMA_041.jpg',
        'assets/images/romance/MARIA_DIEZMA_042.jpg',
        'Refinado diseño con destellos nacarados sobre encaje bordado y escote halter.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000015',
        'Magnolia',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_043.jpg',
        'assets/images/romance/MARIA_DIEZMA_043.jpg',
        'assets/images/romance/MARIA_DIEZMA_044.jpg',
        'assets/images/romance/MARIA_DIEZMA_045.jpg',
        'Diseño de corte sirena con bordados florales artesanales en tul y escote corazón.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000016',
        'Azahar',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_046.jpg',
        'assets/images/romance/MARIA_DIEZMA_046.jpg',
        'assets/images/romance/MARIA_DIEZMA_047.jpg',
        'assets/images/romance/MARIA_DIEZMA_048.jpg',
        'Elegante confección en encaje chantilly y espalda baja con botonadura artesanal.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000017',
        'Organza',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_049.jpg',
        'assets/images/romance/MARIA_DIEZMA_049.jpg',
        'assets/images/romance/MARIA_DIEZMA_050.jpg',
        'assets/images/romance/MARIA_DIEZMA_051.jpg',
        'Volúmenes esculturales con falda en capas de organza de seda y cuerpo entallado.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000018',
        'Rocío',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_052.jpg',
        'assets/images/romance/MARIA_DIEZMA_052.jpg',
        'assets/images/romance/MARIA_DIEZMA_053.jpg',
        'assets/images/romance/MARIA_DIEZMA_054.jpg',
        'Delicado diseño con microperlas cosidas a mano en el escote y falda de gasa pura.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000019',
        'Pétalo',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_055.jpg',
        'assets/images/romance/MARIA_DIEZMA_055.jpg',
        'assets/images/romance/MARIA_DIEZMA_056.jpg',
        'assets/images/romance/MARIA_DIEZMA_057.jpg',
        'Falda con aplicaciones florales tridimensionales que simulan una lluvia de pétalos.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000020',
        'Sauce',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_058.jpg',
        'assets/images/romance/MARIA_DIEZMA_058.jpg',
        'assets/images/romance/MARIA_DIEZMA_059.jpg',
        'assets/images/romance/MARIA_DIEZMA_060.jpg',
        'Líneas caídas y mangas fluidas que evocan ligereza natural en tul de seda.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000021',
        'Encaje',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_061.jpg',
        'assets/images/romance/MARIA_DIEZMA_061.jpg',
        'assets/images/romance/MARIA_DIEZMA_062.jpg',
        'assets/images/romance/MARIA_DIEZMA_063.jpg',
        'Patrón clásico recubierto en exquisito encaje guipur con mangas largas transparentes.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000022',
        'Puntilla',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_064.jpg',
        'assets/images/romance/MARIA_DIEZMA_064.jpg',
        'assets/images/romance/MARIA_DIEZMA_065.jpg',
        'assets/images/romance/MARIA_DIEZMA_066.jpg',
        'Detalles de puntilla artesanal en cuello y puños con sobrio corte imperio.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000023',
        'Trama',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_067.jpg',
        'assets/images/romance/MARIA_DIEZMA_067.jpg',
        'assets/images/romance/MARIA_DIEZMA_068.jpg',
        'assets/images/romance/MARIA_DIEZMA_069.jpg',
        'Textura estructurada en mikado ligero con pliegues arquitectónicos impecables.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000024',
        'Bordado',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_070.jpg',
        'assets/images/romance/MARIA_DIEZMA_070.jpg',
        'assets/images/romance/MARIA_DIEZMA_071.jpg',
        'assets/images/romance/MARIA_DIEZMA_072.jpg',
        'Bordados florales en hilo de seda sobre torso transparente y falda con vuelo suave.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000025',
        'Pliegue',
        'Romance',
        'assets/images/romance/MARIA_DIEZMA_073.jpg',
        'assets/images/romance/MARIA_DIEZMA_073.jpg',
        'assets/images/romance/MARIA_DIEZMA_074.jpg',
        'assets/images/romance/MARIA_DIEZMA_075.jpg',
        'Cuerpo plisado soleil a mano con abertura lateral elegante y sofisticada caída.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),

    -- -------------------------------------------------------------
    -- Colección 2: Nayade de Gala (assets/images/nayade/)
    -- -------------------------------------------------------------
    (
        '30000000-0000-0000-0000-000000000026',
        'Manantial',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_001.jpg',
        'assets/images/nayade/MARIA_DIEZMA_001.jpg',
        'assets/images/nayade/MARIA_DIEZMA_002.jpg',
        'assets/images/nayade/MARIA_DIEZMA_003.jpg',
        'Diseño fluido de gasa de seda con caída natural y delicado escote en V.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000027',
        'Bruma Marina',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_004.jpg',
        'assets/images/nayade/MARIA_DIEZMA_004.jpg',
        'assets/images/nayade/MARIA_DIEZMA_005.jpg',
        'assets/images/nayade/MARIA_DIEZMA_006.jpg',
        'Confeccionado en crepé vaporoso con mangas etéreas y detalles en sutil brillo plateado.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000028',
        'Reflejo de Cristal',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_007.jpg',
        'assets/images/nayade/MARIA_DIEZMA_007.jpg',
        'assets/images/nayade/MARIA_DIEZMA_008.jpg',
        'assets/images/nayade/MARIA_DIEZMA_009.jpg',
        'Silueta estilizada con bordados en microcristales y espalda descubierta con pedrería fina.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000029',
        'Gota de Rocío',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_010.jpg',
        'assets/images/nayade/MARIA_DIEZMA_010.jpg',
        'assets/images/nayade/MARIA_DIEZMA_011.jpg',
        'assets/images/nayade/MARIA_DIEZMA_012.jpg',
        'Vestido midi con falda drapeada soleil y delicados tirantes joya espagueti.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000030',
        'Cascada de Seda',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_013.jpg',
        'assets/images/nayade/MARIA_DIEZMA_013.jpg',
        'assets/images/nayade/MARIA_DIEZMA_014.jpg',
        'assets/images/nayade/MARIA_DIEZMA_015.jpg',
        'Diseño de fiesta de corte imperio en satén fluido con sobrefalda en cascada.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000031',
        'Oleaje Eterno',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_016.jpg',
        'assets/images/nayade/MARIA_DIEZMA_016.jpg',
        'assets/images/nayade/MARIA_DIEZMA_017.jpg',
        'assets/images/nayade/MARIA_DIEZMA_018.jpg',
        'Patrón estructurado con drapeado asimétrico en cuerpo y volante ondeado lateral.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000032',
        'Espuma de Plata',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_019.jpg',
        'assets/images/nayade/MARIA_DIEZMA_019.jpg',
        'assets/images/nayade/MARIA_DIEZMA_020.jpg',
        'assets/images/nayade/MARIA_DIEZMA_021.jpg',
        'Tejido ligero en tul bordado con hilo metálico plateado y falda evasé con movimiento.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000033',
        'Corriente Azul',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_022.jpg',
        'assets/images/nayade/MARIA_DIEZMA_022.jpg',
        'assets/images/nayade/MARIA_DIEZMA_023.jpg',
        'assets/images/nayade/MARIA_DIEZMA_024.jpg',
        'Línea contemporánea en satén sedoso con abertura lateral y caída envolvente.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000034',
        'Marea Calma',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_025.jpg',
        'assets/images/nayade/MARIA_DIEZMA_025.jpg',
        'assets/images/nayade/MARIA_DIEZMA_026.jpg',
        'assets/images/nayade/MARIA_DIEZMA_027.jpg',
        'Elegante diseño minimalista en crepé satinado con escote barco y botonadura artesanal.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000035',
        'Remolino de Tul',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_028.jpg',
        'assets/images/nayade/MARIA_DIEZMA_028.jpg',
        'assets/images/nayade/MARIA_DIEZMA_029.jpg',
        'assets/images/nayade/MARIA_DIEZMA_030.jpg',
        'Falda volumétrica en capas de tul de seda superpuestas con cuerpo entallado.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000036',
        'Brisa Fluvial',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_031.jpg',
        'assets/images/nayade/MARIA_DIEZMA_031.jpg',
        'assets/images/nayade/MARIA_DIEZMA_032.jpg',
        'assets/images/nayade/MARIA_DIEZMA_033.jpg',
        'Confección ligera en georgette de seda con escote halter cruzado y espalda al aire.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000037',
        'Espejo de Agua',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_034.jpg',
        'assets/images/nayade/MARIA_DIEZMA_034.jpg',
        'assets/images/nayade/MARIA_DIEZMA_035.jpg',
        'assets/images/nayade/MARIA_DIEZMA_036.jpg',
        'Vestido de gala en raso brillante de acabado líquido con drapeado cruzado frontal.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000038',
        'Cascada de Luz',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_037.jpg',
        'assets/images/nayade/MARIA_DIEZMA_037.jpg',
        'assets/images/nayade/MARIA_DIEZMA_038.jpg',
        'assets/images/nayade/MARIA_DIEZMA_039.jpg',
        'Bordados de destellos dorados artesanales sobre base translúcida y mangas capa desestructuradas.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000039',
        'Onda Serena',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_040.jpg',
        'assets/images/nayade/MARIA_DIEZMA_040.jpg',
        'assets/images/nayade/MARIA_DIEZMA_041.jpg',
        'assets/images/nayade/MARIA_DIEZMA_042.jpg',
        'Silueta columna con refinado corte al bies y escote asimétrico en tono cálido.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000040',
        'Murmullo de Río',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_043.jpg',
        'assets/images/nayade/MARIA_DIEZMA_043.jpg',
        'assets/images/nayade/MARIA_DIEZMA_044.jpg',
        'assets/images/nayade/MARIA_DIEZMA_045.jpg',
        'Cuerpo drapeado a mano en lino sedoso con tirantes finos y falda de vuelo suave.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000041',
        'Rocío del Alba',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_046.jpg',
        'assets/images/nayade/MARIA_DIEZMA_046.jpg',
        'assets/images/nayade/MARIA_DIEZMA_047.jpg',
        'assets/images/nayade/MARIA_DIEZMA_048.jpg',
        'Diseño romántico con aplicaciones perladas en el torso y falda fluida en gasa multicapa.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000042',
        'Salitre y Oro',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_049.jpg',
        'assets/images/nayade/MARIA_DIEZMA_049.jpg',
        'assets/images/nayade/MARIA_DIEZMA_050.jpg',
        'assets/images/nayade/MARIA_DIEZMA_051.jpg',
        'Espectacular confección en brocado de fiesta con hilos dorados y textura tornasolada.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000043',
        'Manantial de Ensueño',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_052.jpg',
        'assets/images/nayade/MARIA_DIEZMA_052.jpg',
        'assets/images/nayade/MARIA_DIEZMA_053.jpg',
        'assets/images/nayade/MARIA_DIEZMA_054.jpg',
        'Silueta princesa desestructurada con escote corazón y falda etérea de ensueño.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000044',
        'Lluvia de Estrellas',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_055.jpg',
        'assets/images/nayade/MARIA_DIEZMA_055.jpg',
        'assets/images/nayade/MARIA_DIEZMA_056.jpg',
        'assets/images/nayade/MARIA_DIEZMA_057.jpg',
        'Pedrería brillante distribuida en degradé de hombros a falda sobre tul de seda.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000045',
        'Aegle',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_058.jpg',
        'assets/images/nayade/MARIA_DIEZMA_058.jpg',
        'assets/images/nayade/MARIA_DIEZMA_059.jpg',
        'assets/images/nayade/MARIA_DIEZMA_060.jpg',
        'Diseño mitológico en seda salvaje con pliegues esculturales y fajín satinado.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000046',
        'Calírroe',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_061.jpg',
        'assets/images/nayade/MARIA_DIEZMA_061.jpg',
        'assets/images/nayade/MARIA_DIEZMA_062.jpg',
        'assets/images/nayade/MARIA_DIEZMA_063.jpg',
        'Corte helénico en muselina fluida con broche joya en el hombro y espalda infinita.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000047',
        'Castalia',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_064.jpg',
        'assets/images/nayade/MARIA_DIEZMA_064.jpg',
        'assets/images/nayade/MARIA_DIEZMA_065.jpg',
        'assets/images/nayade/MARIA_DIEZMA_066.jpg',
        'Patrón clásico depurado en crepé pesado con majestuosa cola fluida integrada.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000048',
        'Creusa',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_067.jpg',
        'assets/images/nayade/MARIA_DIEZMA_067.jpg',
        'assets/images/nayade/MARIA_DIEZMA_068.jpg',
        'assets/images/nayade/MARIA_DIEZMA_069.jpg',
        'Elegante confección midi con escote palabra de honor y sobrefalda desmontable.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000049',
        'Dafne',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_070.jpg',
        'assets/images/nayade/MARIA_DIEZMA_070.jpg',
        'assets/images/nayade/MARIA_DIEZMA_071.jpg',
        'assets/images/nayade/MARIA_DIEZMA_072.jpg',
        'Inspiración orgánica con bordados de hojas y motivos botánicos en hilo de seda.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000050',
        'Lilaia',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_073.jpg',
        'assets/images/nayade/MARIA_DIEZMA_073.jpg',
        'assets/images/nayade/MARIA_DIEZMA_074.jpg',
        'assets/images/nayade/MARIA_DIEZMA_075.jpg',
        'Diseño evasé en tono sutil con mangas abullonadas desmontables y escote cuadrado.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000051',
        'Melite',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_076.jpg',
        'assets/images/nayade/MARIA_DIEZMA_076.jpg',
        'assets/images/nayade/MARIA_DIEZMA_077.jpg',
        'assets/images/nayade/MARIA_DIEZMA_078.jpg',
        'Silueta sirena confeccionada en satén líquido con abertura pronunciada y tirantes cruzados.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000052',
        'Nomia',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_079.jpg',
        'assets/images/nayade/MARIA_DIEZMA_079.jpg',
        'assets/images/nayade/MARIA_DIEZMA_080.jpg',
        'assets/images/nayade/MARIA_DIEZMA_081.jpg',
        'Línea depurada y minimalista con cuerpo encorsetado y falda con vuelo envolvente.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000053',
        'Peribea',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_082.jpg',
        'assets/images/nayade/MARIA_DIEZMA_082.jpg',
        'assets/images/nayade/MARIA_DIEZMA_083.jpg',
        'assets/images/nayade/MARIA_DIEZMA_084.jpg',
        'Vestido de gala con drapeados diagonales y mangas capa en tul ilusión.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000054',
        'Salmacis',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_085.jpg',
        'assets/images/nayade/MARIA_DIEZMA_085.jpg',
        'assets/images/nayade/MARIA_DIEZMA_086.jpg',
        'assets/images/nayade/MARIA_DIEZMA_087.jpg',
        'Exclusivo diseño en mikado y gasa con escote halter y espalda esculpida.',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        '30000000-0000-0000-0000-000000000055',
        'Egina',
        'Nayade de Gala',
        'assets/images/nayade/MARIA_DIEZMA_088.jpg',
        'assets/images/nayade/MARIA_DIEZMA_088.jpg',
        'assets/images/nayade/MARIA_DIEZMA_089.jpg',
        'assets/images/nayade/MARIA_DIEZMA_090.jpg',
        'Exclusivo diseño en mikado y gasa con escote halter y espalda esculpida.',
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

-- ==========================================================
-- 5. PRESS_ARTICLES (5 artículos de prensa sobre María Diezma)
-- ==========================================================
INSERT INTO press_articles (id, magazine_name, publication_date, title, description, article_url, created_at, updated_at)
VALUES
    (
        '50000000-0000-0000-0000-000000000001',
        'El Español',
        'Marzo 2026',
        'Arranca en Albacete "CLM es Moda" con los desfiles de Félix Ramiro',
        'Este lunes ha tenido lugar la inauguración oficial de la III edición de CLM es Moda en la Fábrica de Harinas de Albacete de la exposición "Materia y Moda. De los oficios artesanos a la moda contemporánea", una muestra que pone en diálogo la tradición artesanal y la creación contemporánea, acercando al público el valor de los oficios y su influencia en el diseño actual.',
        'https://www.elespanol.com/eldigitalcastillalamancha/region/albacete/20260525/arranca-albacete-clm-moda-desfiles-felix-ramiro-raquel-lopez-carmen-alba-maria-diezma/1003744259085_0.html',
        CURRENT_TIMESTAMP - INTERVAL '120 days',
        CURRENT_TIMESTAMP - INTERVAL '120 days'
    ),
    (
        '50000000-0000-0000-0000-000000000002',
        'Lucia Se Casa',
        'Septiembre 2021',
        'La sencillez y el corte clásico son signos de elegancia',
        'La diseñadora María Diezma lleva la costura en sus genes. Heredera de las técnicas de su madre y de su abuela, manejaba las agujas desde muy temprana edad, y ya en su niñez disfrutaba bordando con bastidor, hilvanando o rematando sus diseños',
        'https://luciasecasa.com/novia/vestidos-de-novia/protagonistas-maria-diezma-la-sencillez-y-el-corte-clasico-son-signos-de-elegancia/',
        CURRENT_TIMESTAMP - INTERVAL '180 days',
        CURRENT_TIMESTAMP - INTERVAL '180 days'
    )
ON CONFLICT (id) DO NOTHING;

