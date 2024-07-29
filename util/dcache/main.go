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

func dcacheCmd(cmd *cobra.Command, args []string) error {
	dirty := false
	targetDir := args[0]
	registry := args[1]
	repository := args[2]
	targetDirHash, err := runf("git log -n1 --pretty=format:\"%%h\" %s", targetDir)
	if err != nil {
		return err
	}
	wc, err := runf("git status --porcelain=v1 %s", targetDir)
	if err != nil {
		return err
	}
	if len(wc) > 0 {
		*buildF = true
		dirty = true
	}
	dockerPath := fmt.Sprintf("%s/%s:%s", registry, repository, targetDirHash)
	date := time.Now().Format("20060102030405")
	dockerPathLong := dockerPath + "." + date
	dockerPathLatest := fmt.Sprintf("%s/%s:%s", registry, repository, "latest")
	if !dirty {
		// best effort to pull target image
		_, _ = runf("docker pull %s", dockerPath)
		wc, err = runf("docker image ls -q %s", dockerPath)
		if err != nil {
			return err
		}
		if len(wc) == 0 {
			*buildF = true
		}
	}
	if *buildF {
		_, err := runf("docker build -f %s --tag %s %s", *sourceF, dockerPath, targetDir)
		if err != nil {
			return err
		}
		if !dirty {
			_, err = runf("docker tag %s %s", dockerPath, dockerPathLong)
			if err != nil {
				return err
			}
			if *uploadF {
				_, err = runf("docker push %s", dockerPath)
				if err != nil {
					return err
				}
				_, err = runf("docker push %s", dockerPathLong)
				if err != nil {
					return err
				}
			}
		}
	}
	if *latestF {
		_, err = runf("docker tag %s %s", dockerPath, dockerPathLatest)
		if err != nil {
			return err
		}
	}
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
