# MariaDiezmaBack - Backoffice REST API

Backend y API REST desarrollado en **Go 1.27** para el backoffice de MariaDiezma. Diseñado siguiendo los principios de **Clean Architecture** / **Hexagonal Architecture** y el estándar oficial de layout de proyectos en Go.

---

## 🚀 Características

- **Go 1.27.1**
- **Clean Architecture**: Separación estricta entre Dominio, Casos de Uso (Servicios), Repositorios y Controladores (Handlers).
- **Enrutamiento HTTP**: Basado en [`go-chi/chi/v5`](https://github.com/go-chi/chi), 100% compatible con la librería estándar `net/http`.
- **Autenticación & RBAC**: Tokens JWT (HMAC-SHA256), hashing de contraseñas con `bcrypt` y control de acceso por roles (`admin`, `manager`, `viewer`).
- **Persistencia Dual**:
  - **In-Memory**: Almacenamiento concurrente en memoria listo para funcionar sin dependencias externas (ideal para desarrollo rápido y tests).
  - **PostgreSQL**: Integración nativa de alto rendimiento con `pgx/v5` y pool de conexiones (`pgxpool`).
- **Observabilidad**: Logging estructurado con `log/slog` de Go, Request ID middleware y captura de pánicos (`Recovery`).
- **Gestión de Peticiones REST**: CRUD completo, filtrado por tipo, estado y prioridad, búsqueda textual, notas internas y paginación estructurada.
- **Docker & Compose**: Dockerfile multi-stage ligero y entorno `docker-compose` con base de datos PostgreSQL 17 y migraciones automáticas.

---

## 📁 Estructura del Proyecto

```plaintext
MariaDiezmaBack/
├── cmd/
│   └── api/
│       └── main.go                 # Punto de entrada, inyección de dependencias y Graceful Shutdown
├── internal/
│   ├── config/                     # Carga y validación de variables de entorno (.env)
│   │   └── config.go
│   ├── domain/                     # Entidades del dominio, DTOs y errores tipados
│   │   ├── user.go
│   │   ├── request.go              # Entidad RequestItem y filtros
│   │   └── errors.go
│   ├── handler/                    # Controladores HTTP (REST)
│   │   ├── auth_handler.go
│   │   ├── request_handler.go
│   │   └── health_handler.go
│   ├── middleware/                 # Middlewares (Auth JWT, CORS, Logger, Recovery)
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── repository/                 # Contratos de acceso a datos e implementaciones
│   │   ├── interfaces.go
│   │   ├── memory/                 # Implementación en memoria thread-safe
│   │   └── postgres/               # Implementación PostgreSQL con pgx/v5
│   └── service/                    # Lógica de negocio y casos de uso
│       ├── auth_service.go
│       └── request_service.go
├── pkg/
│   ├── response/                   # Respuestas JSON consistentes y paginación
│   └── validator/                  # Validaciones de entrada de datos
├── migrations/
│   └── 000001_init_schema.up.sql   # DDL de PostgreSQL
├── .env.example                    # Plantilla de variables de entorno
├── docker-compose.yml              # Orquestación de contenedores (API + PostgreSQL)
├── Dockerfile                      # Imagen multi-stage en Alpine
├── Makefile                        # Tareas y comandos de desarrollo
├── go.mod
└── go.sum
```

---

## ⚙️ Configuración (.env)

Copia `.env.example` a `.env`:

```bash
cp .env.example .env
```

| Variable | Descripción | Valor por defecto |
| :--- | :--- | :--- |
| `PORT` | Puerto del servidor HTTP | `8080` |
| `ENV` | Entorno (`development` o `production`) | `development` |
| `DB_DRIVER` | Controlador de BD (`memory` o `postgres`) | `memory` |
| `DATABASE_URL` | Cadena de conexión PostgreSQL | `postgres://postgres:postgres@localhost:5432/mariadiezma?sslmode=disable` |
| `JWT_SECRET` | Clave secreta para firmar los JWT | `super-secret-jwt-key-...` |
| `JWT_EXPIRATION_HOURS` | Tiempo de expiración del token | `24` |
| `ADMIN_EMAIL` | Email del usuario administrador inicial | `admin@mariadiezma.com` |
| `ADMIN_PASSWORD` | Contraseña del admin inicial | `AdminPass123!` |
| `CORS_ALLOWED_ORIGINS` | Orígenes permitidos (CORS) | `*` |

---

## 🛠️ Ejecución y Desarrollo

### 1. Ejecución local rápida (modo en memoria, sin DB externa)
```bash
make run
# o alternativamente:
go run ./cmd/api
```

El servidor arrancará en `http://localhost:8080` con el usuario inicial:
- **Email:** `admin@mariadiezma.com`
- **Contraseña:** `AdminPass123!`

### 2. Ejecución con Docker y PostgreSQL
```bash
make docker-up
```

Para detener los contenedores:
```bash
make docker-down
```

### 3. Inyección de Datos de Prueba (PostgreSQL Seed)
El proyecto incluye un script SQL con datos de prueba realistas para todas las tablas ([`migrations/000002_seed_data.sql`](file:///Users/oscar/Documents/Codes/Web/MariaDiezmaBack/migrations/000002_seed_data.sql)):
- **Sembrador nativo en Go (Recomendado, no requiere tener instalado `psql` en tu Mac)**:
  ```bash
  make seed
  ```
  *(Se conecta a PostgreSQL usando la variable `DATABASE_URL` del archivo `.env` y aplica el esquema y los datos mediante `pgx/v5`).*

- **Desde el contenedor de Docker (donde `psql` viene preinstalado)**:
  ```bash
  make seed-docker
  ```
  *(O automáticamente en el primer arranque al hacer `make docker-up`).*

### 4. Ejecución de Tests
```bash
make test
# O con cobertura:
make test-coverage
```

---

## 📡 Endpoints de la API

### Públicos
- `GET /health` - Estado de salud y uptime del servicio.
- `POST /api/v1/auth/login` - Inicio de sesión y obtención de JWT Bearer token.
- `GET /api/v1/colecciones` (o `/api/v1/collections`) - Obtiene el catálogo de colecciones para la web (nombre, ruta de imagen y descripción) desde base de datos.
- `GET /api/v1/vestidos` (o `/api/v1/dresses`) - Obtiene el catálogo de vestidos (nombre, colección y ruta de imagen). Admite filtro opcional por colección: `?coleccion=...`.
- `GET /api/v1/vestidos/detalle` (o `/api/v1/dresses/detail`) - Obtiene la ficha completa de un vestido (nombre, colección, ruta de imagen 1, ruta de imagen 2 y descripción) a partir de `?nombre=...&coleccion=...`.
- `POST /api/v1/requests` - Envío de peticiones/formularios generales desde la web pública.
- `POST /api/v1/citas` (o `/api/v1/appointments`) - Envío de citas solicitadas desde la web (guarda en backoffice y envía email con los datos de la cita).
- `GET /api/v1/prensa` (o `/api/v1/press`, `/api/v1/articulos-prensa`) - Obtiene el listado de artículos de prensa (nombre de la revista, fecha de publicación, titular, pequeña descripción y enlace del artículo).

### Protegidos (requieren cabecera `Authorization: Bearer <TOKEN>`)
- `GET /api/v1/auth/me` - Datos del usuario autenticado actual.
- `GET /api/v1/requests` - Listado paginado de peticiones (filtros: `?status=&priority=&type=&search=&page=&per_page=`).
- `GET /api/v1/requests/{id}` - Detalle de una petición.
- `PATCH /api/v1/requests/{id}/status` - Cambio de estado (`pending`, `in_progress`, `resolved`, `cancelled`) y notas internas.
- `PUT /api/v1/requests/{id}` - Edición general de la petición o asignación a operadores.
- `DELETE /api/v1/requests/{id}` - Eliminación de petición (Roles permitidos: `admin`, `manager`).

---

## 📝 Ejemplo de Flujo de Trabajo REST

### 1. Autenticación en el Backoffice
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@mariadiezma.com",
    "password": "AdminPass123!"
  }'
```

Respuesta:
```json
{
  "success": true,
  "message": "authenticated successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5c...",
    "expires_at": "2026-09-03T21:00:00Z",
    "user": {
      "id": "...",
      "email": "admin@mariadiezma.com",
      "name": "Administrator",
      "role": "admin"
    }
  }
}
```

### 2. Envío de petición desde el formulario web
```bash
curl -X POST http://localhost:8080/api/v1/requests \
  -H "Content-Type: application/json" \
  -d '{
    "type": "presupuesto",
    "priority": "high",
    "sender_name": "Juan Perez",
    "sender_email": "juan@example.com",
    "sender_phone": "+34 612 345 678",
    "subject": "Solicitud de información",
    "message": "Hola Maria, me gustaría saber precios y disponibilidad.",
    "metadata": {
      "source": "landing_page",
      "preferred_contact": "phone"
    }
  }'
```

### 3. Listado de peticiones en el Backoffice
```bash
curl -X GET "http://localhost:8080/api/v1/requests?status=pending&page=1&per_page=10" \
  -H "Authorization: Bearer <TU_TOKEN>"
```

### 4. Actualización del estado
```bash
curl -X PATCH http://localhost:8080/api/v1/requests/<REQUEST_ID>/status \
  -H "Authorization: Bearer <TU_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "internal_note": "Llamada realizada, pendiente de enviar dossier."
  }'
```

### 5. Envío de Solicitud de Cita desde la Web (Notificación por Email)

Disponible en `POST /api/v1/citas` (o su alias `POST /api/v1/appointments`).

#### Campos admitidos en el cuerpo de la petición (JSON):
- `tipo_cita` (`string`, obligatorio): Tipo de cita elegida (ej. `"Novia a medida"`, `"Invitada"`, `"Primera Consulta"`).
- `fecha` (`string`, obligatorio): Fecha seleccionada para la cita (`YYYY-MM-DD`).
- `franja_horaria` (`string`, obligatorio): Franja horaria elegida (ej. `"Tarde (16:00 - 19:00)"` o `"16:00 - 17:00"`).
- `nombre_apellidos` (`string`, obligatorio): Nombre y apellidos del contacto.
- `telefono_contacto` (`string`, obligatorio): Teléfono de contacto.
- `fecha_estimada` (`date` / `string` o `null`, opcional): Fecha estimada del evento o boda (puede ser `null` o `"YYYY-MM-DD"`).
- `detalles` (`string`, opcional): Comentarios o detalles específicos de la cita.
- `email` (`string`, opcional): Correo electrónico de contacto (si se envía, se valida el formato).

> *Nota: Por retrocompatibilidad, la API también admite los nombres de campos en inglés (`name`, `phone`, `date`, `time_slot`, `type`, `estimated_date`, `details`) y las variantes `tramo_horario`, `nombre`, `telefono`.*

#### Ejemplo de Petición:
```bash
curl -X POST http://localhost:8080/api/v1/citas \
  -H "Content-Type: application/json" \
  -d '{
    "tipo_cita": "Novia a medida",
    "fecha": "2026-10-25",
    "franja_horaria": "Tarde (16:00 - 19:00)",
    "nombre_apellidos": "Lucía Domínguez",
    "telefono_contacto": "+34 678 901 234",
    "fecha_estimada": "2027-05-15",
    "detalles": "Interesada en telas de seda natural y corte sirena"
  }'
```

*(Ejemplo con fecha estimada nula y sin email:)*
```bash
curl -X POST http://localhost:8080/api/v1/citas \
  -H "Content-Type: application/json" \
  -d '{
    "tipo_cita": "Madrina",
    "fecha": "2026-11-10",
    "franja_horaria": "Mañana (10:00 - 13:00)",
    "nombre_apellidos": "Carmen Navarro",
    "telefono_contacto": "+34 611 223 344",
    "fecha_estimada": null,
    "detalles": "Sin mangas"
  }'
```

Respuesta (201 Created):
```json
{
  "success": true,
  "message": "Cita solicitada correctamente",
  "data": {
    "id": "e67bfa37-88df-46fb-a0b2-bb3438adabec",
    "tipo_cita": "Novia a medida",
    "type": "Novia a medida",
    "fecha": "2026-10-25",
    "date": "2026-10-25",
    "franja_horaria": "Tarde (16:00 - 19:00)",
    "time_slot": "Tarde (16:00 - 19:00)",
    "nombre_apellidos": "Lucía Domínguez",
    "name": "Lucía Domínguez",
    "telefono_contacto": "+34 678 901 234",
    "phone": "+34 678 901 234",
    "fecha_estimada": "2027-05-15",
    "estimated_date": "2027-05-15",
    "detalles": "Interesada en telas de seda natural y corte sirena",
    "details": "Interesada en telas de seda natural y corte sirena",
    "status": "pending",
    "created_at": "2026-09-10T14:30:00Z"
  }
}
```
*Al recibir esta petición, los datos se guardan en el repositorio del backoffice (con fecha estimada y detalles en metadata) y se despacha un correo electrónico (HTML y texto plano con todos los datos) a la dirección configurada en `NOTIFICATION_EMAIL`.*

### 6. Consulta de Colecciones para la Web Frontend (Público)
```bash
curl -X GET http://localhost:8080/api/v1/colecciones
```

Respuesta (200 OK):
```json
{
  "success": true,
  "message": "Colecciones obtenidas correctamente",
  "data": [
    {
      "id": "20000000-0000-0000-0000-000000000001",
      "nombre": "Esencia Floral",
      "imagen": "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg",
      "descripcion": "Diseños inspirados en la delicadeza botánica y tonalidades primaverales."
    },
    {
      "id": "20000000-0000-0000-0000-000000000002",
      "nombre": "Atardecer Mediterráneo",
      "imagen": "assets/images/atardecer-mediterraneo/MARIA_DIEZMA_016.jpg",
      "descripcion": "Colección cálida con texturas fluidas y tonos terracota y dorados."
    }
  ]
}
```

### 7. Consulta de Vestidos para la Web Frontend (Público)
```bash
# Obtener todos los vestidos
curl -X GET http://localhost:8080/api/v1/vestidos

# O filtrar por colección
curl -X GET "http://localhost:8080/api/v1/vestidos?coleccion=Esencia+Floral"
```

Respuesta (200 OK):
```json
{
  "success": true,
  "message": "Vestidos obtenidos correctamente",
  "data": [
    {
      "id": "30000000-0000-0000-0000-000000000001",
      "nombre": "Vestido Magnolia",
      "coleccion": "Esencia Floral",
      "ruta_imagen": "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg"
    },
    {
      "id": "30000000-0000-0000-0000-000000000002",
      "nombre": "Vestido Jazmín",
      "coleccion": "Esencia Floral",
      "ruta_imagen": "assets/images/esencia-floral/MARIA_DIEZMA_004.jpg"
    }
  ]
}
```

### 8. Ficha Detallada de un Vestido (Público)
```bash
curl -X GET "http://localhost:8080/api/v1/vestidos/detalle?nombre=Vestido+Magnolia&coleccion=Esencia+Floral"
```

Respuesta (200 OK):
```json
{
  "success": true,
  "message": "Detalle del vestido obtenido correctamente",
  "data": {
    "id": "30000000-0000-0000-0000-000000000001",
    "nombre": "Vestido Magnolia",
    "coleccion": "Esencia Floral",
    "ruta_imagen_1": "assets/images/esencia-floral/MARIA_DIEZMA_001.jpg",
    "ruta_imagen_2": "assets/images/esencia-floral/MARIA_DIEZMA_002.jpg",
    "ruta_imagen_3": "assets/images/esencia-floral/MARIA_DIEZMA_003.jpg",
    "descripcion": "Vestido de corte sirena con bordados florales artesanales en tul y escote corazón."
  }
}
```

### 9. Consulta de Artículos de Prensa (Público)
```bash
curl -X GET http://localhost:8080/api/v1/prensa
```

Respuesta (200 OK):
```json
{
  "success": true,
  "message": "Artículos de prensa obtenidos correctamente",
  "data": [
    {
      "id": "press-001",
      "nombre_revista": "Vogue España",
      "fecha_publicacion": "2024-05-15",
      "titular": "María Diezma: La nueva era de la alta costura nupcial y la artesanía contemporánea",
      "descripcion": "Un recorrido íntimo por el atelier madrileño de María Diezma, donde cada puntada rinde homenaje a la tradición y al patronaje a medida.",
      "enlace_articulo": "https://www.vogue.es/novias/articulos/maria-diezma-alta-costura-nupcial",
      "enlace": "https://www.vogue.es/novias/articulos/maria-diezma-alta-costura-nupcial"
    }
  ]
}
```




