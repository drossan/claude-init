# Skill: Cobra CLI (cobra-cli)

## Propósito
Especialidad en el framework Cobra para crear aplicaciones CLI en Go. Cobra es usado por proyectos como Helm, kubectl, y docker.

## Responsabilidades
- Crear comandos y subcomandos usando Cobra
- Definir flags y argumentos correctamente
- Implementar comandos que retornan errores apropiadamente
- Usar PersistentFlags y LocalFlags apropiadamente

## Estructura Básica

### Comando Raíz

```go
package cmd

var rootCmd = &cobra.Command{
    Use:   "claude-init",
    Short: "CLI para inicializar proyectos con configuración guiada por IA",
    Long: `claude-init es una herramienta CLI escrita en Go que permite
inicializar proyectos con una configuración optimizada para
desarrollo guiado por IA (Claude Code, etc.).`,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runRoot(cmd, args)
    },
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    // Flags globales
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "modo verbose")
    rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "archivo de configuración")
}
```

### Subcomandos

```go
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Inicializa la configuración .claude en el proyecto actual",
    Long: `Inicializa la configuración .claude en el proyecto actual.

Este comando analiza el proyecto, hace preguntas interactivas,
y genera la estructura .claude/ necesaria para desarrollo guiado por IA.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runInit(cmd, args)
    },
}

func init() {
    rootCmd.AddCommand(initCmd)

    // Flags locales (solo para init)
    initCmd.Flags().BoolVar(&force, "force", false, "sobrescribir archivos existentes")
    initCmd.Flags().StringVar(&projectType, "type", "", "tipo de proyecto (auto, frontend, backend)")
}
```

## Patrones de Comandos

### 1. Comando con Validación

```go
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Configura las credenciales de la API de IA",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Validar que haya argumentos
        if len(args) < 2 {
            return cmd.Help()
        }
        return runConfig(cmd, args)
    },
}

func init() {
    rootCmd.AddCommand(configCmd)

    // Subcomandos
    configCmd.AddCommand(configGetCmd)
    configCmd.AddCommand(configSetCmd)
    configCmd.AddCommand(configUnsetCmd)
}
```

### 2. Comando Interactivo

```go
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Inicializa la configuración .claude",
    PreRunE: func(cmd *cobra.Command, args []string) error {
        // Validaciones antes de ejecutar
        if !isInGitRepo() {
            return errors.New("debes ejecutar este comando en un repositorio git")
        }
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        // Ejecutar flujo interactivo
        return runInteractiveInit(cmd.Context())
    },
}
```

### 3. Comando con Pre y Post Run

```go
var detectCmd = &cobra.Command{
    Use:   "detect",
    Short: "Detecta el tipo de proyecto",
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Se ejecuta antes del RunE del comando y sus subcomandos
        fmt.Println("Analizando proyecto...")
        return nil
    },
    PreRunE: func(cmd *cobra.Command, args []string) error {
        // Se ejecuta solo antes del RunE de este comando
        validatePath()
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        return runDetect(cmd, args)
    },
    PostRunE: func(cmd *cobra.Command, args []string) error {
        // Se ejecuta después del RunE
        fmt.Println("Detección completada")
        return nil
    },
}
```

## Flags

### PersistentFlags vs LocalFlags

```go
// ✅ PersistentFlags: disponible en este comando y todos sus subcomandos
rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "modo verbose")

// ✅ LocalFlags: solo disponible en este comando
initCmd.Flags().BoolVar(&force, "force", false, "sobrescribir archivos")

// ✅ Flags requeridos
initCmd.MarkFlagRequired("type")

// ✅ Flags mutuamente excluyentes
initCmd.MarkFlagsMutuallyExclusive("type", "auto-detect")
```

### Tipos de Flags

```go
// Bool
var verbose bool
cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "modo verbose")

// String
var output string
cmd.Flags().StringVarP(&output, "output", "o", "text", "formato de salida (text, json)")

// Int
var port int
cmd.Flags().IntVarP(&port, "port", "p", 8080, "puerto")

// StringSlice
var files []string
cmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "archivos a procesar")

// Count (para flags como -vvv)
var verboseLevel int
cmd.Flags().CountVarP(&verboseLevel, "verbose", "v", "nivel de verbose")
```

## Argumentos

### Validación de Argumentos

```go
var analyzeCmd = &cobra.Command{
    Use:   "analyze [path]",
    Short: "Analiza un proyecto",
    Args:  cobra.ExactArgs(1), // Requiere exactamente 1 argumento
    RunE: func(cmd *cobra.Command, args []string) error {
        path := args[0]
        return analyzeProject(path)
    },
}

// Otras opciones de Args:
// cobra.NoArgs()           // No acepta argumentos
// cobra.ArbitraryArgs()    // Acepta cualquier número de argumentos
// cobra.MinimumNArgs(1)    // Mínimo 1 argumento
// cobra.MaximumNArgs(2)    // Máximo 2 argumentos
// cobra.RangeArgs(1, 3)    // Entre 1 y 3 argumentos
```

### Argumentos con Validación Custom

```go
var validateCmd = &cobra.Command{
    Use: "validate [url]",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return errors.New("requiere exactamente 1 argumento")
        }
        if !isValidURL(args[0]) {
            return fmt.Errorf("'%s' no es una URL válida", args[0])
        }
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        return validateURL(args[0])
    },
}
```

## Help y Completions

### Custom Help

```go
// Sobrescribir el help template
cmd.SetHelpTemplate(`Custom help template...`)

// Añadir ejemplos al help
var exampleCmd = &cobra.Command{
    Use:     "example",
    Short:   "Comando de ejemplo",
    Example: `  claude-init example --type=react /path/to/project
  claude-init example --verbose`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // ...
    },
}
```

### Shell Completions

```go
var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Genera el script de completion",
    Long: `Para cargar completions en Bash:

  $ source <(claude-init completion bash)

  # Para cargar completions en cada sesión, ejecuta:
  $ claude-init completion bash > /etc/bash_completion.d/claude-init`,
    DisableFlagsInUseLine: true,
    ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
    RunE: func(cmd *cobra.Command, args []string) error {
        switch args[0] {
        case "bash":
            return rootCmd.GenBashCompletion(os.Stdout)
        case "zsh":
            return rootCmd.GenZshCompletion(os.Stdout)
        case "fish":
            return rootCmd.GenFishCompletion(os.Stdout, true)
        case "powershell":
            return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
        }
        return nil
    },
}
```

## Errores Comunes

### 1. No usar RunE

```go
// ❌ Mal: usando Run, no se pueden retornar errores
var badCmd = &cobra.Command{
    Use: "bad",
    Run: func(cmd *cobra.Command, args []string) {
        // ¿Qué hacer con el error?
    },
}

// ✅ Bien: usando RunE
var goodCmd = &cobra.Command{
    Use: "good",
    RunE: func(cmd *cobra.Command, args []string) error {
        return doSomething()
    },
}
```

### 2. No validar argumentos

```go
// ❌ Mal: no valida argumentos
var badCmd = &cobra.Command{
    Use: "bad [path]",
    RunE: func(cmd *cobra.Command, args []string) error {
        path := args[0] // panic si no hay argumentos
        return processPath(path)
    },
}

// ✅ Bien: valida argumentos
var goodCmd = &cobra.Command{
    Use:  "good [path]",
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        path := args[0]
        return processPath(path)
    },
}
```

## Checklist de Cobra

- [ ] Los comandos usan `RunE` para manejo de errores
- [ ] Los argumentos están validados con `Args`
- [ ] Los flags tienen descripciones claras
- [ ] Los flags persistentes se definen correctamente
- [ ] El help está completo (Short, Long, Example si es necesario)
- [ ] Los errores se propagan correctamente
- [ ] Se usan comandos anidados apropiadamente
