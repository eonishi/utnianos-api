# utnianos-api

API REST en Go que expone datos del foro [utnianos.com.ar/foro](https://www.utnianos.com.ar/foro/) mediante web scraping.

## Tecnologias

- **Go 1.22**
- [chi](https://github.com/go-chi/chi) - Router HTTP
- [goquery](https://github.com/PuerkitoBio/goquery) - Scraping HTML (estilo jQuery)

## Estructura del proyecto

```
cmd/
  api/
    main.go           # Punto de entrada, configuracion del router
internal/
  handlers/
    handlers.go       # Handlers HTTP
  models/
    models.go         # Estructuras de datos (JSON responses)
  scraper/
    scraper.go        # Logica de scraping del foro
```

## Endpoints

| Metodo | Ruta               | Descripcion                                    |
|--------|--------------------|-------------------------------------------------|
| GET    | `/`                | Informacion general de la API                   |
| GET    | `/foros`           | Lista de foros con topic/message counts         |
| GET    | `/foro/{slug}`     | Temas de un foro (soporta filtros y paginacion) |
| GET    | `/tema/{slug}`     | Detalle de un tema con sus posts                |
| GET    | `/usuario/{username}` | Perfil de un usuario                         |
| GET    | `/search`          | Busqueda de temas                               |

### Parametros de `/foro/{slug}`

| Parametro | Descripcion                      |
|-----------|----------------------------------|
| `page`    | Numero de pagina                 |
| `aporte`  | Filtrar por tipo de aporte       |
| `materia` | Filtrar por materia              |

### Parametros de `/search`

| Parametro | Descripcion                      |
|-----------|----------------------------------|
| `q`       | Palabras clave                   |
| `keywords`| Alternativa a `q`                |
| `fids`    | ID del foro donde buscar         |
| `action`  | Accion de busqueda               |

## Instalacion y ejecucion

```bash
# Clonar repositorio
git clone <url-del-repo>
cd utnianos-api

# Instalar dependencias
go mod download

# Ejecutar
go run ./cmd/api/
```

El servidor inicia en `http://localhost:8080`.

### Desarrollo con hot reload

```bash
# Instalar Air
go install github.com/air-verse/air@latest

# Ejecutar
air
```

### Puerto personalizado

```bash
PORT=3000 go run ./cmd/api/
```

## Respuesta de ejemplo

```
GET /foros
```

```json
{
  "forums": [
    {
      "name": "Analisis Matematico I",
      "url": "https://www.utnianos.com.ar/foro/foro-analisis-matematico-i",
      "slug": "analisis-matematico-i",
      "topic_count": 1523,
      "message_count": 18420,
      "last_post": {
        "url": "https://...",
        "date": "2026-03-31T14:30:00Z",
        "author": "username",
        "author_url": "https://..."
      }
    }
  ]
}
```
