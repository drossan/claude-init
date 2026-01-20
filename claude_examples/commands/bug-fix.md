---
name: bug-fix
description: Identifica, reproduce y soluciona errores en Griddo API. Sigue el ciclo TDD (reproducción obligatoria) y rastrea el fallo a través de las capas hexagonales.
usage: "bug-fix [descripcion-error] [contexto-opcional]"
---

# Comando: Corrección de Errores (Bug Fix)

Este comando orquesta la resolución de fallos en Griddo API, desde el diagnóstico hasta la verificación, asegurando que se identifique la causa raíz en la capa adecuada.

## Flujo de Trabajo Orquestado

### 1. Investigación y Diagnóstico (Debugger Agent)
- Analizar reportes, logs y trazas de error.
- Rastrear el flujo: ¿Es un problema de validación (Infra), lógica (App) o reglas (Domain)?
- **Planificación**: Crear un plan de corrección en la raíz de `.claude/plans/` con `Aprobado: [ ]`, alineado con la `.claude/development_guide.md`.
- **Sesión**: Inicializar la sesión en `.claude/sessions/active/`.
- **DETENCIÓN OBLIGATORIA**: Informar al usuario y esperar aprobación.
- **Activación**: Una vez aprobado, mover el plan a `.claude/plans/active/`.
- **Agente**: `debugger`
- **Skills**: `debug-master`, `typescript`, `node-js`, `typeorm`

### 2. TDD - Fase Roja: Reproducción (Tester Agent)
- **REPRODUCCIÓN OBLIGATORIA**: Crear un test (unitario o de integración) que falle exactamente por el bug.
- **Agente**: `tester`
- **Skills**: `qa-engineer`, `vitest`, `supertest`

### 3. TDD - Fase Verde: Corrección (Developer Agent)
- Implementar la corrección mínima necesaria en la capa correspondiente.
- Asegurar la ausencia de regresiones ejecutando la suite del módulo afectado.
- **Agente**: `developer`
- **Skills**: Dependiendo del área (domain-expert, usecase-developer, infra-specialist)

### 4. Verificación y QA (Reviewer + Debugger)
- **QA**: Ejecutar la suite completa de tests.
- **Auto-Corrección**: Si hay regresiones, corregir automáticamente (máximo 3 ciclos).
- **Agentes**: `reviewer`, `debugger`
- **Skills**: `code-reviewer`, `debug-master`

### 5. Documentación y Versiones (Writer Agent)
- **Documentación**: Actualizar `CHANGELOG.md` (sección `Fixed`).
- **Versiones**: Categorizar el fix (normalmente Patch) y proponer incremento.
- **Agente**: `writer`
- **Skills**: `technical-writer`

### 6. Finalización (Git)
- Crear rama: `git checkout -b fix/{bug-name}`.
- Commits siguiendo la convención: `fix: description`.

## Reglas Críticas
- **Seguimiento en TIEMPO REAL**: Actualizar el archivo de sesión en `.claude/sessions/active/` tras cada acción significativa.
- **Prohibido Fix sin Test**: No se corrige nada sin un test que lo demuestre primero.
- **Invariabilidad de API**: Si el fix altera el contrato de la API, consultar explícitamente.
- **Aislamiento**: Usar mocks para no depender de la base de datos real en tests unitarios.

---

¿Qué error necesitamos diagnosticar y corregir hoy en Griddo API?
