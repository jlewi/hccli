package cmd

import (
	"os"

	"github.com/jlewi/hccli/pkg/config"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var cfgFile string
	var level string
	var jsonLog bool
	rootCmd := &cobra.Command{
		Short: "hccli",
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, config.ConfigFlagName, "", "config file (default is $HOME/.hccli/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&level, config.LevelFlagName, "", "info", "The logging level.")
	rootCmd.PersistentFlags().BoolVarP(&jsonLog, "json-logs", "", false, "Enable json logging.")

	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.AddCommand(NewNLToQuery())
	rootCmd.AddCommand(NewVersionCmd("hccli", os.Stdout))
	return rootCmd
}
