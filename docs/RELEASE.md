# Guía de Release - claude-init CLI

Esta guía describe el proceso completo para crear releases del proyecto claude-init CLI.

## Tabla de Contenidos

1. [Versión Semántica](#versión-semántica)
2. [Proceso de Release Local](#proceso-de-release-local)
3. [Proceso de Release Automatizado](#proceso-de-release-automatizado)
4. [GitHub Actions](#github-actions)
5. [Verificación de Release](#verificación-de-release)
6. [Resolución de Problemas](#resolución-de-problemas)

## Versión Semántica

claude-init CLI sigue [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR**: Cambios incompatibles en la API
- **MINOR**: Nueva funcionalidad backwards-compatible
- **PATCH**: Corrección de bugs backwards-compatible

Formato: `v<MAJOR>.<MINOR>.<PATCH>[-<PRERELEASE>][+<BUILD>]`

Ejemplos:
- `v1.0.0` - Release estable
- `v1.2.3-beta.1` - Pre-release beta
- `v2.0.0-rc.1` - Release candidate

## Proceso de Release Local

### Prerrequisitos

1. **Working tree limpia**: Todos los cambios deben estar commiteados
2. **Tests pasando**: `make test-race` debe pasar
3. **Linting limpio**: `make lint` debe pasar

### Opción 1: Script Automatizado (Recomendado)

El script `scripts/release.sh` automatiza todo el proceso:

```bash
# Para release v0.1.0
./scripts/release.sh v0.1.0

# Para especificar un remote diferente
./scripts/release.sh v0.1.0 upstream
```

El script realiza:
1. Verifica que el working tree esté limpio
2. Ejecuta todos los tests
3. Ejecuta linters
4. Construye para todas las plataformas
5. Crea archivos de release
6. Genera checksums
7. Crea tag git
8. Pregunta si desea pushear el tag

### Opción 2: Manual

```bash
# 1. Verificar que todo esté limpio
git status

# 2. Ejecutar tests
make test-race

# 3. Ejecutar linters
make lint

# 4. Crear release (build + archives + checksums)
make release VERSION=v0.1.0

# 5. Verificar artefactos
ls -lh dist/

# 6. Crear tag anotado
git tag -a v0.1.0 -m "Release v0.1.0"

# 7. Pushear tag (trigger GitHub Actions)
git push origin v0.1.0
```

## Proceso de Release Automatizado

Cuando pusheas un tag (v*), GitHub Actions automáticamente:

1. **CI Pipeline**
   - Ejecuta tests con race detector
   - Ejecuta linters
   - Verifica que el código compile

2. **Build Pipeline**
   - Compila para todas las plataformas
   - Crea archivos `.tar.gz`
   - Genera checksums SHA256

3. **Release Pipeline**
   - Crea GitHub Release
   - Sube todos los artefactos
   - Genera release notes

### Plataformas Soportadas

| OS       | Arquitectura | Archivo Resultante               |
|----------|--------------|----------------------------------|
| Linux    | amd64        | `claude-init-v0.1.0-linux-amd64.tar.gz`   |
| Linux    | arm64        | `claude-init-v0.1.0-linux-arm64.tar.gz`   |
| macOS    | amd64        | `claude-init-v0.1.0-darwin-amd64.tar.gz`  |
| macOS    | arm64        | `claude-init-v0.1.0-darwin-arm64.tar.gz`  |
| Windows  | amd64        | `claude-init-v0.1.0-windows-amd64.tar.gz` |

## GitHub Actions

### CI Workflow (`.github/workflows/ci.yml`)

Se ejecuta en cada push y PR:

```yaml
# Eventos trigger
- push a main/master/develop
- pull requests

# Jobs
- test: Unit tests + race detector + coverage
- lint: go vet + gofmt + golangci-lint
- build: Build multi-plataforma
```

### Release Workflow (`.github/workflows/release.yml`)

Se ejecuta solo en tags:

```yaml
# Eventos trigger
- push de tags (v*)

# Jobs
- test: Tests + race detector
- build: Build multi-plataforma
- release: Crear GitHub Release
```

## Verificación de Release

### Antes de Hacer Release

```bash
# 1. Verificar versión
grep "Version =" cmd/version/version.go

# 2. Verificar changelog
cat CHANGELOG.md

# 3. Ejecutar tests completos
make test-all

# 4. Verificar build local
make build-all
```

### Después del Release

```bash
# 1. Verificar que GitHub Release se creó
gh release view v0.1.0

# 2. Verificar artefactos
gh release view v0.1.0 --json assets -q '.assets[].name'

# 3. Descargar y verificar un binario
curl -LO https://github.com/danielrossellosanchez/claude-init/releases/download/v0.1.0/claude-init-v0.1.0-linux-amd64.tar.gz
curl -LO https://github.com/danielrossellosanchez/claude-init/releases/download/v0.1.0/claude-init-v0.1.0-linux-amd64.tar.gz.sha256

# Verificar checksum
sha256sum -c claude-init-v0.1.0-linux-amd64.tar.gz.sha256

# Extraer y probar
tar -xzf claude-init-v0.1.0-linux-amd64.tar.gz
./claude-init-linux-amd64 version --verbose
```

## Resolución de Problemas

### Tag Ya Existe

```bash
# Eliminar tag local
git tag -d v0.1.0

# Eliminar tag remoto
git push origin :refs/tags/v0.1.0

# Reintentar
./scripts/release.sh v0.1.0
```

### Working Tree Sucio

```bash
# Ver cambios
git status

# Commitear cambios
git add .
git commit -m "Prepare for release v0.1.0"

# O stashear cambios temporalmente
git stash
./scripts/release.sh v0.1.0
git stash pop
```

### Tests Fallan

```bash
# Ejecutar tests verbosos
make test-race TEST_FLAGS=-v

# Si es un test específico
go test -v -race ./path/to/package

# Si es un race condition
go test -race -trace=trace.out ./path/to/package
```

### Build Falla

```bash
# Limpiar y reintentar
make clean-all
make build-all

# Verificar variables de versión
make help | grep VERSION
```

### GitHub Actions Falla

1. Ver el log del workflow:
   ```bash
   gh run list
   gh run view <run-id>
   ```

2. Verificar que el tag se pusheó correctamente:
   ```bash
   git ls-remote --tags origin
   ```

3. Re-ejecutar workflow:
   ```bash
   gh run rerun <run-id>
   ```

## Variables de Versión

Las siguientes variables se inyectan en build time:

- `VERSION`: Número de versión (ej: `v0.1.0`)
- `COMMIT`: Hash corto del commit (ej: `a1b2c3d`)
- `BUILD_DATE`: Fecha de build en UTC (ej: `2024-01-15T10:30:00Z`)

Puedes verlas con:

```bash
claude-init version --verbose
```

O en JSON:

```bash
claude-init version --json
```

## Checklist de Release

Antes de hacer release:

- [ ] Todos los tests pasan (`make test-race`)
- [ ] Linters pasan (`make lint`)
- [ ] Código formateado (`make fmt`)
- [ ] Working tree limpio
- [ ] Changelog actualizado
- [ ] Versión actualizada en `cmd/version/version.go`
- [ ] Release notes preparadas

Después del release:

- [ ] GitHub Release creada
- [ ] Artefactos subidos
- [ ] Checksums generadas
- [ ] Binarios descargados y verificados
- [ ] Documentation actualizada

## Próximos Pasos

Después de un release:

1. **Incrementar versión**: Actualizar `cmd/version/version.go`
2. **Crear rama de desarrollo** (si corresponde)
3. **Merge de hotfixes** (si corresponde)
4. **Actualizar CHANGELOG.md**

## Recursos

- [Semantic Versioning](https://semver.org/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Release Patterns](https://go.dev/doc/distribute)
