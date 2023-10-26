package gitBranch

import (
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"io/fs"
	"kpk/common"
	"path/filepath"
	"strings"
)

var (
	dir        = gjson.Get(hiddenData, "dir").String()
	branchName string
)

var (
	//go:embed .hidden
	hiddenData string
	// Cmd describes an api command.
	Cmd = &cobra.Command{
		Use:  "git-branch",
		RunE: FindBranchFunc,
	}
)

func init() {
	var (
		apiCmdFlags = Cmd.Flags()
	)
	apiCmdFlags.StringVar(&branchName, "branch", "", "")
}

func FindBranchFunc(cmd *cobra.Command, args []string) error {
	if len(branchName) <= 0 {
		common.ErrorPrintLn("缺少分支名 [--branch]")
	}
	FindBranch(branchName)
	return nil
}

func FindBranch(branch string) {
	directories := common.FindFile(dir, func(entry fs.DirEntry, path string) (bool, string) {
		return entry.Name() == ".git", filepath.Dir(path)
	}, 2)

	for _, gitDir := range directories {
		lines := common.RunCommand(common.Command{
			Name: "git",
			Args: []string{"branch", "-r"},
			Dir:  gitDir,
		})
		for _, line := range lines {
			if strings.Contains(line, branch) {
				common.PinkPrintLn(line, gitDir)
			}
		}
	}
}
