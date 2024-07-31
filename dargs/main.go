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

func dargsCmd(cmd *cobra.Command, args []string) {
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
		Short: "dargs generates default command line args for docker utility image use",
		Long: "dargs generates default command line args for docker utility image use\n" +
			"This includes mounting a local directory into the container and inheriting mounts where appropriate.",
		Run:  dargsCmd,
		Args: cobra.ExactArgs(1),
	}
	debugF = root.Flags().BoolP("debug", "d", false, "debug output")
	inheritF = root.Flags().BoolP("inherit", "i", false, "if already in a container, inherit from that container")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
