---
name: pre-flight
description: Comando Gatekeeper para validación final de Griddo API antes de completar una tarea. Verifica el build, la integridad de los tipos y la ausencia de breaking changes en los contratos.
usage: "pre-flight"
---

# Comando: Pre-flight Check (Gatekeeper)

Este comando actúa como la última línea de defensa en Griddo API para garantizar que los cambios no rompan el build de producción ni la compatibilidad de los contratos de la API.

## Verificaciones Orquestadas

### 1. Build y Tipos (Developer Agent)
- **Acción**: Ejecutar `npm run build` (o comando equivalente).
- **Verificar**: Que no haya errores de compilación o transpilación en TypeScript.
- **Agente**: `developer`
- **Skills**: `fullstack-ts-expert`, `typescript`

### 2. Análisis de Breaking Changes (Reviewer Agent)
- **Acción**: Comparar el estado actual con la rama base para detectar cambios en las interfaces de Dominio y DTOs públicos.
- **Checklist**:
    - [ ] ¿Se ha eliminado o renombrado algún campo en un DTO de respuesta?
    - [ ] ¿Se ha añadido un parámetro obligatorio en un Caso de Uso existente?
    - [ ] ¿Se ha cambiado el endpoint de una ruta sin versionado?
- **Agente**: `reviewer`
- **Skills**: `system-architect`, `code-reviewer`, `api-rest`

### 3. Suite Completa de Tests (Tester Agent)
- **Acción**: Ejecutar `npm test` (o comando equivalente) para asegurar que nada se ha roto.
- **Agente**: `tester`
- **Skills**: `vitest`, `supertest`

## Flujo de Ejecución
1. **Análisis de Tipos**: Validación estricta de TypeScript.
2. **Construcción**: Verificación de que el proyecto es "buildable".
3. **Validación de Contratos**: Revisión de DTOs e Interfaces.
4. **Informe Final**:
    - ✅ **Pasa**: Todo correcto.
    - ⚠️ **Advertencia**: Cambios que podrían ser breaking changes.
    - ❌ **Error**: Fallo en build o tests.

## Reglas Críticas
- **Paso Obligatorio**: Ninguna funcionalidad o fix debe darse por concluido sin este chequeo.
- **Auto-corrección**: Si el fallo es trivial (formateo, imports), el desarrollador debe corregir automáticamente.
- **Documentación**: El resultado debe quedar registrado en la sesión.

---

¿Iniciamos la validación final de los cambios en Griddo API?
