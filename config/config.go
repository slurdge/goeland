package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
	BindPFlag(key string, flag *pflag.Flag) error
}

// ReadDefaultConfig reads the configuration file
func ReadDefaultConfig(appName string, configName string) {
	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv()

	// global defaults

	viper.SetDefault("json-logs", false)
	viper.SetDefault("loglevel", "none")
	viper.SetDefault("dry-run", false)
	viper.SetDefault("email-timeout-ms", 5000)
	viper.SetDefault("purge-days", 15)
	viper.SetDefault("auto-purge", false)

	viper.SetConfigFile(configName)
	viper.AddConfigPath("$HOME/.goeland")
	viper.AddConfigPath(".")
	if ex, err := os.Executable(); err == nil {
		exPath := filepath.Dir(ex)
		viper.AddConfigPath(exPath)
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
