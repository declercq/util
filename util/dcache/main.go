package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var buildF *bool
var debugF *bool
var latestF *bool
var uploadF *bool
var sourceF *string

func dcacheCmd(cmd *cobra.Command, args []string) error {
	fmt.Printf("build: %t\n", *buildF)
	fmt.Printf("debug: %t\n", *debugF)
	fmt.Printf("latest: %t\n", *latestF)
	fmt.Printf("upload: %t\n", *uploadF)
	fmt.Printf("source: %s\n", *sourceF)
	fmt.Printf("args: %v\n", args)
	println("dcache called")
	return nil
}

func main() {
	root := &cobra.Command{
		Use:   "dcache",
		Short: "dcache builds & caches docker utility images",
		Long:  "TODO",
		RunE:  dcacheCmd,
		Args:  cobra.ExactArgs(3),
	}
	buildF = root.Flags().BoolP("build", "b", false, "always build")
	debugF = root.Flags().BoolP("debug", "d", false, "debug output")
	latestF = root.Flags().BoolP("latest", "l", false, "tag resulting image as 'latest' (will not push 'latest' tag)")
	uploadF = root.Flags().BoolP("upload", "u", false, "upload image if built and not from a dirty repo")
	sourceF = root.Flags().StringP("source", "s", "Dockerfile", "image source file ")

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
