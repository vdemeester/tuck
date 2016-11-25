package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vdemeester/tuck/process"
)

var (
	version   = "dev"
	gitsha    = "unknown"
	buildDate = ""
)

type tuckOptions struct {
	sourceDir string
	targetDir string
	tuck      bool // Equivalent to --stow
	delete    bool
	version   bool
}

// NewRootCommand returns a new root command
func NewRootCommand() *cobra.Command {
	var opts tuckOptions

	cmd := &cobra.Command{
		Use:   "tuck [OPTION ...] [-S|-D] PACKAGE ... [-S|-D] PACKAGE ...",
		Short: "symlink farm manager à-la-stow",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.version {
				printVersion()
				return nil
			}
			if len(args) == 0 {
				return fmt.Errorf("tuck requires at least 1 argument")
			}
			return runTuck(args, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.version, "version", false, "Print version and exit")
	flags.BoolVarP(&opts.tuck, "tuck", "S", true, "Tuck the package names that follow this option")
	flags.BoolVarP(&opts.delete, "delete", "D", false, "Untuck the package names that follow this option")
	flags.StringVarP(&opts.sourceDir, "dir", "d", "", "Set tuck dir to DIR (default is current dir)")
	flags.StringVarP(&opts.targetDir, "target", "t", "", "Set target to DIR (default is parent of tuck dir)")
	// FIXME(vdemeester) Add support for --ignore
	// FIXME(vdemeester) Add support for -n, --no, --simulate
	// FIXME(vdemeester) Add support for --verbose
	return cmd
}

// FIXME(vdemeester) support * as packages
func runTuck(args []string, opts tuckOptions) error {
	sourceDir := opts.sourceDir
	targetDir := opts.targetDir
	if sourceDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		sourceDir = cwd
	}
	if targetDir == "" {
		targetDir = filepath.Clean(filepath.Join(sourceDir, "../"))
	}
	mode := process.TuckMode
	if opts.delete {
		mode = process.UntuckMode
	}
	return process.Install(process.Config{
		Source:  sourceDir,
		Target:  targetDir,
		Modules: args,
		Mode:    mode,
	})
}

func printVersion() {
	fmt.Printf("tuck version %s (build: %s, date: %s)\n", version, gitsha, buildDate)
}
