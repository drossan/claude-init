---
name: pre-flight
description: Ejecuta verificaciones finales antes de hacer commit o release. Verifica que el build pasa, los tests están verdes y no hay cambios breaking.
usage: "pre-flight"
---

# Comando: Pre-Flight (Verificación Final)

Este comando ejecuta todas las verificaciones necesarias antes de considerar una tarea completa o antes de hacer commit.

## Verificaciones

### 1. Build
```bash
go build ./...
```
- [ ] El proyecto compila sin errores
- [ ] No hay warnings de compilación

### 2. Linters
```bash
gofmt -s -w .
golangci-lint run
```
- [ ] El código pasa gofmt
- [ ] El código pasa golangci-lint

### 3. Tests
```bash
go test ./...
go test -race ./...
go test -cover ./...
```
- [ ] Todos los tests pasan
- [ ] No hay race conditions
- [ ] Cobertura >80% para código crítico

### 4. Documentación
```bash
go doc -all ./...
```
- [ ] Todos los exports tienen godoc
- [ ] Los ejemplos en godoc son correctos

### 5. Cambios Breaking
- [ ] No hay cambios breaking sin version bump
- [ ] Los cambios están documentados en CHANGELOG

## Salida

El comando genera un reporte con el estado de cada verificación:

```
✓ Build: PASSED
✓ Linters: PASSED
✓ Tests: PASSED (coverage: 85%)
✓ Documentation: PASSED
✓ Breaking Changes: NONE

All checks passed! Ready to commit.
```

o

```
✗ Build: FAILED
  - main.go:45: undefined: SomeFunction

Fix the issues above before committing.
```

## Reglas Críticas

- **Todo Verde**: Todas las verificaciones deben pasar
- **Sin Warnings**: No debe haber warnings de linters
- **Cobertura**: Mantener cobertura >80%
- **Documentación**: Todo el código exportado debe estar documentado

---

Ejecutando verificaciones pre-flight...
