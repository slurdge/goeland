package cmd

import (
	"github.com/slurdge/goeland/internal/goeland/filters"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func purge(cmd *cobra.Command, args []string) {
	log.Debugln("Purging...")
	config := viper.GetViper()
	sources := config.GetStringMapString("sources")
	for source := range sources {
		filters.PurgeUnseen(config, source)
	}
}

// versionCmd represents the version command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge the database from old entries",
	Run:   purge,
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}
