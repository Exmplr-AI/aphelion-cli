package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `Generate shell completion script for the specified shell.

The completion script for bash can be generated with:
  aphelion completion bash

To load completions:

Bash:
  $ source <(aphelion completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ aphelion completion bash > /etc/bash_completion.d/aphelion
  # macOS:
  $ aphelion completion bash > /usr/local/etc/bash_completion.d/aphelion

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ aphelion completion zsh > "${fpath[1]}/_aphelion"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ aphelion completion fish | source

  # To load completions for each session, execute once:
  $ aphelion completion fish > ~/.config/fish/completions/aphelion.fish

PowerShell:
  PS> aphelion completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> aphelion completion powershell > aphelion.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	return cmd
}