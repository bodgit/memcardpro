package cmd

import "go.szostok.io/version/extension"

func init() {
	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("bodgit", "memcardpro"),
		),
	)
}
