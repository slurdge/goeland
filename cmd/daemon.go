package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runPipe(pipe string) {
	args := []string{pipe}
	run(runCmd, args)
}

func daemon(cmd *cobra.Command, args []string) {
	config := viper.GetViper()
	pipes := config.GetStringMapString("pipes")
	scheduler := cron.New()
	found := 0
	for pipe := range pipes {
		schedule := config.GetString(fmt.Sprintf("pipes.%s.cron", pipe))
		if schedule != "" {
			log.Infof("Scheduling pipe:%s to run at: %s", pipe, schedule)
			scheduler.AddFunc(schedule, func() { runPipe(pipe) })
			found += 1
		}
	}
	exitsignal := make(chan os.Signal, 1)
	signal.Notify(exitsignal, syscall.SIGINT, syscall.SIGTERM)
	log.Infof("Scheduled %d jobs", found)
	scheduler.Start()
	log.Infof("Press CTRL+C to exit program")
	<-exitsignal
	log.Infof("Ending daemon")
	context := scheduler.Stop()
	<-context.Done()
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Put the program in daemon mode",
	Long: `This command will run the program in the foreground indefinitely.
		The run command will be called according to the cron schedule defined in the configuration file`,
	Run: daemon,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
