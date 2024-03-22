package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/jlewi/hccli/pkg"
	"github.com/jlewi/hccli/pkg/app"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/spf13/cobra"
)

// NewCreateQuery creates a command to generate queries
func NewCreateQuery() *cobra.Command {
	var dataset string
	var query string
	var queryFile string
	cmd := &cobra.Command{
		Use: "createquery",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := app.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

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
					return errors.Wrapf(err, "Error unmarshalling query")
				}

				hc, err := pkg.NewHoneycombClient(*app.Config)
				if err != nil {
					return err
				}

				qid, err := hc.CreateQuery(dataset, *hcq)
				if err != nil {
					return err
				}

				fmt.Printf("Created query :\n%v\n", qid)

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

	util.IgnoreError(cmd.MarkFlagRequired("dataset"))
	return cmd
}
