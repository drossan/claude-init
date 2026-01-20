---
name: refactor
description: Refactoriza código existente del claude-init CLI para mejorar su calidad, legibilidad o mantenibilidad sin cambiar su funcionalidad externa.
usage: "refactor <descripcion-del-refactor>"
---

# Comando: Refactorización (Refactor)

Este comando orquesta la refactorización de código existente, mejorando la calidad interna sin cambiar la funcionalidad externa.

## Flujo de Refactorización

### 1. Análisis (Architect + Reviewer)
- Identificar código que necesita mejora
- Analizar patrones de diseño aplicables
- Considerar implicaciones del cambio

### 2. Planificación (Planning Agent)
- Crear un plan de refactorización
- Identificar pasos secuenciales
- Asegurar que los tests existentes protejan la funcionalidad

### 3. Refactorización (Developer)
- Aplicar cambios pequeños e incrementales
- Ejecutar tests después de cada cambio
- Mantener la funcionalidad externa idéntica

### 4. Verificación (Tester)
- Ejecutar suite completa de tests
- Verificar que no hay cambios en la funcionalidad
- Verificar que no hay regresiones

### 5. Documentación (Writer)
- Documentar cambios de arquitectura
- Actualizar godoc si es necesario

## Tipos Comunes de Refactorización

### 1. Extracción de Funciones
```go
// Antes
func process() {
    // mucho código...
}

// Después
func process() {
    data := readData()
    result := transformData(data)
    writeResult(result)
}
```

### 2. Introducción de Interfaces
```go
// Antes
func process(c *Client) error { ... }

// Después
type Processor interface {
    Process() error
}

func process(p Processor) error { ... }
```

### 3. Simplificación de Estructuras
```go
// Antes
type Config struct {
    Field1 string
    Field2 string
    Field3 string
    // muchos campos...
}

// Después
type Config struct {
    Common  CommonConfig
    Network NetworkConfig
    Storage StorageConfig
}
```

## Reglas Críticas

- **Sin Cambios Externos**: La API no debe cambiar
- **Tests Protegen**: Los tests existentes deben seguir pasando
- **Cambios Pequeños**: Hacer cambios incrementales
- **Tests Verdes**: Mantener los tests verdes todo el tiempo

---

¿Qué código deseas refactorizar?
