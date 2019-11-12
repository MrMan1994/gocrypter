package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                "gocrypter",
	Aliases:            nil,
	SuggestFor:         nil,
	Short:              "a command line interface tool for file cryptography",
	Long:               "",
	Example:            "",
	ValidArgs:          nil,
	Args:               nil,
	ArgAliases:         nil,
	PersistentPreRun:   nil,
	PersistentPreRunE:  nil,
	PreRun:             nil,
	PreRunE:            nil,
	Run:                nil,
	RunE:               nil,
	PostRun:            nil,
	PostRunE:           nil,
	PersistentPostRun:  nil,
	PersistentPostRunE: nil,
}

// Execute the rootCmd
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
