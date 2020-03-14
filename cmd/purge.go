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
	numOfDays := config.GetInt("purge-days")
	if numOfDays < 0 {
		numOfDays = 0
	}
	for source := range sources {
		filters.PurgeUnseen(config, source, numOfDays)
	}
}

// versionCmd represents the version command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge the database from old entries",
	Run:   purge,
}

func init() {
	purgeCmd.Flags().Int("purge-days", 15, "Number of days to keep for the purge command")
	viper.GetViper().BindPFlag("purge-days", purgeCmd.Flags().Lookup("purge-days"))
	rootCmd.AddCommand(purgeCmd)
}
