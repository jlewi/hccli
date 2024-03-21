package cmd

import (
	"fmt"
	"os"

	"github.com/jlewi/hccli/pkg"
	"github.com/jlewi/hccli/pkg/app"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/spf13/cobra"
)

// CreateQuery creates a command to generate queries
func CreateQuery() *cobra.Command {
	var nlq string
	var cols string
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

				p := pkg.Predictor{
					Config: app.Config,
				}

				resp, err := p.Predict(pkg.QueryInput{
					NLQ:  nlq,
					COLS: cols,
				})
				if err != nil {
					return err
				}
				if resp.Output != nil {
					fmt.Printf("The query is:\n%v\n", *resp.Output)
				} else {
					fmt.Printf("No query was returned:\n%v\n", resp.Error)
				}

				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&nlq, "query", "", "", "The honeycomb query")
	util.IgnoreError(cmd.MarkFlagRequired("query"))
	return cmd
}
