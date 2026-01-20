---
name: new-feature
description: Planifica e implementa una nueva funcionalidad en Griddo API siguiendo la Arquitectura Hexagonal y TDD. Coordina los agentes de planificación, desarrollo y testing.
usage: "new-feature [nombre-funcionalidad] [descripcion]"
---

# Comando: Nueva Funcionalidad (New Feature)

Este comando orquesta la creación de nuevas funcionalidades en Griddo API, asegurando que se respeten las capas de Dominio, Aplicación e Infraestructura, el tipado estricto y la metodología TDD.

## Flujo de Implementación Orquestado

### 1. Investigación y Planificación (Planning Agent)
- Analizar el impacto en los módulos existentes.
- Crear un plan detallado en la raíz de `.claude/plans/` estructurado por capas hexagonales y siguiendo estrictamente la `.claude/development_guide.md`.
- El plan debe incluir `Aprobado: [ ]`.
- Inicializar la sesión en `.claude/sessions/active/`.
- **DETENCIÓN OBLIGATORIA**: Informar al usuario y esperar aprobación explícita del plan.
- **Activación**: Una vez aprobado, mover el plan a `.claude/plans/active/` mediante `/plan-manage activate`.
- **Agente**: `planning-agent`
- **Skills**: `system-architect`, `domain-expert`, `typescript`, `api-rest`

### 2. TDD - Fase Roja (Tester Agent)
- Escribir pruebas unitarias que fallen para los Casos de Uso (Vitest).
- Asegurar que el patrón AAA se cumpla.
- **Agente**: `tester`
- **Skills**: `tdd-champion`, `qa-engineer`, `vitest`, `typescript`

### 3. Desarrollo - Fase Verde (Developer Agent)
- Implementar la lógica mínima necesaria en el Caso de Uso para que los tests pasen.
- Implementar Entidades y DTOs según el plan.
- **Agente**: `developer`
- **Skills**: `domain-expert`, `usecase-developer`, `typescript`, `node-js`

### 4. Infraestructura y Persistencia (Developer Agent)
- Implementar Repositorios (TypeORM), Controladores (Express) y Rutas.
- Validar inputs con Zod.
- **Agente**: `developer`
- **Skills**: `infra-specialist`, `db-expert`, `express-js`, `typeorm`, `zod`, `mysql`

### 5. Integración y Calidad (Tester + Reviewer)
- Escribir y ejecutar tests de integración (Supertest).
- **Revisión**: El `reviewer` valida el cumplimiento de la arquitectura hexagonal y SOLID.
- **Auto-Corrección**: Si hay fallos, el desarrollador corrige (máximo 3 ciclos).
- **Agentes**: `tester`, `reviewer`, `debugger`
- **Skills**: `qa-engineer`, `supertest`, `code-reviewer`, `debug-master`

### 6. Documentación y Versiones (Writer Agent)
- **Documentación**: Actualizar `CHANGELOG.md` y JSDoc.
- **Versiones**: Proponer incremento de versión (SemVer).
- **Agente**: `writer`
- **Skills**: `technical-writer`

### 7. Finalización (Git)
- Crear rama: `git checkout -b feat/{feature-name}`.
- Commits siguiendo la convención: `feat: description`.

## Reglas Críticas
- **Seguimiento en TIEMPO REAL**: Actualizar el archivo de sesión en `.claude/sessions/active/` tras cada acción significativa.
- **Arquitectura Hexagonal**: Prohibido el acoplamiento de infraestructura en el dominio.
- **TDD Obligatorio**: No escribir código de producción sin tests previos.
- **Cero 'any'**: Mantener el tipado estricto.
- **Naming**: Seguir las convenciones (`save.usecase.ts`, `index.entity.ts`).

---

¿Cuál es la nueva funcionalidad que deseas planificar para Griddo API?
