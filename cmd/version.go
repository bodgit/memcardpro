package cmd

import "go.szostok.io/version/extension"

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("bodgit", "memcardpro"),
		),
	)
}
