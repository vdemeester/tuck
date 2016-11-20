package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	gitsha = "unknown"
	buildDate = ""
)

type tuckOptions struct {
	version bool
}

// NewRootCommand returns a new root command
func NewRootCommand() *cobra.Command {
	var opts tuckOptions
	
	cmd:= &cobra.Command{
		Use: "tuck [OPTION ...] [-D|-S|-R] PACKAGE ... [-D|-S|-R] PACKAGE ...",
		Short: "symlink farm manager à-la-stow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTuck(opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.version, "version", false, "Print version and exit")
	return cmd
}

func runTuck(opts tuckOptions) error {
	if opts.version {
		printVersion()
		return nil
	}
	return nil
}

func printVersion() {
	fmt.Printf("tuck version %s (build: %s, date: %s)\n", version, gitsha, buildDate)
}
