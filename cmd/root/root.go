package root

import (
	"github.com/danielrossellosanchez/claude-init/cmd/completion"
	configcmd "github.com/danielrossellosanchez/claude-init/cmd/config"
	"github.com/danielrossellosanchez/claude-init/cmd/generate"
	initcmd "github.com/danielrossellosanchez/claude-init/cmd/init"
	"github.com/danielrossellosanchez/claude-init/cmd/version"
	"github.com/danielrossellosanchez/claude-init/internal/logger"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	log     *logger.Logger
)

var rootCmd = &cobra.Command{
	Use:   "claude-init",
	Short: "CLI para inicializar proyectos con configuración guiada por IA",
	Long: `claude-init es una herramienta CLI escrita en Go que permite
inicializar proyectos con una configuración optimizada para
desarrollo guiado por IA (Claude Code, etc.).

El CLI analiza el proyecto, hace preguntas interactivas, se conecta
a una API de IA (Claude, OpenAI, z.ai) y genera la estructura .claude/
necesaria para el proyecto.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configurar nivel de logging según el flag verbose
		if verbose && log != nil {
			log.SetLevel(logger.DEBUGLevel)
			log.Debug("Verbose mode enabled")
		}
	},
}

// Execute ejecuta el comando raíz con el logger proporcionado.
func Execute(l *logger.Logger) error {
	log = l
	return rootCmd.Execute()
}

// GetLogger retorna el logger configurado para uso de subcomandos.
func GetLogger() *logger.Logger {
	return log
}

// GetRootCmd retorna el comando raíz para testing.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Flag de verbosidad (persistente para todos los subcomandos)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (debug level)")

	// Añadir subcomandos
	generate.Execute(rootCmd, log)
	initcmd.Execute(rootCmd, log)

	// Añadir comandos adicionales
	rootCmd.AddCommand(version.NewVersionCommand())
	rootCmd.AddCommand(completion.NewCompletionCommand(rootCmd))
	rootCmd.AddCommand(configcmd.Cmd)
}
