package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var buildF *bool
var debugF *bool
var latestF *bool
var uploadF *bool
var sourceF *string

func run(cmd string) (string, error) {
	out, err := exec.Command(cmd).CombinedOutput()
	sout := string(out)
	if *debugF {
		_, _ = fmt.Fprintf(os.Stderr, "cmd: %s\n", cmd)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", sout)
	}
	if err != nil {
		return "", fmt.Errorf("cmd (%s) failed: %w - %s", cmd, err, sout)
	}
	return sout, nil
}

func runf(format string, a ...any) (string, error) {
	out, err := run(fmt.Sprintf(format, a...))
	if err != nil {
		return "", err
	}
	return out, err
}

func runfs(format string, a ...any) error {
	_, err := runf(format, a...)
	return err
}

func fatalErr(err error) {
	if err != nil {
		panic(err)
	}
}

func dcacheCmd(cmd *cobra.Command, args []string) (ret error) {
	defer func() {
		if r := recover(); r != nil {
			ret = r.(error)
		}
	}()
	dirty := false
	targetDir := args[0]
	registry := args[1]
	repository := args[2]
	targetDirHash, err := runf("git log -n1 --pretty=format:\"%%h\" %s", targetDir)
	fatalErr(err)
	wc, err := runf("git status --porcelain=v1 %s", targetDir)
	fatalErr(err)
	if len(wc) > 0 {
		*buildF = true
		dirty = true
	}
	dockerPath := fmt.Sprintf("%s/%s:%s", registry, repository, targetDirHash)
	date := time.Now().Format("20060102030405")
	dockerPathLong := fmt.Sprintf("%s.%s", dockerPath, date)
	dockerPathLatest := fmt.Sprintf("%s/%s:%s", registry, repository, "latest")
	if !dirty {
		// best effort to pull target image
		_ = runfs("docker pull %s", dockerPath)
		wc, err = runf("docker image ls -q %s", dockerPath)
		fatalErr(err)
		if len(wc) == 0 {
			*buildF = true
		}
	}
	if *buildF {
		fatalErr(runfs("docker build -f %s --tag %s %s", *sourceF, dockerPath, targetDir))
		if !dirty {
			fatalErr(runfs("docker tag %s %s", dockerPath, dockerPathLong))
			if *uploadF {
				fatalErr(runfs("docker push %s", dockerPath))
				fatalErr(runfs("docker push %s", dockerPathLong))
			}
		}
	}
	if *latestF {
		fatalErr(runfs("docker tag %s %s", dockerPath, dockerPathLatest))
	}
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
	sourceF = root.Flags().StringP("source", "s", "Dockerfile", "image source file ")

	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
