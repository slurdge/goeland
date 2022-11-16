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
	runAtStartup := config.GetBool("run-at-startup")
	if runAtStartup {
		log.Infof("Running all the pipes once as requested")
		for pipe := range pipes {
			runPipe(pipe)
		}
	}
	scheduler := cron.New()
	found := 0
	for pipe := range pipes {
		schedule := config.GetString(fmt.Sprintf("pipes.%s.cron", pipe))
		if schedule != "" {
			log.Infof("Scheduling pipe:%s to run at: %s", pipe, schedule)
			pipe := pipe
			_, err := scheduler.AddFunc(schedule, func() { runPipe(pipe) })
			if err != nil {
				log.Warnf("Failed to schedule pipe:%s: %v", pipe, err)
				continue
			}
			found++
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
	daemonCmd.Flags().Bool("run-at-startup", false, "Run all the enabled pipes once at startup, before scheduling them")
	viper.BindPFlag("run-at-startup", daemonCmd.Flags().Lookup("run-at-startup"))
	bindFlags(daemonCmd, viper.GetViper())
	rootCmd.AddCommand(daemonCmd)
}
