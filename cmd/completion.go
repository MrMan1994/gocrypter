package cmd

import(
	"github.com/spf13/cobra"
	"gocrypter/log"
	"os"
)

func init() {
	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:                        "completion",
	Aliases:                    nil,
	SuggestFor:                 nil,
	Short:                      "",
	Long:                       "",
	Example:                    "",
	ValidArgs:                  []string{"bash", "zsh", "powershell"},
	Args:                       cobra.ExactValidArgs(1),
	ArgAliases:                 nil,
	PersistentPreRun:           nil,
	PreRun:                     nil,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "bash" {
			if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
				log.Fatal(err)
			}
		} else if args[0] == "zsh" {
			if err := rootCmd.GenZshCompletion(os.Stdout); err != nil {
				log.Fatal(err)
			}
		} else if args[0] == "powershell" {
			if err := rootCmd.GenPowerShellCompletion(os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
	},
	PostRun:                    nil,
	PersistentPostRun:          nil,
}
