package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "run app as cli program",
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(demoEvtCmd)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
