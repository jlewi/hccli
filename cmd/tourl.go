package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-logr/zapr"
	"github.com/jlewi/hccli/pkg"
	"github.com/jlewi/hccli/pkg/app"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewQueryToURL creates a command to turn queries into URLs
func NewQueryToURL() *cobra.Command {
	var dataset string
	var query string
	var queryFile string
	var baseURL string
	var open bool
	var outFile string
	var chromePort int
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
				if open {
					if err := browser.OpenURL(hc); err != nil {
						return errors.Wrapf(err, "Error opening URL %v", hc)
					}
				}
				if outFile != "" {
					if err := pkg.SaveHoneycombGraph(hc, outFile, chromePort); err != nil {
						return err
					}
				}
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
	cmd.Flags().StringVarP(&outFile, "out-file", "", "", "Save a PNG of the page to this file")
	cmd.Flags().IntVarP(&chromePort, "port", "", 9222, "Port chrome developer tools is running on. This only matters if you are saving a PNG of the page.")
	cmd.Flags().StringVarP(&baseURL, "base-url", "", "", "The base URL for your honeycomb URLs. It should be something like https://ui.honeycomb.io/${ORG}/environments/${ENVIRONMENT}")
	cmd.Flags().BoolVarP(&open, "open", "", false, "Open the URL in a browser")
	util.IgnoreError(cmd.MarkFlagRequired("dataset"))
	return cmd
}
