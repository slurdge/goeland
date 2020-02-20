package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	fmt.Println("Running...")
}

// versionCmd represents the version command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch the RSS and emails it",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(runCmd)
}
