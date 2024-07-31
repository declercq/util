package commands

import (
	"fmt"
	"time"
	"util/lib/runner"
)

// DockerCache efficiently builds docker utility images, caching where possible
func DockerCache(targetDir string, registry string, repository string, build bool, r runner.Runner, latest bool,
	upload bool, source string) (out string, ret error) {
	defer func() {
		if r := recover(); r != nil {
			ret = r.(error)
		}
	}()

	// the image tag is based on the git hash of the source directory
	targetDirHash := r.Runfp("git log -n1 --pretty=format:\"%%h\" %s", targetDir)
	if targetDirHash == "" {
		return "", fmt.Errorf("unable to get git status of target directory: %s", targetDir)
	}
	dirty := false
	// determine if the target dir is dirty - always build if dirty as we don't know the real hash
	if len(r.Runfp("git status --porcelain=v1 %s", targetDir)) > 0 {
		build = true
		dirty = true
	}
	dockerPath := fmt.Sprintf("%s/%s:%s", registry, repository, targetDirHash)
	dockerPathLatest := fmt.Sprintf("%s/%s:%s", registry, repository, "latest")
	if !dirty {
		// best effort to pull target image - centrally cached is assumed to be best and should override local cache
		_, _ = r.Runf("docker pull %s", dockerPath)
		// check to see if the image is present - if not, build it
		if len(r.Runfp("docker image ls -q %s", dockerPath)) == 0 {
			build = true
		}
	} else {
		dockerPath += ".DIRTY"
	}

	if build {
		_ = r.Runfp("docker build -f %s --tag %s %s", source, dockerPath, targetDir)
		if !dirty {
			// tagging with the date helps us track the image source
			dockerPathLong := fmt.Sprintf("%s.%s", dockerPath, time.Now().Format("20060102030405"))
			_ = r.Runfp("docker tag %s %s", dockerPath, dockerPathLong)
			if upload {
				_ = r.Runfp("docker push %s", dockerPath)
				_ = r.Runfp("docker push %s", dockerPathLong)
			}
		}
	}
	if latest {
		r.Runfp("docker tag %s %s", dockerPath, dockerPathLatest)
	}
	return dockerPath, nil
}

// DockerArgs outputs default args to be used when running directories mounted into utility docker images
func DockerArgs(r runner.Runner, inherit bool, dir string) string {
	inheritS := ""
	// check to see if we are running in a container
	if inherit {
		inheritS, _ = r.Run("grep -e kubepods -e docker /proc/self/cgroup | head -n1 | cut -d: -f3 | rev | cut -d/ -f1 | rev")
	}
	if inheritS != "" {
		return fmt.Sprintf("--rm --workdir %s --volumes-from %s ", dir, inheritS)
	}
	return fmt.Sprintf("--rm --workdir /working/ --mount type=bind,source=\"%s\",target=/working/", dir)
}
