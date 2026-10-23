-- ==============================================================================
-- Migration: 000003_add_secretos_collection.sql
-- Histórico de versiones: Inserción de la nueva colección "Secretos" y sus 15 vestidos.
--
-- Estructura de imágenes (001 a 047):
--   - 001 a 045: 15 vestidos con 3 imágenes consecutivas cada uno (1 portada/imagen1 + imagen2 + imagen3).
--   - 046: Fotografía de la diseñadora (assets/images/secretos/MARIA_DIEZMA_046.jpg - uso directo en frontend/perfil).
--   - 047: Portada de la colección Secretos (assets/images/secretos/MARIA_DIEZMA_047.jpg - registrada en tabla collections).
--
-- Descripciones: Dejadas temporalmente en blanco ('') según especificación.
-- ==============================================================================

-- ==========================================================
-- 1. COLECCIÓN: Secretos
-- ==========================================================
INSERT INTO collections (id, name, image_path, description, created_at, updated_at)
VALUES
    (
        '20000000-0000-0000-0000-000000000003',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_047.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (id) DO UPDATE 
SET 
    name = EXCLUDED.name,
    image_path = EXCLUDED.image_path,
    description = CASE WHEN collections.description = '' THEN EXCLUDED.description ELSE collections.description END,
    updated_at = CURRENT_TIMESTAMP;

-- ==========================================================
-- 2. VESTIDOS: 15 Vestidos de la Colección Secretos (001 - 045)
-- ==========================================================
INSERT INTO dresses (id, name, collection, image_path, image1_path, image2_path, image3_path, description, created_at, updated_at)
VALUES
    -- 1. Confidencia (001 - 003)
    (
        '30000000-0000-0000-0000-000000000056',
        'Confidencia',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_001.jpg',
        'assets/images/secretos/MARIA_DIEZMA_001.jpg',
        'assets/images/secretos/MARIA_DIEZMA_002.jpg',
        'assets/images/secretos/MARIA_DIEZMA_003.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 2. Promesa (004 - 006)
    (
        '30000000-0000-0000-0000-000000000057',
        'Promesa',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_004.jpg',
        'assets/images/secretos/MARIA_DIEZMA_004.jpg',
        'assets/images/secretos/MARIA_DIEZMA_005.jpg',
        'assets/images/secretos/MARIA_DIEZMA_006.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 3. Suspiro (007 - 009)
    (
        '30000000-0000-0000-0000-000000000058',
        'Suspiro',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_007.jpg',
        'assets/images/secretos/MARIA_DIEZMA_007.jpg',
        'assets/images/secretos/MARIA_DIEZMA_008.jpg',
        'assets/images/secretos/MARIA_DIEZMA_009.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 4. Alba (010 - 012)
    (
        '30000000-0000-0000-0000-000000000059',
        'Alba',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_010.jpg',
        'assets/images/secretos/MARIA_DIEZMA_010.jpg',
        'assets/images/secretos/MARIA_DIEZMA_011.jpg',
        'assets/images/secretos/MARIA_DIEZMA_012.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 5. Velo (013 - 015)
    (
        '30000000-0000-0000-0000-000000000060',
        'Velo',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_013.jpg',
        'assets/images/secretos/MARIA_DIEZMA_013.jpg',
        'assets/images/secretos/MARIA_DIEZMA_014.jpg',
        'assets/images/secretos/MARIA_DIEZMA_015.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 6. Susurro (016 - 018)
    (
        '30000000-0000-0000-0000-000000000061',
        'Susurro',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_016.jpg',
        'assets/images/secretos/MARIA_DIEZMA_016.jpg',
        'assets/images/secretos/MARIA_DIEZMA_017.jpg',
        'assets/images/secretos/MARIA_DIEZMA_018.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 7. Íntima (019 - 021)
    (
        '30000000-0000-0000-0000-000000000062',
        'Íntima',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_019.jpg',
        'assets/images/secretos/MARIA_DIEZMA_019.jpg',
        'assets/images/secretos/MARIA_DIEZMA_020.jpg',
        'assets/images/secretos/MARIA_DIEZMA_021.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 8. Silencio (022 - 024)
    (
        '30000000-0000-0000-0000-000000000063',
        'Silencio',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_022.jpg',
        'assets/images/secretos/MARIA_DIEZMA_022.jpg',
        'assets/images/secretos/MARIA_DIEZMA_023.jpg',
        'assets/images/secretos/MARIA_DIEZMA_024.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 9. Hechizo (025 - 027)
    (
        '30000000-0000-0000-0000-000000000064',
        'Hechizo',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_025.jpg',
        'assets/images/secretos/MARIA_DIEZMA_025.jpg',
        'assets/images/secretos/MARIA_DIEZMA_026.jpg',
        'assets/images/secretos/MARIA_DIEZMA_027.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 10. Enigma (028 - 030)
    (
        '30000000-0000-0000-0000-000000000065',
        'Enigma',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_028.jpg',
        'assets/images/secretos/MARIA_DIEZMA_028.jpg',
        'assets/images/secretos/MARIA_DIEZMA_029.jpg',
        'assets/images/secretos/MARIA_DIEZMA_030.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 11. Olivo (031 - 033)
    (
        '30000000-0000-0000-0000-000000000066',
        'Olivo',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_031.jpg',
        'assets/images/secretos/MARIA_DIEZMA_031.jpg',
        'assets/images/secretos/MARIA_DIEZMA_032.jpg',
        'assets/images/secretos/MARIA_DIEZMA_033.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 12. Áurea (034 - 036)
    (
        '30000000-0000-0000-0000-000000000067',
        'Áurea',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_034.jpg',
        'assets/images/secretos/MARIA_DIEZMA_034.jpg',
        'assets/images/secretos/MARIA_DIEZMA_035.jpg',
        'assets/images/secretos/MARIA_DIEZMA_036.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 13. Sombra (037 - 039)
    (
        '30000000-0000-0000-0000-000000000068',
        'Sombra',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_037.jpg',
        'assets/images/secretos/MARIA_DIEZMA_037.jpg',
        'assets/images/secretos/MARIA_DIEZMA_038.jpg',
        'assets/images/secretos/MARIA_DIEZMA_039.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 14. Penumbra (040 - 042)
    (
        '30000000-0000-0000-0000-000000000069',
        'Penumbra',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_040.jpg',
        'assets/images/secretos/MARIA_DIEZMA_040.jpg',
        'assets/images/secretos/MARIA_DIEZMA_041.jpg',
        'assets/images/secretos/MARIA_DIEZMA_042.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    -- 15. Relicario (043 - 045)
    (
        '30000000-0000-0000-0000-000000000070',
        'Relicario',
        'Secretos',
        'assets/images/secretos/MARIA_DIEZMA_043.jpg',
        'assets/images/secretos/MARIA_DIEZMA_043.jpg',
        'assets/images/secretos/MARIA_DIEZMA_044.jpg',
        'assets/images/secretos/MARIA_DIEZMA_045.jpg',
        '',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (id) DO UPDATE
SET
    name = EXCLUDED.name,
    collection = EXCLUDED.collection,
    image_path = EXCLUDED.image_path,
    image1_path = EXCLUDED.image1_path,
    image2_path = EXCLUDED.image2_path,
    image3_path = EXCLUDED.image3_path,
    description = CASE WHEN dresses.description = '' THEN EXCLUDED.description ELSE dresses.description END,
    updated_at = CURRENT_TIMESTAMP;

-- ==============================================================================
-- NOTA COMPLEMENTARIA SOBRE LAS IMÁGENES:
-- - assets/images/secretos/MARIA_DIEZMA_046.jpg: Imagen de la diseñadora (sin inserción en dresses/collections, disponible para frontend).
-- - assets/images/secretos/MARIA_DIEZMA_047.jpg: Imagen principal de la colección Secretos (registrada en collections).
-- ==============================================================================
