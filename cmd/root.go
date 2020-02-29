package cmd

import (
	"fmt"
	"os"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goeland run",
	Short: "goeland is a simple rss to email program.",
	Long: `goeland is a simple rss to email program.
	
It was inspired by rss2email, but is an alternative with some cool features, such as filters.

The simple way to use it is to copy the provided config.toml.sample file and customize it.

The filters available are:
	- all: default, include all entries
	- none: include no entries
	- today: only includes entries with are today
	- digest: combines all the entries into one
	- digest2: combines all the entires into one, but with one less level (<h2> instead of <h1>)
	- combine: like digest, but use the first item as the title of the digest,
	- links: rewrite src="// and href="// to have an https:// prefix
	- lebrief: fetch full articles for LeBrief by NextINpact
	- wikipedia: remove unnecessary text from wikipedia entries`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//      Run: func(cmd *cobra.Command, args []string) { },
}

func fatalErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fatalErr(err)
	}
}

func createDefaultConfig(cfgFile string) {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		var configDefault pkging.File
		if configDefault, err = pkger.Open("/config.default.toml"); err != nil {
			fatalErr(fmt.Errorf("no config.toml and no default file present"))
		}
		info, err := configDefault.Stat()
		if err != nil {
			fatalErr(fmt.Errorf("cannot get default file stats"))
		}
		var configFile *os.File
		if configFile, err = os.Create(cfgFile); err != nil {
			fatalErr(fmt.Errorf("cannot open config.toml for writing"))
		}
		content := make([]byte, info.Size())
		configDefault.Read(content)
		configFile.Write(content)
	}
}

func initConfig() {
	createDefaultConfig(cfgFile)
	config.ReadDefaultConfig("goeland", cfgFile)
	log.SetDefaultLogger(log.NewLogger(viper.GetViper()))
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file (default is config.toml)")
}
