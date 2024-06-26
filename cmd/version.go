package cmd

import (
	"fmt"
	"io"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

// These constants are set by goreleaser
// https://goreleaser.com/cookbooks/using-main.version/?h=using+main.version
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func NewVersionCmd(name string, w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Return version",
		Example: fmt.Sprintf("%s  version", name),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(w, "%s %s, commit %s, built at %s by %s\n", name, version, commit, date, builtBy)
		},
	}
	return cmd
}

func logVersion() {
	log := zapr.NewLogger(zap.L())
	log.Info("binary version", "version", version, "commit", commit, "date", date, "builtBy", builtBy)
}
