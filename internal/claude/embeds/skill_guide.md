# Guía Técnica para la Creación de Skills en Claude Code

## 1. Principio Fundamental: Skills como Conocimiento Inyectable

Una **Skill es un paquete modular y auto-contenido** que extiende las capacidades de Claude proporcionando:

- ✅ **Conocimiento especializado** (no razonamiento)
- ✅ **Lenguajes de programación** (los que requiera el proyecto)
- ✅ **Plugins, librerías o herramientas** (los que requiera el proyecto)
- ✅ **Workflows específicos de dominio**
- ✅ **Integraciones con herramientas**
- ✅ **Recursos reutilizables** (scripts, referencias, assets)

### Skills NO son:

- ❌ Agentes (que razonan)
- ❌ Commands (que orquestan)
- ❌ Tools (que ejecutan acciones en el sistema)

**Las Skills transforman a Claude de un agente generalista a un especialista equipado con conocimiento procedural que
ningún modelo puede poseer completamente.**

---

## 2. Anatomía de una Skill

### Estructura de Directorios

```
skill-name/
├── SKILL.md              (OBLIGATORIO)
│   ├── YAML frontmatter  (metadata)
│   └── Markdown body     (instrucciones)
│
└── Recursos Opcionales
    ├── scripts/          (código ejecutable)
    │   ├── rotate_pdf.py
    │   └── process_data.sh
    │
    ├── references/       (documentación para contexto)
    │   ├── api_docs.md
    │   ├── schemas.md
    │   └── policies.md
    │
    └── assets/           (archivos para output)
        ├── templates/
        ├── logo.png
        └── boilerplate/
```

---

## 3. Principio de Progressive Disclosure

Las skills usan un **sistema de carga de tres niveles** para gestionar el contexto eficientemente:

### Nivel 1: Metadata (siempre en contexto)

- `name` + `description` del frontmatter
- ~100 palabras
- **Determina cuándo se activa la skill**

### Nivel 2: SKILL.md body (cuando skill se activa)

- Instrucciones procedurales
- < 5k palabras recomendadas
- **Cómo usar la skill**

### Nivel 3: Recursos (cuando Claude los necesita)

- Scripts, referencias, assets
- Tamaño ilimitado*
- **Claude decide cuándo cargarlos**

*Ilimitado porque scripts pueden ejecutarse sin leer al contexto.

---

## 4. Tipos de Recursos

### 4.1 Scripts (`scripts/`)

**Propósito**: Código ejecutable para tareas que requieren fiabilidad determinista o se reescriben repetidamente.

**Cuándo incluirlos**:

- ✅ El mismo código se reescribe constantemente
- ✅ Se necesita fiabilidad determinista (procesamiento de archivos, cálculos)
- ✅ Operaciones complejas que no deben reinventarse cada vez

**Beneficios**:

- Token-efficient (no ocupan contexto al ejecutarse)
- Deterministas y testeables
- Reutilizables sin reescritura

**Nota**: Scripts pueden necesitar leerse para parches o ajustes específicos del entorno.

**Ejemplos**:

```
scripts/
├── rotate_pdf.py          # Rotación de PDFs
├── optimize_images.sh     # Optimización de imágenes
├── validate_schema.py     # Validación de esquemas
└── generate_report.js     # Generación de reportes
```

---

### 4.2 Referencias (`references/`)

**Propósito**: Documentación y material de referencia que Claude debe consultar mientras trabaja.

**Cuándo incluirlas**:

- ✅ Documentación que Claude debe referenciar durante el trabajo
- ✅ Información detallada que no cabe en SKILL.md
- ✅ Conocimiento que cambia con frecuencia

**Beneficios**:

- Mantiene SKILL.md conciso
- Se cargan solo cuando Claude las necesita
- Fácilmente actualizables

**Casos de uso**:

- Esquemas de base de datos
- Documentación de APIs
- Políticas de la empresa
- Guías de workflows detalladas
- Convenciones de código del proyecto

**Mejores prácticas**:

- Si archivos > 10k palabras, incluir patrones de búsqueda grep en SKILL.md
- **Evitar duplicación**: información vive en SKILL.md O en referencias, no en ambos
- Preferir referencias para info detallada; SKILL.md solo para procedimientos core

**Ejemplos**:

```
references/
├── api_docs.md            # Documentación de API
├── database_schema.md     # Esquemas de BD
├── policies.md            # Políticas de la empresa
├── conventions.md         # Convenciones de código
└── workflows.md           # Workflows detallados
```

---

### 4.3 Assets (`assets/`)

**Propósito**: Archivos que NO se cargan en contexto, sino que se usan en el output que Claude produce.

**Cuándo incluirlos**:

- ✅ La skill necesita archivos que estarán en el output final
- ✅ Templates, boilerplates, recursos visuales

**Beneficios**:

- Separa recursos de output de documentación
- Claude puede usar archivos sin cargarlos en contexto
- Acelera desarrollo al evitar recrear boilerplate

**Casos de uso**:

- Templates (HTML, React, documentos)
- Imágenes (logos, iconos)
- Boilerplate de código
- Fuentes tipográficas
- Documentos de muestra

**Ejemplos**:

```
assets/
├── templates/
│   ├── slides.pptx        # Template de presentaciones
│   ├── report.docx        # Template de reportes
│   └── email.html         # Template de emails
├── frontend-boilerplate/  # Proyecto React base
│   ├── package.json
│   ├── src/
│   └── public/
├── logo.png               # Logo de la empresa
└── fonts/
    └── brand-font.ttf
```

---

## 5. Proceso de Creación de Skills

### Paso 1: Entender la Skill con Ejemplos Concretos

**Objetivo**: Clarificar patrones de uso antes de construir.

**Saltar este paso solo si**: Los patrones de uso ya son muy claros.

**Preguntas a responder**:

1. ¿Qué funcionalidad debe soportar la skill?
2. ¿Puedes dar ejemplos de cómo se usaría?
3. ¿Qué diría un usuario que debería activar esta skill?

**Ejemplo - skill `image-editor`**:

```
Q: ¿Qué funcionalidad debería soportar?
A: Edición, rotación, redimensionado, optimización

Q: ¿Ejemplos de uso?
A: 
- "Elimina los ojos rojos de esta imagen"
- "Rota esta imagen 90 grados"
- "Redimensiona a 800x600"
- "Optimiza este PNG"

Q: ¿Qué activaría la skill?
A: Cualquier petición de manipulación de imágenes
```

**Concluir cuando**: Hay claridad sobre la funcionalidad que debe soportar.

---

### Paso 2: Planificar los Contenidos Reutilizables

**Objetivo**: Convertir ejemplos concretos en recursos de la skill.

**Proceso para cada ejemplo**:

1. Considerar cómo ejecutar el ejemplo desde cero
2. Identificar qué scripts, referencias o assets serían útiles al repetir estos workflows

**Ejemplo - skill `pdf-editor`**:

```
Query: "Ayúdame a rotar este PDF"

Análisis:
1. Rotar un PDF requiere reescribir el mismo código cada vez
2. Solución: Script `scripts/rotate_pdf.py`

Recursos a incluir:
- scripts/rotate_pdf.py
- scripts/merge_pdfs.py
- scripts/compress_pdf.py
```

**Ejemplo - skill `frontend-webapp-builder`**:

```
Queries: 
- "Construye una todo app"
- "Crea un dashboard para trackear pasos"

Análisis:
1. Escribir webapp frontend requiere mismo boilerplate HTML/React
2. Solución: Template con estructura de proyecto base

Recursos a incluir:
- assets/hello-world/ (boilerplate React)
- assets/templates/dashboard.html
- references/component-patterns.md
```

**Ejemplo - skill `bigquery`**:

```
Query: "¿Cuántos usuarios se logearon hoy?"

Análisis:
1. Consultar BigQuery requiere redescubrir schemas cada vez
2. Solución: Documentación de schemas

Recursos a incluir:
- references/schema.md (esquemas de tablas)
- references/query-patterns.md
- scripts/validate_query.py
```

**Resultado**: Lista de recursos reutilizables: scripts, referencias, assets.

---

### Paso 3: Inicializar la Skill

**Saltar este paso solo si**: La skill ya existe y solo se necesita iterar.

**Para skills nuevas**: SIEMPRE usar el script de inicialización.

```bash
# Crear nueva skill
scripts/init_skill.py <skill-name> --path <output-directory>

# Ejemplo
scripts/init_skill.py pdf-editor --path ./skills/
```

**El script genera**:

```
pdf-editor/
├── SKILL.md              # Template con frontmatter y TODOs
├── scripts/
│   └── example.py        # Script de ejemplo (borrar si no se usa)
├── references/
│   └── example.md        # Referencia de ejemplo (borrar si no se usa)
└── assets/
    └── example.txt       # Asset de ejemplo (borrar si no se usa)
```

**Después de inicialización**:

- Personalizar o eliminar archivos de ejemplo
- Completar TODOs en SKILL.md

---

### Paso 4: Editar la Skill

**Recordar**: La skill es para **otra instancia de Claude**. Incluir información que sería útil y no-obvia para Claude.

#### 4.1 Implementar Recursos Reutilizables

**Orden recomendado**:

1. Crear scripts en `scripts/`
2. Documentar en `references/`
3. Agregar assets en `assets/`
4. Eliminar archivos de ejemplo no necesarios

**Nota**: Este paso puede requerir input del usuario (assets, documentación interna, etc.)

#### 4.2 Actualizar SKILL.md

**Estilo de escritura**:

- ✅ **Forma imperativa/infinitiva** (verb-first)
- ✅ Lenguaje objetivo e instructivo
- ❌ NO usar segunda persona ("tú debes")
- ❌ NO usar primera persona ("yo haré")

**Ejemplo correcto**:

```markdown
Para rotar un PDF:

1. Ejecutar `scripts/rotate_pdf.py` con el archivo de entrada
2. Especificar el ángulo de rotación (90, 180, 270)
3. Validar que el output se generó correctamente
```

**Ejemplo incorrecto**:

```markdown
Debes rotar el PDF usando el script.
Tú especificarás el ángulo.
```

**Preguntas a responder en SKILL.md**:

1. ¿Cuál es el propósito de la skill? (2-3 frases)
2. ¿Cuándo debe usarse la skill?
3. ¿Cómo debe Claude usar la skill en la práctica?
4. ¿Cómo se usan los recursos incluidos (scripts, references, assets)?

---

### Paso 5: Empaquetar la Skill

**Objetivo**: Crear un zip distribuible validado automáticamente.

```bash
# Empaquetar skill
scripts/package_skill.py <path/to/skill-folder>

# Con directorio de output personalizado
scripts/package_skill.py <path/to/skill-folder> ./dist
```

**El script realiza**:

#### 5.1 Validación Automática

Verifica:

- ✅ Formato de frontmatter YAML
- ✅ Campos requeridos (name, description)
- ✅ Convenciones de nomenclatura
- ✅ Estructura de directorios
- ✅ Completitud de description
- ✅ Organización de archivos
- ✅ Referencias a recursos en SKILL.md

#### 5.2 Empaquetado (si validación pasa)

Crea:

- `{skill-name}.zip` con todos los archivos
- Mantiene estructura de directorios
- Listo para distribución

**Si validación falla**:

- Script reporta errores
- No crea package
- Corregir errores y re-ejecutar

---

### Paso 6: Iterar

**Workflow de iteración**:

1. Usar la skill en tareas reales
2. Notar dificultades o ineficiencias
3. Identificar cómo mejorar SKILL.md o recursos
4. Implementar cambios
5. Testear de nuevo

**Común después de**: Primera vez usando la skill, con contexto fresco de su desempeño.

---

## 6. Template Oficial de Skills

```markdown
---
name: {skill-name}
description: Esta skill debe usarse cuando el usuario necesite {acción específica}. Se activa con peticiones como: "{trigger-1}", "{trigger-2}".
license: Complete terms in LICENSE.txt
version: 1.0.0
author: {team/person}
category: {development | data | design | automation | business}
tags: [tag1, tag2, tag3]
---

# {Skill Name}

{Breve descripción del propósito de la skill en 2-3 frases.
Explicar qué problema resuelve y qué capacidades añade a Claude.}

## Cuándo Usar Esta Skill

Esta skill debe usarse cuando:

- {Escenario 1 específico}
- {Escenario 2 específico}
- {Escenario 3 específico}

Triggers comunes:

- "{ejemplo de query del usuario 1}"
- "{ejemplo de query del usuario 2}"
- "{ejemplo de query del usuario 3}"

## Workflow Principal

{Instrucciones paso a paso en modo IMPERATIVO.
Enfocarse en el proceso lógico que debe seguir Claude.}

### 1. Análisis Inicial

Antes de proceder:

1. Identificar {variable/requisito específico}
2. Verificar que {condición necesaria} está presente
3. Determinar si {decisión clave}

### 2. Ejecución

Para completar la tarea:

1. Ejecutar `scripts/{script-name}` para {propósito}
2. Consultar `references/{doc-name}` si se necesita {tipo de info}
3. Aplicar template de `assets/{template-name}` como base
4. {Paso adicional específico del workflow}

### 3. Validación

Verificar que:

- [ ] {Criterio de validación 1}
- [ ] {Criterio de validación 2}
- [ ] {Criterio de validación 3}

### 4. Output

Presentar resultados:

- Formato: {especificar formato esperado}
- Incluir: {qué debe contener el output}
- Omitir: {qué NO debe incluirse}

## Recursos de la Skill

### Scripts (`scripts/`)

Scripts ejecutables disponibles:

#### `scripts/{script-name}.{ext}`

**Propósito**: {Qué hace el script}

**Uso**:

```bash
{lenguaje} scripts/{script-name}.{ext} [argumentos]
```

**Parámetros**:

- `arg1`: {Descripción del argumento}
- `arg2`: {Descripción del argumento}

**Output**: {Qué retorna/genera}

**Ejemplo**:

```bash
python scripts/rotate_pdf.py input.pdf --angle 90 --output rotated.pdf
```

---

### Referencias (`references/`)

Documentación técnica que cargar según necesidad:

#### `references/{doc-name}.md`

**Contenido**: {Tipo de información que contiene}

**Cuándo consultar**: {Bajo qué circunstancias debe leerse}

**Estructura**: {Breve descripción de cómo está organizado}

**Búsqueda rápida**: Para archivos grandes, usar patrones grep:

```bash
grep -i "{patrón relevante}" references/{doc-name}.md
```

---

### Assets (`assets/`)

Archivos para usar en output final:

#### `assets/{asset-name}`

**Tipo**: {Template / Image / Boilerplate / etc.}

**Uso**: {Cómo debe usarse este asset}

**Modificaciones**: {Qué partes deben personalizarse}

## Ejemplos de Uso

### Ejemplo 1: {Caso de uso común}

**Input del usuario**:
> "{query exacto del usuario}"

**Proceso**:

1. {Paso que Claude debe seguir}
2. {Paso que Claude debe seguir}
3. {Paso que Claude debe seguir}

**Output esperado**:

```
{Ejemplo de output que Claude debería generar}
```

---

### Ejemplo 2: {Otro caso de uso}

**Input del usuario**:
> "{query exacto del usuario}"

**Proceso**:

1. {Paso específico}
2. {Paso específico}

**Output esperado**:

```
{Ejemplo de output}
```

## Presentación de Resultados

Al completar la tarea:

1. **Resumir cambios**: "Completada {acción} aplicando {método}. Resultados: {resumen}"
2. **Formato de output**: {Especificar formato exacto}
3. **Incluir métricas**: {Si aplica, qué métricas mostrar}
4. **Adjuntar archivos**: {Si se generaron archivos, cómo presentarlos}

**Ejemplo de resumen**:

```
Rotado PDF correctamente 90° en sentido horario.
- Archivo de entrada: document.pdf (1.2 MB)
- Archivo de salida: document-rotated.pdf (1.2 MB)
- Tiempo de procesamiento: 0.3s
```

## Troubleshooting

### Problema: {Error común 1}

**Síntoma**: {Cómo se manifiesta}

**Causa**: {Por qué ocurre}

**Solución**:

```bash
{Comando o acción para resolverlo}
```

---

### Problema: {Error común 2}

**Síntoma**: {Cómo se manifiesta}

**Causa**: {Por qué ocurre}

**Solución**:

1. {Paso de resolución}
2. {Paso de resolución}

---

### Problema: Referencias muy grandes

**Síntoma**: Archivo de referencia > 10k palabras

**Solución**: Usar grep para buscar secciones específicas:

```bash
grep -A 10 "{término clave}" references/{archivo}.md
```

---

### Problema: Permisos de ejecución

**Síntoma**: Error "Permission denied" al ejecutar script

**Solución**:

```bash
chmod +x scripts/*.sh
chmod +x scripts/*.py
```

## Consideraciones Especiales

### Rendimiento

- {Nota sobre rendimiento si aplica}
- {Limitaciones conocidas}

### Seguridad

- {Consideraciones de seguridad}
- {Datos sensibles que no deben exponerse}

### Compatibilidad

- {Requisitos de versiones}
- {Dependencias externas}

## Mejoras Futuras (Roadmap)

- [ ] {Feature o mejora planificada 1}
- [ ] {Feature o mejora planificada 2}
- [ ] {Feature o mejora planificada 3}

```

---

## 7. Calidad del Frontmatter (Metadata)

**CRÍTICO**: El `name` y `description` determinan **cuándo Claude usará la skill**.

### 7.1 Name (nombre)

**Formato**: `{dominio}-{acción}` (kebab-case)

✅ **Buenos nombres**:
- `pdf-editor`
- `bigquery-analyst`
- `frontend-builder`
- `brand-guidelines`

❌ **Malos nombres**:
- `PDFEditor` (PascalCase)
- `pdf_editor` (snake_case)
- `editor` (demasiado genérico)
- `do-everything` (no específico)

---

### 7.2 Description (descripción)

**Estilo**: Tercera persona, específico sobre uso

✅ **Buena description**:
```yaml
description: Esta skill debe usarse cuando el usuario necesite editar, rotar, o comprimir archivos PDF. Se activa con peticiones como "rota este PDF", "combina estos PDFs", o "reduce el tamaño de este documento".
```

❌ **Mala description**:

```yaml
description: Use esta skill para PDFs.  # Demasiado corta, no específica
```

❌ **Mala description**:

```yaml
description: Esta skill permite a los usuarios trabajar con archivos PDF realizando diversas operaciones incluyendo pero no limitadas a...  # Demasiado verbosa, no clara
```

**Elementos clave en description**:

1. **Cuándo usarse** ("cuando el usuario necesite...")
2. **Triggers específicos** ("peticiones como...")
3. **Alcance claro** (qué incluye y qué no)

---

## 8. Mejores Prácticas

### 8.1 Evitar Duplicación

**Regla**: Información vive en SKILL.md **O** en referencias, **nunca en ambos**.

**Preferir referencias para**:

- Documentación detallada (> 500 palabras)
- Info que cambia frecuentemente
- Esquemas complejos
- Especificaciones técnicas exhaustivas

**Mantener en SKILL.md solo**:

- Instrucciones procedurales core
- Guía de workflow
- Referencias a dónde encontrar info detallada

---

### 8.2 Granularidad de Scripts

**Crear script cuando**:

- ✅ Mismo código se reescribe 3+ veces
- ✅ Necesita ser determinista (procesamiento de archivos)
- ✅ Lógica compleja que no debe reinventarse

**NO crear script para**:

- ❌ Tareas simples que Claude hace bien (ej: formatear JSON)
- ❌ Lógica que varía mucho caso por caso
- ❌ Una sola vez / uso único

---

### 8.3 Organización de Referencias

**Para archivos grandes** (> 10k palabras):

- Incluir tabla de contenidos en SKILL.md
- Proporcionar patrones de búsqueda grep
- Dividir en múltiples archivos por tema

**Ejemplo en SKILL.md**:

```markdown
### `references/api_documentation.md` (grande - 25k palabras)

**Estructura**:

- Sección 1: Authentication (líneas 1-500)
- Sección 2: User Endpoints (líneas 501-1200)
- Sección 3: Data Endpoints (líneas 1201-2000)

**Búsqueda rápida**:

```bash
# Buscar endpoint específico
grep -i "POST /api/users" references/api_documentation.md

# Buscar info de autenticación
grep -A 20 "## Authentication" references/api_documentation.md
```

```

---

### 8.4 Uso de Assets

**Assets deben**:
- ✅ Ser archivos finales o templates listos para usar
- ✅ Estar organizados por tipo (templates/, images/, etc.)
- ✅ Incluir instrucciones de personalización en SKILL.md

**Assets NO deben**:
- ❌ Ser documentación (eso va en references/)
- ❌ Ser código ejecutable (eso va en scripts/)
- ❌ Cargarse en contexto (solo usarse en output)

---

## 9. Anti-patrones en Skills

### ❌ God Skill
**Problema**: Skill que intenta hacer demasiado

**Ejemplo**:
```yaml
name: developer-assistant
description: Esta skill hace desarrollo completo de software, testing, deployment, y todo lo relacionado a programación.
```

**Por qué es malo**: Demasiado genérica, no se activa apropiadamente

**Solución**: Dividir en skills específicas:

- `backend-api-developer`
- `frontend-component-builder`
- `test-automation`
- `deployment-manager`

---

### ❌ Skill Sin Workflow Claro

**Problema**: Solo lista recursos sin explicar cómo usarlos

**Ejemplo**:

```markdown
## Recursos

- scripts/script1.py
- scripts/script2.sh
- references/doc.md
```

**Por qué es malo**: Claude no sabe cuándo/cómo usar cada recurso

**Solución**: Workflow explícito:

```markdown
## Workflow

### 1. Análisis

Ejecutar `scripts/analyze.py` para evaluar el input

### 2. Procesamiento

Según el tipo identificado:

- Si es PDF → usar `scripts/process_pdf.py`
- Si es imagen → usar `scripts/process_image.py`

### 3. Validación

Consultar `references/validation_rules.md` para criterios
```

---

### ❌ Descripciones Vagas

**Problema**: Triggers no específicos

**Ejemplo**:

```yaml
description: Para trabajar con archivos.
```

**Por qué es malo**: No queda claro cuándo activarse

**Solución**:

```yaml
description: Esta skill debe usarse cuando el usuario necesite convertir, comprimir, o validar archivos PDF. Se activa con peticiones como "convierte este Word a PDF", "reduce el tamaño de este PDF", o "verifica que este PDF es válido".
```

---

### ❌ Duplicación de Contenido

**Problema**: Misma info en SKILL.md y en references/

**Ejemplo**:

```markdown
# SKILL.md

## API Endpoints

POST /api/users - Crear usuario
GET /api/users/:id - Obtener usuario
...

# references/api_docs.md

## API Endpoints

POST /api/users - Crear usuario
GET /api/users/:id - Obtener usuario
...
```

**Por qué es malo**: Desperdicia tokens, info desincronizada

**Solución**:

```markdown
# SKILL.md

## Workflow

Para consultar endpoints de API, referirse a `references/api_docs.md` sección "Endpoints".

# references/api_docs.md

## Endpoints

[Documentación completa aquí]
```

---

### ❌ Scripts Sin Documentación

**Problema**: Scripts sin explicar parámetros o uso

**Ejemplo**:

```markdown
### Scripts

- `scripts/process.py`
```

**Por qué es malo**: Claude no sabe cómo invocar el script

**Solución**:

```markdown
### `scripts/process.py`

**Propósito**: Procesar archivos CSV y generar reporte JSON

**Uso**:

```bash
python scripts/process.py <input.csv> [--output report.json] [--verbose]
```

**Parámetros**:

- `input.csv`: Archivo CSV de entrada (obligatorio)
- `--output`: Nombre del archivo JSON de salida (opcional, default: output.json)
- `--verbose`: Modo verbose para debugging (opcional)

**Ejemplo**:

```bash
python scripts/process.py data/sales.csv --output reports/sales_summary.json
```

```

---

## 10. Categorías de Skills

### 10.1 Development Skills
**Propósito**: Desarrollo de software, código, APIs

**Ejemplos**:
- `backend-api-builder`
- `frontend-component-generator`
- `database-migration-manager`
- `test-automation-creator`

---

### 10.2 Data Skills
**Propósito**: Procesamiento, análisis, visualización de datos

**Ejemplos**:
- `bigquery-analyst`
- `data-visualization-builder`
- `csv-processor`
- `sql-query-optimizer`

---

### 10.3 Design Skills
**Propósito**: Diseño, assets visuales, branding

**Ejemplos**:
- `brand-guidelines-enforcer`
- `presentation-builder`
- `image-optimizer`
- `icon-generator`

---

### 10.4 Automation Skills
**Propósito**: Automatización de tareas repetitivas

**Ejemplos**:
- `email-template-generator`
- `report-automator`
- `deployment-orchestrator`
- `backup-manager`

---

### 10.5 Business Skills
**Propósito**: Procesos de negocio, documentación corporativa

**Ejemplos**:
- `contract-generator`
- `invoice-creator`
- `meeting-notes-formatter`
- `policy-enforcer`

---

## 11. Validación de Skills

### Checklist Pre-Empaquetado

Antes de empaquetar, verificar:

#### Estructura (5/5)
- [ ] SKILL.md existe y tiene frontmatter válido
- [ ] name y description completos y específicos
- [ ] Directorios (scripts/, references/, assets/) presentes si se usan
- [ ] No hay archivos de ejemplo sin personalizar
- [ ] Estructura sigue convenciones de nomenclatura

#### Contenido (7/7)
- [ ] Description usa tercera persona y especifica triggers
- [ ] Workflow está en modo imperativo
- [ ] Todos los recursos (scripts/references/assets) están documentados en SKILL.md
- [ ] Scripts tienen documentación de uso y parámetros
- [ ] Referencias tienen descripción de cuándo consultarlas
- [ ] Assets tienen instrucciones de personalización
- [ ] Ejemplos de uso incluidos

#### Calidad (4/4)
- [ ] No hay duplicación entre SKILL.md y references/
- [ ] Workflow es claro y paso a paso
- [ ] Troubleshooting cubre errores comunes
- [ ] Description activará la skill apropiadamente

**Mínimo para validación**: 16/16 ✅

---

## 12. Ejemplo Real Completo

### Ejemplo Completo: `pdf-editor` Skill

```markdown
---
name: pdf-editor
description: Esta skill debe usarse cuando el usuario necesite manipular archivos PDF (rotar, combinar, comprimir, dividir, extraer páginas). Se activa con peticiones como "rota este PDF 90 grados", "combina estos dos PDFs", "comprime este documento", o "extrae las páginas 1-5".
license: Complete terms in LICENSE.txt
version: 1.0.0
author: platform-team
category: automation
tags: [pdf, documents, file-processing]
---

# PDF Editor

Skill para manipulación avanzada de archivos PDF. Permite rotar, combinar, comprimir, dividir y extraer páginas de documentos PDF usando scripts optimizados y deterministas.

## Cuándo Usar Esta Skill

Esta skill debe usarse cuando:
- El usuario necesite rotar páginas de un PDF
- Se deban combinar múltiples PDFs en uno solo
- Un PDF necesite comprimirse para reducir tamaño
- Se requiera dividir un PDF en múltiples archivos
- Necesite extraerse un rango específico de páginas

Triggers comunes:
- "Rota este PDF 90 grados en sentido horario"
- "Combina estos tres PDFs en uno solo"
- "Reduce el tamaño de este documento PDF"
- "Divide este PDF en archivos separados por página"
- "Extrae las páginas 5 a 10 de este PDF"

## Workflow Principal

### 1. Análisis Inicial

Antes de proceder:
1. Identificar el tipo de operación solicitada (rotar/combinar/comprimir/dividir/extraer)
2. Verificar que el archivo PDF de entrada está accesible
3. Determinar parámetros específicos:
   - Para rotación: ángulo (90, 180, 270)
   - Para combinación: orden de archivos
   - Para compresión: nivel de calidad deseado
   - Para división/extracción: rango de páginas

### 2. Ejecución

Según el tipo de operación:

#### Rotación
```bash
python scripts/rotate_pdf.py <input.pdf> --angle <90|180|270> --output <output.pdf>
```

#### Combinación

```bash
python scripts/merge_pdfs.py <pdf1> <pdf2> [pdf3 ...] --output <merged.pdf>
```

#### Compresión

```bash
python scripts/compress_pdf.py <input.pdf> --quality <low|medium|high> --output <compressed.pdf>
```

#### División

```bash
python scripts/split_pdf.py <input.pdf> --output-dir <output_directory>
```

#### Extracción

```bash
python scripts/extract_pages.py <input.pdf> --pages <start-end> --output <extracted.pdf>
```

### 3. Validación

Verificar que:

- [ ] El archivo de salida se generó correctamente
- [ ] El tamaño del archivo es razonable (no aumentó inesperadamente)
- [ ] El PDF resultante es válido (puede abrirse sin errores)
- [ ] La operación logró el objetivo (páginas rotadas, archivos combinados, etc.)

### 4. Output

Presentar resultados:

- Formato: Resumen textual + estadísticas de la operación
- Incluir: Nombre del archivo de salida, tamaño, número de páginas
- Omitir: Detalles técnicos internos del procesamiento

## Recursos de la Skill

### Scripts (`scripts/`)

#### `scripts/rotate_pdf.py`

**Propósito**: Rotar todas las páginas de un PDF en el ángulo especificado

**Uso**:

```bash
python scripts/rotate_pdf.py <input.pdf> --angle <angle> --output <output.pdf>
```

**Parámetros**:

- `input.pdf`: Archivo PDF de entrada (obligatorio)
- `--angle`: Ángulo de rotación - 90, 180, o 270 grados (obligatorio)
- `--output`: Nombre del archivo de salida (opcional, default: input_rotated.pdf)

**Output**: PDF rotado en la ubicación especificada

**Ejemplo**:

```bash
python scripts/rotate_pdf.py document.pdf --angle 90 --output document_rotated.pdf
```

---

#### `scripts/merge_pdfs.py`

**Propósito**: Combinar múltiples archivos PDF en uno solo

**Uso**:

```bash
python scripts/merge_pdfs.py <pdf1> <pdf2> [pdf3 ...] --output <merged.pdf>
```

**Parámetros**:

- `pdf1, pdf2, ...`: Archivos PDF a combinar en orden (mínimo 2)
- `--output`: Nombre del archivo combinado (opcional, default: merged.pdf)

**Output**: PDF único con todos los documentos combinados

**Ejemplo**:

```bash
python scripts/merge_pdfs.py intro.pdf content.pdf appendix.pdf --output complete_document.pdf
```

---

#### `scripts/compress_pdf.py`

**Propósito**: Reducir el tamaño de un PDF optimizando imágenes y eliminando metadata innecesaria

**Uso**:

```bash
python scripts/compress_pdf.py <input.pdf> --quality <level> --output <compressed.pdf>
```

**Parámetros**:

- `input.pdf`: Archivo PDF a comprimir (obligatorio)
- `--quality`: Nivel de compresión - low/medium/high (opcional, default: medium)
    - `low`: Máxima compresión, menor calidad
    - `medium`: Balance compresión/calidad
    - `high`: Mínima compresión, máxima calidad
- `--output`: Nombre del archivo comprimido (opcional, default: input_compressed.pdf)

**Output**: PDF comprimido

**Ejemplo**:

```bash
python scripts/compress_pdf.py large_document.pdf --quality medium --output optimized.pdf
```

---

#### `scripts/split_pdf.py`

**Propósito**: Dividir un PDF en múltiples archivos, uno por página

**Uso**:

```bash
python scripts/split_pdf.py <input.pdf> --output-dir <directory>
```

**Parámetros**:

- `input.pdf`: Archivo PDF a dividir (obligatorio)
- `--output-dir`: Directorio donde guardar las páginas (opcional, default: output/)

**Output**: Múltiples archivos PDF (page_1.pdf, page_2.pdf, etc.)

**Ejemplo**:

```bash
python scripts/split_pdf.py document.pdf --output-dir ./pages/
```

---

#### `scripts/extract_pages.py`

**Propósito**: Extraer un rango específico de páginas de un PDF

**Uso**:

```bash
python scripts/extract_pages.py <input.pdf> --pages <start-end> --output <extracted.pdf>
```

**Parámetros**:

- `input.pdf`: Archivo PDF fuente (obligatorio)
- `--pages`: Rango de páginas a extraer, formato: start-end (obligatorio)
- `--output`: Nombre del archivo con páginas extraídas (opcional, default: extracted.pdf)

**Output**: PDF con solo las páginas especificadas

**Ejemplo**:

```bash
python scripts/extract_pages.py report.pdf --pages 5-10 --output summary.pdf
```

---

### Referencias (`references/`)

#### `references/pdf_standards.md`

**Contenido**: Especificaciones técnicas del formato PDF, versiones soportadas, y limitaciones conocidas

**Cuándo consultar**: Al encontrar PDFs con características especiales (encriptación, formularios, anotaciones)

**Estructura**:

- Versiones de PDF (1.4 - 2.0)
- Características soportadas/no soportadas
- Manejo de PDFs encriptados
- Limitaciones de compresión

---

#### `references/troubleshooting_guide.md`

**Contenido**: Guía detallada de resolución de problemas comunes con PDFs

**Cuándo consultar**: Cuando un script falla o produce resultados inesperados

**Búsqueda rápida**:

```bash
# Buscar error específico
grep -i "encryption error" references/troubleshooting_guide.md

# Buscar por tipo de problema
grep -A 10 "## Compression Issues" references/troubleshooting_guide.md
```

---

### Assets (`assets/`)

No se incluyen assets en esta skill ya que todas las operaciones trabajan directamente con PDFs del usuario sin
necesidad de templates.

## Ejemplos de Uso

### Ejemplo 1: Rotación de PDF escaneado incorrectamente

**Input del usuario**:
> "Este PDF está girado 90 grados, rótaloadecuadamente"

**Proceso**:

1. Identificar que necesita rotación
2. Determinar ángulo correcto (probablemente 270° para corregir rotación de 90°)
3. Ejecutar: `python scripts/rotate_pdf.py scanned.pdf --angle 270 --output scanned_corrected.pdf`
4. Validar que el PDF resultante está correctamente orientado

**Output esperado**:

```
PDF rotado correctamente 270° en sentido horario.
- Archivo de entrada: scanned.pdf (2.4 MB, 15 páginas)
- Archivo de salida: scanned_corrected.pdf (2.4 MB, 15 páginas)
- Tiempo de procesamiento: 0.8s
- Estado: ✓ Completado exitosamente
```

---

### Ejemplo 2: Combinación de documentos para envío

**Input del usuario**:
> "Combina mi CV, carta de presentación y referencias en un solo PDF"

**Proceso**:

1. Identificar los tres archivos: cv.pdf, cover_letter.pdf, references.pdf
2. Determinar orden apropiado (CV primero, carta segundo, referencias último)
3. Ejecutar: `python scripts/merge_pdfs.py cv.pdf cover_letter.pdf references.pdf --output application_complete.pdf`
4. Validar que todas las páginas se combinaron correctamente

**Output esperado**:

```
PDFs combinados exitosamente en orden secuencial.
- Archivos combinados: 3 (cv.pdf, cover_letter.pdf, references.pdf)
- Total de páginas: 12 (2 + 1 + 9)
- Archivo de salida: application_complete.pdf (1.8 MB)
- Tiempo de procesamiento: 0.5s
- Estado: ✓ Completado exitosamente
```

---

### Ejemplo 3: Compresión de PDF grande para email

**Input del usuario**:
> "Este PDF es muy pesado para enviar por email, comprímelo"

**Proceso**:

1. Verificar tamaño del archivo (ej: 15 MB)
2. Determinar nivel de compresión (medium por defecto, o preguntar al usuario)
3. Ejecutar: `python scripts/compress_pdf.py large_report.pdf --quality medium --output report_compressed.pdf`
4. Validar reducción de tamaño y calidad aceptable

**Output esperado**:

```
PDF comprimido exitosamente con calidad media.
- Archivo original: large_report.pdf (15.2 MB, 45 páginas)
- Archivo comprimido: report_compressed.pdf (3.8 MB, 45 páginas)
- Reducción: 75% (11.4 MB ahorrados)
- Calidad: Media (adecuada para email)
- Tiempo de procesamiento: 2.1s
- Estado: ✓ Completado exitosamente
```

## Presentación de Resultados

Al completar cualquier operación de PDF:

1. **Resumir acción**: "Completada {operación} aplicando {método}. Resultados: {resumen}"
2. **Formato de output**: Texto estructurado con estadísticas clave
3. **Incluir métricas**:
    - Tamaño de archivos (entrada/salida)
    - Número de páginas
    - Tiempo de procesamiento
    - Reducción de tamaño (para compresión)
4. **Adjuntar archivos**: Indicar ubicación del archivo generado

**Ejemplo de resumen completo**:

```
Operación: Rotación de PDF
- Entrada: document.pdf (1.2 MB, 8 páginas)
- Salida: document_rotated.pdf (1.2 MB, 8 páginas)
- Ángulo aplicado: 90° (sentido horario)
- Tiempo: 0.3s
- Estado: ✓ Éxito
```

## Troubleshooting

### Problema: Error de permisos al ejecutar scripts

**Síntoma**: `Permission denied` al ejecutar cualquier script Python

**Causa**: Scripts no tienen permisos de ejecución

**Solución**:

```bash
chmod +x scripts/*.py
```

---

### Problema: PDF encriptado o protegido

**Síntoma**: Error "PDF is encrypted" o "Password required"

**Causa**: El PDF tiene protección con contraseña

**Solución**:

1. Solicitar contraseña al usuario
2. Consultar `references/pdf_standards.md` sección "Encrypted PDFs"
3. Si no hay contraseña disponible, informar al usuario que no se puede procesar

---

### Problema: Compresión no reduce tamaño significativamente

**Síntoma**: PDF comprimido tiene casi el mismo tamaño que el original

**Causa**: El PDF original ya está optimizado o contiene mayormente texto (no imágenes)

**Solución**:

1. Informar al usuario que el PDF ya está optimizado
2. Explicar que PDFs con texto plano no comprimen mucho
3. Sugerir alternativas si el tamaño sigue siendo problema (dividir, extraer páginas)

---

### Problema: Script falla con PDFs muy grandes

**Síntoma**: Error de memoria o timeout

**Causa**: PDF demasiado grande (> 100 MB o > 500 páginas)

**Solución**:

1. Dividir el PDF en chunks más pequeños
2. Procesar cada chunk individualmente
3. Combinar resultados al final
4. Consultar `references/troubleshooting_guide.md` sección "Large Files"

---

### Problema: Calidad visual degradada después de compresión

**Síntoma**: Texto borroso o imágenes pixeladas en PDF comprimido

**Causa**: Nivel de compresión demasiado agresivo

**Solución**:

1. Re-comprimir con nivel `high` en vez de `medium` o `low`
2. Explicar trade-off entre tamaño y calidad al usuario

```bash
python scripts/compress_pdf.py input.pdf --quality high --output better_quality.pdf
```

## Consideraciones Especiales

### Rendimiento

- PDFs < 5 MB: procesamiento instantáneo (< 1s)
- PDFs 5-20 MB: procesamiento rápido (1-3s)
- PDFs > 20 MB: puede tomar 5-10s según operación
- Combinación de muchos archivos: tiempo proporcional al número de archivos

### Seguridad

- Los scripts NO almacenan ni transmiten contenido de PDFs
- Archivos temporales se eliminan automáticamente después del procesamiento
- No se accede a metadatos sensibles sin autorización explícita
- PDFs encriptados requieren contraseña proporcionada por el usuario

### Compatibilidad

- Requisitos: Python 3.8+
- Dependencias: PyPDF2, pikepdf (instaladas automáticamente)
- Formatos soportados: PDF 1.4 - 2.0
- Limitaciones conocidas: PDFs con formularios XFA no soportados (consultar `references/pdf_standards.md`)

## Mejoras Futuras (Roadmap)

- [ ] Soporte para extracción de texto de PDFs escaneados (OCR)
- [ ] Marca de agua (watermarking) en batch
- [ ] Conversión PDF a imágenes (PNG/JPG)
- [ ] Soporte para firmas digitales
- [ ] Optimización de PDFs para web (linearization)

```

---

## 13. Resumen Ejecutivo

### Principios Fundamentales

📦 **Skill = Conocimiento Modular** - Paquetes auto-contenidos de expertise  
🎯 **Metadata = Activación** - name + description determinan cuándo se usa  
📊 **Progressive Disclosure** - 3 niveles de carga (metadata → SKILL.md → recursos)  
🛠️ **Scripts = Determinismo** - Para código que se reescribe constantemente  
📚 **Referencias = Documentación** - Cargada solo cuando se necesita  
🎨 **Assets = Output** - Archivos para usar, no para leer  
🚫 **No Duplicación** - Info vive en SKILL.md O referencias, nunca ambos  
✍️ **Modo Imperativo** - Instrucciones verb-first, objetivas  

---

## 14. Comandos de Utilidad

### Inicializar Skill
```bash
scripts/init_skill.py <skill-name> --path <output-directory>
```

### Empaquetar Skill

```bash
scripts/package_skill.py <path/to/skill-folder> [output-dir]
```

### Validar Skill (pre-empaquetado)

```bash
scripts/validate_skill.py <path/to/skill-folder>
```

---

## 15. Recursos Adicionales

### Plantillas

- `templates/basic-skill.md` - Skill básica sin recursos
- `templates/script-heavy-skill.md` - Skill con múltiples scripts
- `templates/documentation-skill.md` - Skill orientada a referencias

### Ejemplos

- `examples/pdf-editor/` - Manipulación de archivos
- `examples/bigquery-analyst/` - Análisis de datos
- `examples/frontend-builder/` - Generación de código

### Documentación

- `docs/skill-best-practices.md` - Mejores prácticas detalladas
- `docs/metadata-guide.md` - Guía de frontmatter efectivo
- `docs/progressive-disclosure.md` - Cómo optimizar carga de contexto

---

**Versión de la guía**: 2.0.0  
**Última actualización**: 2025-01-20  
**Basado en**: Documentación oficial de Anthropic Skills