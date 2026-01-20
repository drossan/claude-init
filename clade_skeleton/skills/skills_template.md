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