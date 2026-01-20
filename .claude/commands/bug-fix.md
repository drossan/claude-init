---
name: bug-fix
description: Investiga y corrige errores en el claude-init CLI. Requiere la creación de un test de reproducción antes de aplicar la solución.
usage: "bug-fix <descripcion-del-error>"
---

# Comando: Corrección de Errores (Bug Fix)

Este comando orquesta la corrección de errores en el claude-init CLI, asegurando que se encuentre la causa raíz y se prevengan regresiones.

## Flujo de Corrección

### 1. Investigación (Debugger)
- Analizar el error reportado
- Identificar la causa raíz
- Recopilar información relevante (stack traces, logs)

### 2. Reproducción (Tester)
- Crear un test que reproduzca el error
- Verificar que el test falla
- Documentar los pasos para reproducir

### 3. Planificación (Planning Agent)
- Crear un plan para corregir el error
- Identificar qué cambios son necesarios
- Considerar efectos secundarios

### 4. Corrección (Developer)
- Implementar la solución
- Asegurar que el test de reproducción pasa
- No introducir nuevos errores

### 5. Verificación (Tester)
- Ejecutar suite completa de tests
- Verificar que no hay regresiones
- Añadir tests adicionales si es necesario

### 6. Documentación (Writer)
- Documentar el error y la solu
- Actualizar CHANGELOG.md

## Reglas Críticas

- **Test de Reproducción**: Obligatorio antes de corregir
- **Causa Raíz**: Entender y corregir la causa, no el síntoma
- **Sin Regresiones**: Todos los tests deben pasar
- **Documentación**: Documentar cambios significativos

---

¿Qué error necesitas corregir?
