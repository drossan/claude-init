---
name: plan-manage
description: Gestiona el ciclo de vida de un plan para Griddo API (iniciar, retomar, estado, finalizar). Organiza el contexto entre agentes y carpetas de estado.
usage: "plan-manage [start | resume | status | finish] [feature-name]"
---

# Comando: Gestión de Planes (Plan Manage)

Este comando facilita la gestión de los planes de desarrollo y su seguimiento a través de las capas de la arquitectura hexagonal, asegurando que el contexto se mantenga íntegro.

## Estructura de Organización

- **Raíz (`.claude/plans/`)**: Planes recién creados. Deben contener `Aprobado: [ ]`.
- **Activos (`.claude/plans/active/`)**: Planes aprobados y en implementación.
- **Sesiones (`.claude/sessions/active/`)**: Registro del progreso actual en tiempo real.
- **Completados (`.claude/plans/completed/` y `.claude/sessions/completed/`)**: Historial archivado.

## Modos de Uso

### 1. Iniciar (`start`)
- **Acción**: Crea el archivo de plan en la raíz de `.claude/plans/` con `Aprobado: [ ]`. Inicializa la sesión relacionada en `.claude/sessions/active/`.
- **Agente**: `planning-agent`

### 2. Aprobar y Activar (`approve` / `activate`)
- **Acción**: Tras la marca manual del usuario `Aprobado: [x]` en el plan de la raíz, este modo mueve el archivo a `.claude/plans/active/` para comenzar el desarrollo.
- **Agente**: `orchestrator-agent`

### 3. Retomar (`resume`)
- **Acción**: Lee el plan en `active/` y la sesión en `active/` para reconstruir el estado mental de los agentes.
- **Agente**: `planning-agent`

### 4. Estado (`status`)
- **Acción**: Analiza marcas de progreso `[x]` en el plan activo y los hitos en la sesión.
- **Agente**: `planning-agent`

### 5. Finalizar (`finish`)
- **Verificación OBLIGATORIA**:
    - [ ] `pre-flight` pasado con éxito.
    - [ ] `CHANGELOG.md` y JSDoc actualizados.
- **Acción**: 
    1. Mueve el plan a `completed/` añadiendo la fecha de finalización.
    2. Mueve la sesión a `completed/` añadiendo fecha de fin, resumen de cambios y hash del commit.
- **Agente**: `planning-agent`, `writer`, `orchestrator-agent`

## Reglas Críticas
- **Actualización en TIEMPO REAL**: Es obligatorio actualizar el archivo de sesión tras cada acción significativa para asegurar que el trabajo pueda reanudarse fácilmente.
- **No duplicidad**: No se permiten dos planes activos para la misma funcionalidad.

---

¿Qué plan o tarea de Griddo API necesitamos gestionar ahora?
