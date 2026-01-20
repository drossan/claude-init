---
name: improve-tests
description: Añade o mejora las pruebas unitarias e integración en Griddo API sin modificar el código de producción. Se centra en el patrón AAA y en aumentar la cobertura.
usage: "improve-tests [modulo/archivo] [objetivo-cobertura]"
---

# Comando: Mejorar Pruebas (Improve Tests)

Este comando está diseñado para fortalecer la suite de pruebas de Griddo API sin alterar el comportamiento de producción, coordinando a los agentes de calidad para alcanzar objetivos de cobertura.

## Flujo de Trabajo Orquestado

### 1. Investigación y Preparación (Tester Agent)
- Identificar áreas con baja cobertura o lógica crítica no probada.
- **Planificación**: Crear un plan de testing en la raíz de `.claude/plans/` con `Aprobado: [ ]`, siguiendo la estrategia de testing de la `.claude/development_guide.md`.
- **Sesión**: Inicializar la sesión en `.claude/sessions/active/`.
- **DETENCIÓN OBLIGATORIA**: Informar al usuario y esperar aprobación.
- **Activación**: Una vez aprobado, mover el plan a `.claude/plans/active/`.
- Crear una rama específica: `git checkout -b test/{nombre-del-test}`.
- **Agente**: `tester`
- **Skills**: `qa-engineer`, `tdd-champion`, `vitest`, `supertest`

### 2. Escritura de Tests (Tester Agent)
- Escribir pruebas siguiendo el patrón **AAA** (Arrange, Act, Assert).
- **Arrange**: Configurar mocks (repositorios, bus de eventos, etc.) y factories.
- **Act**: Ejecutar el Caso de Uso o Endpoint bajo prueba.
- **Assert**: Verificar resultados y efectos secundarios.
- **PROHIBIDO modificar código de producción**.
- **Agente**: `tester`
- **Skills**: `qa-engineer`, `vitest`, `supertest`, `typescript`

### 3. Verificación y Cobertura (Tester Agent)
- Ejecutar los tests y comprobar el aumento de cobertura.
- **Agente**: `tester`
- **Skills**: `vitest`, `supertest`

### 4. QA y Revisión (Reviewer Agent)
- Validar que los tests sean legibles, independientes y sigan las mejores prácticas.
- **Agente**: `reviewer`
- **Skills**: `code-reviewer`

### 5. Documentación y Git (Writer Agent)
- Commits siguiendo la convención: `test: description`.
- Actualizar guías de testing si se introducen nuevos patrones.
- **Agente**: `writer`
- **Skills**: `technical-writer`

## Reglas Críticas
- **Seguimiento en TIEMPO REAL**: Actualizar el archivo de sesión en `.claude/sessions/active/` tras cada acción significativa.
- **Integridad de Producción**: Bajo ninguna circunstancia se debe tocar el código de producción.
- **Independencia**: Cada test debe poder ejecutarse de forma aislada.
- **Determinismo**: Evitar depender de datos volátiles o servicios externos no mockeados.

---

¿Qué suite de pruebas de Griddo API vamos a fortalecer hoy?
