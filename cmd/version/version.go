// Package version implementa el comando para mostrar información de la versión.
//
// El comando version muestra el número de versión del CLI, junto con información adicional
// como el commit hash y la fecha de build cuando se usa el flag --verbose.
package version

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// Version es el número de versión del CLI.
// Se puede establecer en tiempo de compilación usando -ldflags.
var Version = "0.1.0"

// Commit es el hash del commit desde el cual se compiló el binario.
// Se establece en tiempo de compilación usando -ldflags.
var Commit = "unknown"

// BuildDate es la fecha en la que se compiló el binario.
// Se establece en tiempo de compilación usando -ldflags.
var BuildDate = "unknown"

// VersionOptions contiene las opciones configurables del comando version.
type VersionOptions struct {
	Short   bool // Short indica si se debe mostrar solo el número de versión
	JSON    bool // JSON indica si la salida debe ser en formato JSON
	Verbose bool // Verbose indica si se debe mostrar información detallada
}

// NewVersionCommand crea una nueva instancia del comando version.
// El comando version muestra la versión del CLI junto con información adicional si se solicita.
func NewVersionCommand() *cobra.Command {
	opts := &VersionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long: `Show the version of claude-init.

By default, shows the version number. With --verbose, also shows commit hash and build date.
With --json, outputs in JSON format.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(cmd, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Short, "short", "s", false, "Show only version number")
	cmd.Flags().BoolVarP(&opts.JSON, "json", "j", false, "Output in JSON format")
	cmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Show detailed information (commit, build date)")

	return cmd
}

// runVersion ejecuta la lógica principal del comando version.
// Muestra la información de versión en el formato solicitado.
func runVersion(cmd *cobra.Command, opts *VersionOptions) error {
	if opts.JSON {
		return outputJSON(cmd, opts)
	}
	return outputText(cmd, opts)
}

// outputJSON genera y muestra la información de versión en formato JSON.
func outputJSON(cmd *cobra.Command, opts *VersionOptions) error {
	output := map[string]interface{}{
		"version": Version,
	}

	// Añadir campos adicionales si no es short o si es verbose
	if !opts.Short || opts.Verbose {
		if Commit != "" && Commit != "unknown" {
			output["commit"] = Commit
		}
		if BuildDate != "" && BuildDate != "unknown" {
			output["build_date"] = BuildDate
		}
	}

	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// outputText genera y muestra la información de versión en formato texto legible.
func outputText(cmd *cobra.Command, opts *VersionOptions) error {
	if opts.Short {
		// Solo mostrar el número de versión
		fmt.Fprintln(cmd.OutOrStdout(), Version)
		return nil
	}

	// Salida estándar
	fmt.Fprintf(cmd.OutOrStdout(), "claude-init version %s\n", Version)

	// Información adicional verbose
	if opts.Verbose {
		if Commit != "" && Commit != "unknown" {
			fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", Commit)
		}
		if BuildDate != "" && BuildDate != "unknown" {
			fmt.Fprintf(cmd.OutOrStdout(), "built at: %s\n", BuildDate)
		}
	}

	return nil
}
