package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"kpk/common"
	gitBranch "kpk/git-branch"
	"runtime"
)

var (
	rootCmd = &cobra.Command{
		Use: "kpk",
	}
)

// Execute executes the given command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		common.ErrorPrintLn(err.Error())
	}
}

func init() {
	rootCmd.Version = fmt.Sprintf(
		"%s %s/%s", "0.0.0",
		runtime.GOOS, runtime.GOARCH)

	rootCmd.AddCommand(gitBranch.Cmd)
}
