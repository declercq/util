package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"util/lib/commands"
	"util/lib/runner"
)

var buildF *bool
var debugF *bool
var latestF *bool
var uploadF *bool
var sourceF *string

func dcacheCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	targetDir := args[0]
	registry := args[1]
	repository := args[2]
	if *sourceF == "" {
		*sourceF = targetDir + "/Dockerfile"
	}
	var d io.Writer
	if *debugF {
		d = os.Stderr
	}
	out, err := commands.DockerCache(targetDir, registry, repository, *buildF, runner.New(d), *latestF, *uploadF, *sourceF)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func main() {
	root := &cobra.Command{
		Use:   "dcache <dir> <registry> <repository>",
		Short: "dcache builds & caches docker utility images",
		Long: "dcache builds & caches docker utility images\n" +
			"The image will only be built if there is change detected in the Dockerfile or build context.\n" +
			"This is based on the git hash of the target directory.",
		RunE: dcacheCmd,
		Args: cobra.ExactArgs(3),
	}
	buildF = root.Flags().BoolP("build", "b", false, "always build")
	debugF = root.Flags().BoolP("debug", "d", false, "debug output")
	latestF = root.Flags().BoolP("latest", "l", false, "tag resulting image as 'latest' (will not push 'latest' tag)")
	uploadF = root.Flags().BoolP("upload", "u", false, "upload image if built and not from a dirty repo")
	sourceF = root.Flags().StringP("source", "s", "", "image source file ")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
