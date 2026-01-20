---
name: new-feature
description: Planifica e implementa una nueva funcionalidad en el claude-init CLI. Coordina los agentes de planificación, desarrollo y testing.
usage: "new-feature [nombre-funcionalidad] [descripcion]"
---

# Comando: Nueva Funcionalidad (New Feature)

Este comando orquesta la creación de nuevas funcionalidades en el claude-init CLI, asegurando que se respeten las mejores prácticas de Go y la arquitectura definida.

## Flujo de Implementación

### 1. Planificación (Planning Agent)
- Analizar el impacto de la nueva funcionalidad
- Crear un plan detallado en `.claude/plans/`
- Descomponer la tarea en pasos manejables
- Inicializar sesión en `.claude/sessions/active/`
- **DETENCIÓN OBLIGATORIA**: Esperar aprobación del plan

### 2. TDD - Fase Roja (Tester)
- Escribir tests que fallen para la nueva funcionalidad
- Definir los casos de uso principales
- Crear mocks para dependencias externas

### 3. Desarrollo - Fase Verde (Developer)
- Implementar la lógica mínima para que los tests pasen
- Seguir las convenciones de Go
- Manejar errores apropiadamente

### 4. Refactorización (Developer + Reviewer)
- Limpiar el código
- Aplicar patrones de diseño apropiados
- Revisar que cumple con las mejores prácticas

### 5. Integración (Tester)
- Escribir tests de integración
- Verificar que no hay regresiones
- Ejecutar suite completa de tests

### 6. Documentación (Writer)
- Añadir godoc comments
- Actualizar CHANGELOG.md
- Actualizar README si es necesario

## Reglas Críticas

- **Plan Aprobado**: No empezar a implementar sin un plan aprobado
- **TDD**: Escribir tests antes que el código
- **Go Idiomático**: Seguir las convenciones de Go
- **Cobertura**: Mantener >80% de cobertura
- **Documentación**: Documentar todo el código exportado

---

¿Cuál es la nueva funcionalidad que deseas implementar?
