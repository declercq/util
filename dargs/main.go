package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"util/lib/commands"
	"util/lib/runner"
)

var debugF *bool
var inheritF *bool

func dcacheCmd(cmd *cobra.Command, args []string) {
	cmd.SilenceUsage = true
	targetDir := args[0]
	var d io.Writer
	if *debugF {
		d = os.Stderr
	}
	out := commands.DockerArgs(runner.New(d), *inheritF, targetDir)
	fmt.Println(out)
}

func main() {
	root := &cobra.Command{
		Use:   "dargs",
		Short: "dcache builds & caches docker utility images",
		Long: "dcache builds & caches docker utility images\n" +
			"The image will only be built if there is change detected in the Dockerfile or build context.\n" +
			"This is based on the git hash of the target directory.",
		Run:  dcacheCmd,
		Args: cobra.ExactArgs(1),
	}
	debugF = root.Flags().BoolP("debug", "d", false, "debug output")
	inheritF = root.Flags().BoolP("inherit", "i", false, "if already in a container, inherit from that container")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
