package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "memcardpro",
	Short: "",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//
		os.Exit(1)
	}
}
