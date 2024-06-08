package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/hccli/pkg"
	"github.com/jlewi/hccli/pkg/app"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// NewQueryToURL creates a command to turn queries into URLs
func NewQueryToURL() *cobra.Command {
	var dataset string
	var query string
	var queryFile string
	var baseURL string
	cmd := &cobra.Command{
		Use: "querytourl",
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

				if (query == "" && queryFile == "") || (query != "" && queryFile != "") {
					return errors.New("Exactly one of --query and --query-file must be specified")
				}

				if queryFile != "" {
					data, err := os.ReadFile(queryFile)
					if err != nil {
						return errors.Wrapf(err, "Error reading query file %v", queryFile)
					}
					query = string(data)
				}

				hcq := &pkg.HoneycombQuery{}

				if err := json.Unmarshal([]byte(query), hcq); err != nil {
					log.Error(err, "Error unmarshalling query", "query", query)
					return errors.Wrapf(err, "Error unmarshalling query")
				}

				hc, err := pkg.QueryToURL(*hcq, baseURL, dataset)
				if err != nil {
					return err
				}
				fmt.Printf("Honeycomb URL:\n%v\n", hc)

				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&query, "query", "", "", "The honeycomb query")
	cmd.Flags().StringVarP(&queryFile, "query-file", "", "", "A file containing the honeycomb query")
	cmd.Flags().StringVarP(&dataset, "dataset", "", "", "The dataset slug to create the query in")
	cmd.Flags().StringVarP(&baseURL, "base-url", "", "", "The base URL for your honeycomb URLs. It should be something like https://ui.honeycomb.io/${ORG}/environments/${ENVIRONMENT}")
	util.IgnoreError(cmd.MarkFlagRequired("dataset"))
	return cmd
}
