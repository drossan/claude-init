---
name: init
description: Inicializa un nuevo componente del CLI claude-init. Detecta el tipo de proyecto, crea la estructura de paquetes necesaria y genera los tests iniciales siguiendo TDD.
usage: "init <nombre-componente> <descripcion>"
---

# Comando: Init (Inicialización)

Este comando inicializa un nuevo componente del claude-init CLI, creando toda la estructura necesaria siguiendo las mejores prácticas de Go.

## Flujo de Implementación

### 1. Planificación (Planning Agent)
- Analizar el componente a implementar
- Crear un plan detallado en `.claude/plans/`
- El plan debe incluir la estructura de paquetes y las interfaces
- **DETENCIÓN OBLIGATORIA**: Informar al usuario y esperar aprobación

### 2. Diseño de Interfaces (Architect)
- Definir las interfaces necesarias
- Definir las estructuras de datos
- Establecer las dependencias entre paquetes

### 3. TDD - Tests Primero (Tester)
- Escribir tests para las interfaces
- Definir los casos de prueba principales
- Crear mocks para dependencias externas

### 4. Implementación (Developer)
- Implementar las interfaces definidas
- Seguir las convenciones de Go (gofmt, nombres)
- Manejar errores apropiadamente

### 5. Integración y Calidad (Reviewer + Tester)
- Revisión de código para asegurar calidad
- Ejecutar tests completos
- Verificar cobertura de código

### 6. Documentación (Writer)
- Añadir godoc comments
- Actualizar README si es necesario

## Reglas Críticas

- **Go Idiomático**: El código debe seguir "Effective Go"
- **Interfaces First**: Definir interfaces antes de implementar
- **TDD**: Escribir tests antes que el código
- **Error Handling**: Nunca ignorar errores
- **Godoc**: Documentar todos los exports

---

¿Qué componente deseas inicializar?
