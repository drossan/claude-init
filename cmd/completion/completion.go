// Package completion implementa el comando para generar scripts de autocompletado.
//
// El comando completion genera scripts de autocompletado para diferentes shells (bash, zsh, fish, powershell).
// Estos scripts permiten al usuario autocompletar comandos y flags del CLI.
package completion

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

// NewCompletionCommand crea una nueva instancia del comando completion.
// El comando completion genera scripts de autocompletado para el shell especificado.
// Requiere el comando raíz del CLI para generar las completiones apropiadas.
func NewCompletionCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `Generate shell completion script.

To load completions:

Bash:
  $ source <(ia-start completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ ia-start completion bash > /etc/bash_completion.d/ia-start
  # macOS:
  $ ia-start completion bash > /usr/local/etc/bash_completion.d/ia-start

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ ia-start completion zsh > "${fpath[1]}/_ia-start"

  # You will need to start a new shell for this setup to take effect.

fish:
  $ ia-start completion fish | source

  # To load completions for each session, execute once:
  $ ia-start completion fish > ~/.config/fish/completions/ia-start.fish

PowerShell:
  PS> ia-start completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> ia-start completion powershell > ia-start.ps1
  # and source this file from your PowerShell profile.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompletion(rootCmd, args[0])
		},
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// runCompletion ejecuta la lógica principal del comando completion.
// Genera y muestra el script de autocompletado para el shell especificado.
func runCompletion(rootCmd *cobra.Command, shell string) error {
	// Normalizar el shell a minúsculas para comparación
	shellLower := strings.ToLower(strings.TrimSpace(shell))

	out := rootCmd.OutOrStdout()

	switch shellLower {
	case "bash":
		return rootCmd.GenBashCompletion(out)
	case "zsh":
		return rootCmd.GenZshCompletion(out)
	case "fish":
		return rootCmd.GenFishCompletion(out, true)
	case "powershell", "pwsh":
		return rootCmd.GenPowerShellCompletionWithDesc(out)
	default:
		return errors.New("unsupported shell type: " + shell + ". Supported shells: bash, zsh, fish, powershell")
	}
}
