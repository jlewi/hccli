package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jlewi/hccli/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// NewConfigCmd creates a command to configure various settings
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config",
	}

	cmd.AddCommand(NewConfigSetCmd())
	cmd.AddCommand(NewConfigGetCmd())
	return cmd
}

// NewConfigSetCmd creates a command to configure various settings
func NewConfigSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set <name>=<value>",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			err := func() error {
				if err := config.InitViper(cmd); err != nil {
					return err
				}

				parts := strings.Split(args[0], "=")
				if len(parts) < 2 {
					return errors.New("Invalid argument; argument is not in the form <name>=<value>")
				}
				name := parts[0]

				// N.B. We use a switch state because in the future if we have associated arrays we will need to
				// special case them because viper doesn't support them.
				var cfg *config.Config
				switch name {
				default:
					value := parts[1]
					viper.Set(name, value)
					cfg = config.GetConfig()
				}

				cfgFile := viper.ConfigFileUsed()
				if cfgFile == "" {
					cfgFile = config.DefaultConfigFile()
				}
				fmt.Printf("Writing configuration to %s\n", cfgFile)
				return cfg.Write(cfgFile)
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

// NewConfigGetCmd  creates a command to get the configuration
func NewConfigGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the configuration",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				if err := config.InitViper(cmd); err != nil {
					return err
				}
				cfg := config.GetConfig()

				if err := yaml.NewEncoder(os.Stdout).Encode(cfg); err != nil {
					return err
				}

				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}
