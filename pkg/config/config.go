package config

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// N.B. We are currently using viper to manage configuration. viper takes care of merging configuration from different
// locations e.g. a file, environment variables, and command line flags. We then have viper unmarshal the configuration
// to the Configuration struct which is what we pass around in the application.

const (
	ConfigFlagName = "config"
	LevelFlagName  = "level"
	ConfigDir      = ".hccli"
)

// Config is the configuration data that gets persisted for kubedr.
//
// N.B. Right now the ondisk format and the in memory format are the same. We might want ot change that in the future
// to make it easier to change the on disk format and also to store derived values in memory that shouldn't be persisted
// to disk.
type Config struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion" yamltags:"required"`
	Kind       string `json:"kind" yaml:"kind" yamltags:"required"`

	// AIEndpoint is the endpoint of the model that turns natural language into queries
	AIEndpoint string `json:"aiEndpoint" yaml:"aiEndpoint"`

	// HoneycombAPIKeyFile contains the URI of the APIKey for HoneyComb
	HoneycombAPIKeyFile string `json:"honeycombAPIKeyFile" yaml:"honeycombAPIKeyFile"`

	Logging Logging `json:"logging" yaml:"logging"`
}

type Logging struct {
	Level string `json:"level" yaml:"level"`
}

func (c *Config) GetLogLevel() string {
	if c.Logging.Level == "" {
		return "info"
	}
	return c.Logging.Level
}

// GetConfigDir returns the configuration directory
func (c *Config) GetConfigDir() string {
	return filepath.Dir(viper.ConfigFileUsed())
}

// IsValid returns any errors with the configuration. The return is a list of configuration problems
func (c *Config) IsValid() []string {
	problems := make([]string, 0, 1)
	if c.HoneycombAPIKeyFile == "" {
		problems = append(problems, "No HoneycombAPIKeyFile key file specified. Please set one by running:\n\thccli config set honeycombApiKeyFile <path>")
	}
	return problems
}

// InitViper reads in config file and ENV variables if set.
// The results are stored inside viper. Call GetConfig to get a configuration.
// The cmd is passed in so we can bind to command flags
func InitViper(cmd *cobra.Command) error {
	// TODO(jeremy): Should we use a variable to ensure this is only called once?

	// TODO(jeremy): Should we be setting defaults?
	// see https://github.com/spf13/viper#establishing-defaults
	viper.SetEnvPrefix("hccli")
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath("$HOME/.hccli") // adding home directory as first search path

	// Makes overriding with environment variables work
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// Bind to the command line flag if it was specified.
	keyToflagName := map[string]string{
		ConfigFlagName:             ConfigFlagName,
		"logging." + LevelFlagName: LevelFlagName,
	}

	if cmd != nil {
		for key, flag := range keyToflagName {
			if err := viper.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
				return err
			}
		}
	}

	// We want to make sure the config file path gets set.
	// This is because we use viper to persist the location of the config file so can save to it.
	cfgFile := viper.GetString(ConfigFlagName)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		// TODO(jeremy): Is this the right semantics? We don't throw an error if the config file doesn't exist
		// We will just initialize an empty configuration. If the command requires configuration then the user
		// should get an error telling them to create the configuration.

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log := zapr.NewLogger(zap.L())
			log.Error(err, "Config file not found", "file", cfgFile)
			return nil
		}
		if _, ok := err.(*fs.PathError); ok {
			log := zapr.NewLogger(zap.L())
			log.Error(err, "Config file not found", "file", cfgFile)
			return nil
		}
		return err
	}
	return nil
}

// GetConfig returns the configuration instantiated from the viper configuration.
func GetConfig() *Config {
	// N.B. THis is a bit of a hacky way to load the configuration while allowing values to be overwritten by viper
	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal configuration; error %v", err))
	}

	return cfg
}

func binHome() string {
	log := zapr.NewLogger(zap.L())
	usr, err := user.Current()
	homeDir := ""
	if err != nil {
		log.Error(err, "Failed to get current user; falling back to temporary directory for homeDir", "homeDir", os.TempDir())
		homeDir = os.TempDir()
	} else {
		homeDir = usr.HomeDir
	}
	p := filepath.Join(homeDir, ConfigDir)

	return p
}

// Write writes the configuration to the specified file.
func (c *Config) Write(cfgFile string) error {
	log := zapr.NewLogger(zap.L())
	if cfgFile == "" {
		return errors.Errorf("No config file specified")
	}
	configDir := filepath.Dir(cfgFile)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Info("Creating config directory", "dir", configDir)
		if err := os.Mkdir(configDir, 0700); err != nil {
			return errors.Wrapf(err, "Failed to create config directory %s", configDir)
		}
	}

	f, err := os.Create(cfgFile)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(c)
}

func DefaultConfigFile() string {
	return binHome() + "/config.yaml"
}
