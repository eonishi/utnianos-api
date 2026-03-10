# Documentación API - UTNianos Forum

## Información General

**Base URL:** `https://www.utnianos.com.ar/foro/`

Este es un foro **MyBB** (versión antigua) de la comunidad de estudiantes de UTN (Universidad Tecnológica Nacional).

---

## Endpoints Principales

### 1. Listado de Foros
```
GET /
GET /index.php
```
Muestra el listado principal de todas las categorías y foros.

### 2. Ver Foro Específico
```
GET /foro-{nombre}
```
- **Ejemplo:** `/foro-sistemas` → Foro de Sistemas (fid=93)
- **URL larga:** `/forumdisplay.php?fid={id}`

### 3. Ver Tema/Thread
```
GET /tema-{slug}
```
- **Ejemplo:** `/tema-final-ads-03-03-2026`
- **URL larga:** `/showthread.php?tid={id}`

### 4. Búsqueda
```
GET /search.php
POST /search.php
```
- **Action:** `do_search` (POST)
- **Búsqueda rápida:** `getnew`, `getdaily`, `unanswered`

### 5. RSS/Atom Feed
```
GET /syndication.php
GET /syndication.php?fid={forum_id}
GET /syndication.php?type=atom1.0
```

### 6. Perfil de Usuario
```
GET /usuario-{username}
GET /member.php?action=profile&uid={uid}
```

---

## Parámetros de Query (para foros)

### Parámetros de Filtrado (Custom de UTNianos)

| Parámetro | Tipo | Descripción |
|------------|------|-------------|
| `search` | string | Texto de búsqueda en títulos |
| `filtertf_tipo_aporte[]` | int[] | **Tipo de aporte** (múltiple) |
| `filtertf_materia[]` | int[] | **Materia** (múltiple) |
| `page` | int | Número de página |
| `sortby` | string | Campo de ordenamiento |
| `order` | string | Dirección (`asc` o `desc`) |
| `datecut` | int | Filtrar por fecha (días) |

### Valores: `filtertf_tipo_aporte[]` (Tipo de Aporte)

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

### Valores: `filtertf_materia[]` (Materias - Ingeniería en Sistemas)

#### Materias de Ciclo Básico (1°-2° año)
| ID | Nombre |
|----|--------|
| -8 | General para todo Ingeniería en Sistemas |
| 5 | Álgebra y Geometría Analítica |
| 2 | Algoritmos y Estructuras de Datos |
| 7 | Análisis Matemático I |
| 14 | Análisis Matemático II |
| 3 | Arquitectura de Computadoras |
| 67 | Física I |
| 4 | Ingeniería y Sociedad |
| 72 | Legislación |
| 6 | Matemática Discreta |
| 63 | Química (Sistemas) |
| 69 | Química General |
| 66 | Sistemas de Representación |
| 1 | Sistemas y Organizaciones |

#### Materias de Ciclo Superior (3°-5° año)
| ID | Nombre |
|----|--------|
| 10 | Análisis de Sistemas |
| 13 | Física (Sistemas) |
| 68 | Física II |
| 64 | Inglés I |
| 65 | Inglés II |
| 8 | Paradigmas de Programación |
| 12 | Probabilidad y Estadística |
| 11 | Sintaxis y Semántica de los Lenguajes |
| 9 | Sistemas Operativos |

#### Materias Optativas/Especialización
| ID | Nombre |
|----|--------|
| 422 | Algoritmos Complejos para Estructuras de Datos Avanzadas |
| 16 | Comunicaciones |
| 76 | Comunicaciones y Redes |
| 20 | Creatividad |
| 15 | Diseño de Sistemas |
| 71 | Economía |
| 17 | Gestión de Datos |
| 21 | Gestión de Recursos Humanos |
| 22 | Innovación Tecnológica |
| 23 | Investigación Tecnológica |
| 70 | Matemática Superior |
| 75 | Metodología de la Conducción de Equipos de Trabajo |
| 19 | Modelos Numéricos |
| 18 | Redes de Información |
| 24 | Seguridad Informática |
| 25 | Técnicas Avanzadas de Programación |
| 26 | Técnicas de Gráficos por Computadoras |
| 358 | Investigación operativa (Sistemas) |
| 360 | Simulación (Sistemas) |
| 28 | Sistemas de Gestión I |
| 37 | Sistemas de Información Geográfica |
| 38 | Sistemas Distribuidos |
| 36 | Tecnologías Avanzadas en la Construcción de Software |
| 29 | Teoría de Control |
| 477 | Sistemas Embebidos Aplicados a Robótica |
| 73 | Ingeniería de Software |
| 32 | Introducción a la Ingeniería de Software |
| 40 | Gerenciamiento de Proyectos Informáticos |
| 44 | Inteligencia Artificial |
| 59 | Inteligencia de Negocios |
| 41 | Proyecto (Sistemas) |
| 50 | Sistemas Avanzados de Bases de Datos |
| 53 | Sistemas de Costos y Presupuestos |
| 74 | Sistemas de Gestión |
| 42 | Sistemas de Gestión II |
| 61 | Tecnologías Para la Explotación de Información |

### Parámetros de Ordenamiento

| Valor `sortby` | Descripción |
|----------------|-------------|
| `subject` | Por título del tema |
| `starter` | Por autor |
| `replies` | Por número de respuestas |
| `views` | Por vistas |
| `lastpost` | Por último mensaje |
| `tyl_tnumtyls` | Por agradecimientos |

### Parámetros de Búsqueda (search.php)

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `action` | string | Acción: `do_search`, `getnew`, `getdaily`, `unanswered` |
| `keywords` | string | Palabras clave |
| `author` | string | Buscar por usuario |
| `forums[]` | int[] | Foros a buscar (o `all`) |
| `postthread` | int | 1 = buscar todo, 2 = solo títulos |
| `findthreadst` | int | 1 = mínimo respuestas, 2 = máximo |
| `numreplies` | int | Número de respuestas |
| `postdate` | int | Antigüedad: 0=ualquiera, 1=ayer, 7=1 semana, 30=1 mes, etc. |
| `sortby` | string | Orden: `lastpost`, `starter`, `forum` |
| `sortordr` | string | `asc` o `desc` |
| `showresults` | string | `threads` o `posts` |
| `fids` | string | Forum IDs separados por coma |

---

## Ejemplos de Uso

### 1. Buscar parciales de una materia específica
```
GET /foro-sistemas?filtertf_tipo_aporte[]=1&filtertf_materia[]=2
```
Busca parciales (tipo 1) de Algoritmos (materia 2).

### 2. Búsqueda con texto
```
GET /foro-sistemas?search=automatizaciones&filtertf_tipo_aporte[]=1&filtertf_materia[]=1
```
La URL proporcionada: `?search=automatizaciones&filtertf_tipo_aporte%5B%5D=1&filtertf_materia%5B%5D=1`

### 3. Pagination
```
GET /foro-sistemas?page=2
GET /foro-sistemas?page=3
```

### 4. RSS de un foro específico
```
GET /syndication.php?fid=93
```

### 5. Temas sin respuesta
```
GET /search.php?action=unanswered
GET /search.php?action=unanswered&fids=93
```

### 6. Mensajes nuevos desde último visita
```
GET /search.php?action=getnew
GET /search.php?action=getnew&fids=93
```

---

## Notas Importantes

1. **URLs amigables**: El foro usa URLs amigables (rewrite), pero también acepta los parámetros tradicionales (`fid`, `tid`, etc.)

2. **Formato de array**: Los filtros múltiples usan la notación `[]`:
   - `filtertf_tipo_aporte[]=1&filtertf_tipo_aporte[]=2`
   - `filtertf_materia[]=1&filtertf_materia[]=8`

3. **Foro ID**: El foro de Sistemas tiene `fid=93`. Otros foros comunes:
   - Básicas: `fid=85`
   - Civil: `fid=86`
   - Eléctrica: `fid=87`
   - Electrónica: `fid=88`
   - Industrial: `fid=89`
   - Mecánica: `fid=90`
   - Química: `fid=92`
   - Aeronáutica: `fid=171`

4. **Autenticación**: No hay API pública. El scraping debe respetar los ToS del sitio.

5. **Rate Limiting**: El sitio puede bloquear requests excesivos.

---

## Construcción de URLs

### URL base + parámetros
```
https://www.utnianos.com.ar/foro/foro-sistemas
    ?search=palabra
    &filtertf_tipo_aporte[]=1
    &filtertf_tipo_aporte[]=2
    &filtertf_materia[]=1
    &filtertf_materia[]=8
    &page=1
    &sortby=lastpost
    &order=desc
```
