---
name: improve-tests
description: Mejora la cobertura y calidad de los tests existentes del claude-init CLI. Prohibido modificar código de producción con este comando.
usage: "improve-tests [paquete-o-funcion]"
---

# Comando: Mejora de Tests (Improve Tests)

Este comando mejora la calidad y cobertura de los tests existentes sin modificar el código de producción.

## Flujo de Mejora

### 1. Análisis (Tester)
- Identificar código sin cobertura suficiente
- Analizar qué casos de prueba faltan
- Revisar calidad de tests existentes

### 2. Planificación (Planning Agent)
- Crear un plan para mejorar los tests
- Identificar casos de prueba faltantes
- Priorizar código crítico

### 3. Implementación de Tests (Tester)
- Añadir tests para código sin cobertura
- Mejorar tests existentes
- Añadir tests de tabla para múltiples casos
- Crear mocks para dependencias externas

### 4. Verificación (Tester)
- Ejecutar suite completa de tests
- Verificar mejora de cobertura
- Asegurar que todos los tests pasan

## Tipos de Tests a Añadir

### 1. Tests de Unidad Faltantes
```go
func TestFunction_EdgeCase(t *testing.T) {
    // Test caso borde
}
```

### 2. Tests de Tabla
```go
func TestFunction_MultipleCases(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        // múltiples casos...
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test...
        })
    }
}
```

### 3. Tests de Integración
```go
func TestIntegration_CompleteFlow(t *testing.T) {
    // test del flujo completo
}
```

## Reglas Críticas

- **Sin Cambios de Producción**: Prohibido modificar código de producción
- **Cobertura**: Objetivo >80% para código crítico
- **Calidad**: Los tests deben ser claros y mantenibles
- **Mantener**: No eliminar tests existentes sin razón

---

¿Qué tests deseas mejorar?
