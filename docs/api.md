# UTNianos API

API REST que envuelve el contenido del foro UTNianos (MyBB), devolviendo los datos en formato JSON en lugar del HTML original.

## Base URL

```
http://localhost:8080
```

## Endpoints

### 1. Listado de Foros

```
GET /foros
```

Devuelve todas las categorías y foros disponibles.

**Respuesta:**
```json
{
  "forums": [
    {
      "name": "Carreras de Grado",
      "forums": [
        {
          "id": 93,
          "name": "Sistemas",
          "description": "Foro de Sistemas",
          "url": "https://www.utnianos.com.ar/foro/foro-sistemas",
          "slug": "sistemas"
        }
      ]
    }
  ]
}
```

---

### 2. Ver Foro

```
GET /foro/{slug}
GET /foro/{slug}?page=2
GET /foro/{slug}?aporte=1&materia=2
```

Parámetros de URL:

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `slug` | path | Identificador del foro (ej: `sistemas`) |
| `page` | query | Número de página |
| `search` | query | Texto de búsqueda en títulos |
| `aporte` | query | Filtro por tipo de aporte (múltiple) |
| `materia` | query | Filtro por materia (múltiple) |
| `sortby` | query | Campo de ordenamiento |
| `order` | query | Dirección (`asc` o `desc`) |

**Valores: `aporte`**

| ID | Nombre |
|----|--------|
| 1 | Parciales |
| 2 | Finales |
| 3 | Trabajo practico |
| 4 | Apuntes y Guias |
| 5 | Libro |
| 6 | Profesores |
| 7 | Ejercicios |
| 8 | Dudas y recomendaciones |
| 9 | Consultas administrativas |
| 10 | Otro |
| 11 | Guias CEIT |

**Valores: `materia`** (ver documentación completa en `docs/utnianos-api.md`)

| ID | Nombre |
|----|--------|
| 2 | Algoritmos y Estructuras de Datos |
| 7 | Análisis Matemático I |
| 8 | Paradigmas de Programación |
| 9 | Sistemas Operativos |
| 10 | Análisis de Sistemas |
| 15 | Diseño de Sistemas |
| ... | (más materias) |

**Respuesta:**
```json
{
  "topics": [
    {
      "id": 39447,
      "title": "[APORTE] Parciales de la cursada de DdSI 2025",
      "url": "https://www.utnianos.com.ar/foro/tema-aporte-parciales-ddsi-2025--39447",
      "author": "chilaverto",
      "author_url": "https://www.utnianos.com.ar/foro/usuario-chilaverto",
      "replies": 0,
      "views": 0,
      "last_post": "2026-03-06T13:23:00Z",
      "last_post_by": "chilaverto",
      "is_pinned": false,
      "is_closed": false,
      "thanked_count": 0
    }
  ],
  "total": 20,
  "page": 1,
  "total_pages": 240
}
```

**Ejemplos:**

```bash
# Foro de Sistemas
curl http://localhost:8080/foro/sistemas

# Parciales de Algoritmos
curl "http://localhost:8080/foro/sistemas?aporte=1&materia=2"

# Buscar en el foro
curl "http://localhost:8080/foro/sistemas?search=final"

# Segunda página
curl "http://localhost:8080/foro/sistemas?page=2"

# Múltiples filtros
curl "http://localhost:8080/foro/sistemas?aporte=1&aporte=2&materia=2"
```

---

### 3. Ver Tema

```
GET /tema/{slug}
```

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `slug` | path | Identificador del tema |

**Respuesta:**
```json
{
  "id": 39443,
  "title": "Final ADS 03/03/2026",
  "url": "https://www.utnianos.com.ar/foro/tema-final-ads-03-03-2026--39443",
  "author": "Sebastian1",
  "author_url": "https://www.utnianos.com.ar/foro/usuario-Sebastian1",
  "replies": 0,
  "views": 0,
  "last_post": "2026-03-04T10:20:00Z",
  "last_post_by": "Sebastian1",
  "content": "Contenido del primer post...",
  "posts": [
    {
      "id": 0,
      "author": "Sebastian1",
      "author_url": "https://www.utnianos.com.ar/foro/usuario-Sebastian1",
      "content": "Contenido del post...",
      "date": "2026-03-04T10:20:00Z",
      "thanked": 0,
      "number": 1
    }
  ],
  "is_pinned": false,
  "is_closed": false
}
```

---

### 4. Búsqueda

```
GET /search
GET /search?action=getnew
GET /search?action=unanswered
GET /search?q=parcial&fids=93
GET /search?keywords=final&action=do_search
```

Parámetros de query:

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `q` | query | Palabras clave |
| `keywords` | query | Palabras clave (alternativo) |
| `action` | query | Acción: `do_search`, `getnew`, `getdaily`, `unanswered` |
| `fids` | query | ID del foro(s), separado por coma |
| `forums[]` | query | Foros a buscar |
| `postthread` | query | 1 = buscar todo, 2 = solo títulos |
| `sortby` | query | Orden: `lastpost`, `starter`, `forum` |
| `sortordr` | query | `asc` o `desc` |

**Respuesta:**
```json
{
  "topics": [
    {
      "id": 39447,
      "title": "[APORTE] Parciales de la cursada de DdSI 2025",
      "url": "...",
      "author": "chilaverto",
      "replies": 0,
      "views": 0,
      "last_post": "2026-03-06T13:23:00Z",
      "last_post_by": "chilaverto",
      "is_pinned": false,
      "is_closed": false
    }
  ],
  "total": 20,
  "page": 1,
  "total_pages": 0
}
```

**Ejemplos:**

```bash
# Temas nuevos
curl "http://localhost:8080/search?action=getnew&fids=93"

# Temas sin respuesta
curl "http://localhost:8080/search?action=unanswered&fids=93"

# Buscar con keywords
curl "http://localhost:8080/search?keywords=parcial&fids=93"
```

---

### 5. Perfil de Usuario

```
GET /usuario/{username}
```

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `username` | path | Nombre de usuario |

**Respuesta:**
```json
{
  "id": 0,
  "username": "Dem0",
  "join_date": "...",
  "posts": 1234,
  "location": "...",
  "website": "https://...",
  "bio": "...",
  "avatar": "https://...",
  "signature": "..."
}
```

---

### 6. Raíz

```
GET /
```

Devuelve información básica de la API.

---

## CORS

La API soporta CORS para requests desde cualquier origen.

## Errores

Los errores se devuelven en formato:

```json
{
  "error": "Error description",
  "message": "Detailed error message"
}
```

Códigos de estado HTTP:
- `200` - OK
- `400` - Bad Request
- `404` - Not Found
- `500` - Internal Server Error

## Notas

- El sitio original (utnianos.com.ar) puede bloquear requests excesivos. Usar con moderación.
- Algunos campos pueden estar vacíos si el sitio cambia su estructura HTML.
- La API hace scraping del HTML del foro, por lo que es sensible a cambios en la estructura del sitio.
