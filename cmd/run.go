package cmd

import (
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
	"github.com/slurdge/indigo/config"
	"github.com/spf13/cobra"
)

const tmpURL = "https://www.nextinpact.com/rss/lebrief.xml"

func run(cmd *cobra.Command, args []string) {
	fmt.Println("Running...")
	defaultConfig := config.Config()

	// todo: remove this line
	fmt.Println(defaultConfig.GetBool("main.dry-run"))

	feeds := defaultConfig.GetStringMapString("feeds")
	for feedname, _ := range feeds {
		filename := defaultConfig.GetString(fmt.Sprintf("feeds.%s.filename", feedname))
		//todo : handle error
		file, _ := os.Open(filename)
		defer file.Close()
		fp := gofeed.NewParser()
		//todo: handle error
		feed, _ := fp.Parse(file)
		fmt.Println(feed.Title)
		//	for _, item := range feed.Items {
		//		fmt.Println(item.Title)
		//		fmt.Println(item.Description)
		//	}
	}

}

// versionCmd represents the version command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch the RSS and emails it",
	Run:   run,
}

func init() {
	runCmd.Flags().Bool("dry-run", false, "Do only a dry-run")
	config.Config().BindPFlag("main.dry-run", runCmd.Flags().Lookup("dry-run"))
	rootCmd.AddCommand(runCmd)
}
