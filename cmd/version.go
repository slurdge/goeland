package cmd

import (
	"fmt"

	"github.com/slurdge/goeland/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
	},
}

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Print the full changelog",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.ChangeLog)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(changelogCmd)
}
