package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/zapr"
	"github.com/jlewi/hccli/pkg"
	"github.com/jlewi/hccli/pkg/app"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewNLToQuery creates a command to generate queries
func NewNLToQuery() *cobra.Command {
	var nlq string
	var cols string
	var dataset string
	var output string
	cmd := &cobra.Command{
		Use: "nltoq",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := app.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				log := zapr.NewLogger(zap.L())

				logVersion()

				var translator pkg.Translator
				if app.Config.Replicate != nil {
					log.Info("Using Replicate translator")
					rc, err := pkg.NewReplicateClient(*app.Config)
					if err != nil {
						return err
					}
					translator = rc
				} else {
					log.Info("Using model on K8s")
					translator = &pkg.Predictor{
						Config: app.Config,
					}
				}

				hc, err := pkg.NewHoneycombClient(*app.Config)
				if err != nil {
					return err
				}
				if cols == "" {
					if dataset == "" {
						return errors.New("dataset must be specified if cols isn't specified")
					}
					log.Info("No columns specified; fetching columns from Honeycomb")
					columns, err := hc.GetColumns(dataset)
					if err != nil {
						return err
					}

					names := make([]string, 0, len(columns))

					for _, c := range columns {
						names = append(names, c.KeyName)
					}
					log.Info("Fetched list of columns", "names", names)

					b, err := json.Marshal(columns)
					if err != nil {
						return errors.Wrapf(err, "Failed to serialize columns")
					}
					cols = string(b)
				}

				queryStr, err := translator.Translate(pkg.QueryInput{
					NLQ:  nlq,
					COLS: cols,
				})
				if err != nil {
					return err
				}
				if queryStr != "" {
					fmt.Printf("The query is:\n%v\n", *queryStr)
					// Escaped query is to support copying the query inside a notebook to the command to create the
					// query
					// This is a bit of a hack. We replace ' with " so on the command line we can enclose the whole
					// thing in single quotes
					escaped := queryStr
					escaped = strings.Replace(escaped, "'", "\"", -1)
					fmt.Printf("Escaped query :\n%v\n", escaped)
				} else {
					fmt.Printf("No query was returned:\n%v\n", err)
				}

				if output != "" {
					if err := os.WriteFile(output, []byte(*queryStr), 0644); err != nil {
						return err
					}
					fmt.Printf("Wrote query to %v\n", output)
				}

				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&nlq, "nlq", "", "", "Natural language query")
	cmd.Flags().StringVarP(&cols, "cols", "", "", "Columns")
	cmd.Flags().StringVarP(&dataset, "dataset", "", "", "Honeycomb dataset to fetch columns for. Only required if cols isn't specified")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file to write the query to")
	util.IgnoreError(cmd.MarkFlagRequired("nlq"))
	return cmd
}
