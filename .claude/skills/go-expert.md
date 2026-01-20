# Skill: Go Expert (go-expert)

## Propósito
Garantizar código Go idiomático que sigue las convenciones y mejores prácticas de la comunidad Go.

## Responsabilidades
- Seguir "Effective Go" y las convenciones de la comunidad
- Usar `gofmt` y `golangci-lint` para asegurar estilo consistente
- Aplicar patrones de diseño idiomáticos de Go
- Evitar "habitos de otros lenguajes" en Go

## Convenciones de Nombres

### Paquetes
```go
// ✅ Bien: lowercase, sin guiones ni underscores
package detector
package config
package aiclient

// ❌ Mal: mezcla de mayúsculas, guiones
package detector
package config_loader
```

### Interfaces
```go
// ✅ Bien: verbos con sufijo -er
type Detector interface{}
type Generator interface{}
type Reader interface{}

// ✅ Bien: nombres compuestos aceptados
type http.ResponseWriter interface{}
type io.Reader interface{}

// ❌ Mal: nombres no verbales para interfaces
type Detect interface{}
type Generate interface{}
```

### Exportaciones vs Privados
```go
// ✅ Bien: PascalCase para exportados
type ProjectInfo struct {}
func (d *Detector) Detect() error {}

// ✅ Bien: camelCase para privados
type projectInfo struct {}
func (d *detector) detect() error {}
```

## Patrones Idiomáticos

### 1. Interfaces First

```go
// ✅ Bien: definir la interfaz primero
type Detector interface {
    Detect(path string) (ProjectInfo, error)
    Supports(path string) bool
}

// Implementación después
type NodeJSDetector struct {
    // ...
}

func (n *NodeJSDetector) Detect(path string) (ProjectInfo, error) {
    // ...
}
```

### 2. Error Handling

```go
// ✅ Bien: siempre devolver errores
func (d *Detector) Detect(path string) (ProjectInfo, error) {
    info, err := d.readFile(path)
    if err != nil {
        return ProjectInfo{}, fmt.Errorf("failed to read file: %w", err)
    }
    return info, nil
}

// ❌ Mal: ignorar errores
func (d *Detector) Detect(path string) (ProjectInfo, error) {
    info, _ := d.readFile(path)
    return info, nil
}
```

### 3. Defer para Cleanup

```go
// ✅ Bien: usar defer para cleanup
func (c *Client) Process() error {
    file, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer file.Close()

    // ... procesar archivo
}

// ✅ Bien: múltiples defers (LIFO)
func (c *Client) ProcessMultiple() error {
    db, err := sql.Open("sqlite3", "db.sqlite")
    if err != nil {
        return err
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // no-op si Commit() fue exitoso

    // ... procesar

    return tx.Commit()
}
```

### 4. Context para Operaciones Cancelables

```go
// ✅ Bien: aceptar context como primer parámetro
func (c *Client) GenerateRecommendations(ctx context.Context, req Request) (*Response, error) {
    if err := ctx.Err(); err != nil {
        return nil, err
    }

    resp, err := c.httpPost(ctx, url, req)
    if err != nil {
        return nil, fmt.Errorf("failed to generate: %w", err)
    }
    return resp, nil
}
```

### 5. Goroutines con WaitGroup

```go
// ✅ Bien: usar WaitGroup para esperar goroutines
func (p *Processor) ProcessItems(items []Item) error {
    var wg sync.WaitGroup
    errs := make(chan error, len(items))

    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()
            if err := p.process(i); err != nil {
                errs <- err
            }
        }(item)
    }

    wg.Wait()
    close(errs)

    for err := range errs {
        if err != nil {
            return err
        }
    }
    return nil
}
```

## Estructura de Archivos

```go
// ✅ Bien: orden estándar
package detector

// 1. Imports (agrupados y ordenados)
import (
    "fmt"
    "os"
)

// 2. Constants
const (
    DefaultTimeout = 30
    MaxRetries     = 3
)

// 3. Type definitions
type Detector struct {
    path    string
    timeout time.Duration
}

// 4. Interface methods
func (d *Detector) Detect() (ProjectInfo, error) {
    // ...
}

// 5. Private methods
func (d *Detector) readFile() (string, error) {
    // ...
}
```

## Errores Comunes a Evitar

### 1. No usar defer para cleanup

```go
// ❌ Mal
file, _ := os.Open("file.txt")
// procesar
file.Close()

// ✅ Bien
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()
```

### 2. Panic en código normal

```go
// ❌ Mal
func (d *Detector) Detect(path string) ProjectInfo {
    if path == "" {
        panic("path is empty")
    }
    // ...
}

// ✅ Bien
func (d *Detector) Detect(path string) (ProjectInfo, error) {
    if path == "" {
        return ProjectInfo{}, errors.New("path is empty")
    }
    // ...
}
```

### 3. No respetar el context

```go
// ❌ Mal: ignora el context
func (c *Client) Fetch(ctx context.Context) (*Response, error) {
    resp, _ := http.Get("https://api.example.com") // ignora ctx
    return parseResponse(resp), nil
}

// ✅ Bien: respeta el context
func (c *Client) Fetch(ctx context.Context) (*Response, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    return parseResponse(resp), nil
}
```

## Checklist de Go Idiomático

- [ ] El código pasa `gofmt -s -w .`
- [ ] El código pasa `golangci-lint run`
- [ ] Los errores nunca se ignoran
- [ ] Los recursos se liberan con `defer`
- [ ] Las operaciones que pueden tomar tiempo aceptan `context.Context`
- [ ] Los nombres siguen las convenciones (paquetes lowercase, exports PascalCase)
- [ ] Las interfaces se definen antes que las implementaciones
- [ ] No hay `panic` en código normal (solo en init)
- [ ] Los tests están en archivos `*_test.go`
- [ ] La documentación usa godoc format
