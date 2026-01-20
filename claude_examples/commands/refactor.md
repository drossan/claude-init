---
name: refactor
description: Mejora la estructura interna de Griddo API sin cambiar su comportamiento externo. Asegura la integridad mediante la suite de tests existente y los principios SOLID.
usage: "refactor [modulo-o-archivo] [descripcion-mejora]"
---

# Comando: Refactorización (Refactor)

Este comando orquesta la mejora de la calidad, legibilidad o rendimiento del código en Griddo API, manteniendo intacta su funcionalidad externa y asegurando el cumplimiento de la arquitectura hexagonal.

## Flujo de Trabajo Orquestado

### 1. Preparación y Análisis (Planning Agent)
- Verificar que los tests actuales pasen (Estado Verde inicial).
- Analizar oportunidades de mejora (DRY, SOLID, Patrones de diseño) basándose en la `.claude/development_guide.md`.
- **Planificación**: Crear un plan de refactorización en la raíz de `.claude/plans/` con `Aprobado: [ ]`.
- **Sesión**: Inicializar la sesión en `.claude/sessions/active/`.
- **DETENCIÓN OBLIGATORIA**: Informar al usuario y esperar aprobación.
- **Activación**: Una vez aprobado, mover el plan a `.claude/plans/active/`.
- **Agente**: `planning-agent`
- **Skills**: `system-architect`, `fullstack-ts-expert`

### 2. Ejecución Incremental (Developer Agent)
- Aplicar cambios quirúrgicos manteniendo la compatibilidad de las interfaces.
- Mejorar el tipado y la estructura de las capas afectadas.
- **Agente**: `developer`
- **Skills**: `fullstack-ts-expert`, `domain-expert`, `usecase-developer`

### 3. Validación y QA (Tester Agent)
- Ejecutar tests unitarios e integración frecuentemente.
- Asegurar que la cobertura no disminuya.
- **Agente**: `tester`
- **Skills**: `qa-engineer`, `vitest`, `supertest`

### 4. Revisión y Auto-Corrección (Reviewer + Debugger)
- **Revisión**: Validar que la nueva estructura siga las Guidelines.
- **Auto-Corrección**: Corregir automáticamente si se detectan fallos (máximo 3 ciclos).
- **Agentes**: `reviewer`, `debugger`
- **Skills**: `code-reviewer`, `debug-master`

### 5. Documentación y Registro (Writer Agent)
- **Documentación**: Actualizar guías de diseño en `docs_dev/` si aplica.
- **Registro**: Actualizar `CHANGELOG.md` (sección `Changed`).
- **Agente**: `writer`
- **Skills**: `technical-writer`

### 6. Finalización (Git)
- Crear rama: `git checkout -b refactor/{nombre-del-refactor}`.
- Commits siguiendo la convención: `refactor: description`.

## Reglas Críticas
- **Seguimiento en TIEMPO REAL**: Actualizar el archivo de sesión en `.claude/sessions/active/` tras cada acción significativa.
- **Prohibido añadir Features**: No se introducen nuevas funcionalidades durante un refactor.
- **Invariabilidad de Comportamiento**: El cliente de la API no debe notar cambios.
- **API Pública Inalterada**: Mantener la compatibilidad hacia atrás rigurosamente.
- **Tests en Verde**: No concluir el refactor si hay algún test fallando.

---

¿Qué parte de la arquitectura hexagonal de Griddo API vamos a optimizar hoy?
