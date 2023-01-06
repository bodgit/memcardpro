package cmd

import (
	"github.com/bodgit/memcardpro/internal/split"
	"github.com/spf13/cobra"
)

var (
	useFlashID bool
	useSize    bool
)

func init() {
	splitCmd := &cobra.Command{
		Use:                   "split DIRECTORY FILE...",
		DisableFlagsInUseLine: true,
		Short:                 "Split generic virtual memory cards into multiple per-game cards",
		Long:                  ``,
		Args:                  cobra.MinimumNArgs(2), //nolint:gomnd
		RunE: func(cmd *cobra.Command, args []string) error {
			return split.MemoryCards(args[0], args[1:], useSize, useFlashID) //nolint:wrapcheck
		},
	}

	splitCmd.Flags().BoolVar(&useFlashID, "use-flash-id", false, "use the source memory card flash ID")
	splitCmd.Flags().BoolVar(&useSize, "use-size", false, "use the source memory card size")

	rootCmd.AddCommand(splitCmd)
}
