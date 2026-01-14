package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCompletionCommand creates the completion command for generating shell completion scripts.
func NewCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for agentmgr.

To load completions:

Bash:
  $ source <(agentmgr completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ agentmgr completion bash > /etc/bash_completion.d/agentmgr
  # macOS:
  $ agentmgr completion bash > $(brew --prefix)/etc/bash_completion.d/agentmgr

Zsh:
  $ source <(agentmgr completion zsh)
  # To load completions for each session, execute once:
  $ agentmgr completion zsh > "${fpath[1]}/_agentmgr"

Fish:
  $ agentmgr completion fish | source
  # To load completions for each session, execute once:
  $ agentmgr completion fish > ~/.config/fish/completions/agentmgr.fish

PowerShell:
  PS> agentmgr completion powershell | Out-String | Invoke-Expression
  # To load completions for each session, add to your profile:
  PS> agentmgr completion powershell > agentmgr.ps1
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
