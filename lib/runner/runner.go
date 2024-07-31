package runner

import (
	"fmt"
	"io"
	"os/exec"
)

// This snippet can be used with FatalErr and Runfp to simulate 'bash set -e' behavior (fail/exit on error)
// Tried to make some reusable variant of this, but defer/recover are very picky
//
//	defer func() {
//	 	if r := recover(); r != nil {
// 			err = r.(error)
//		}
//	}()

type runner struct {
	debug io.Writer
}

// Runner utilities for commands
type Runner interface {
	Run(string) (string, error)
	Runf(string, ...interface{}) (string, error)
	Runfp(string, ...interface{}) string
}

// New Runner
func New(debug io.Writer) Runner {
	return &runner{debug: debug}
}

// FatalErr panics if passed err is not nil
func FatalErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Run cmd via shell, handle debug output if necessary, properly format errors
func (r *runner) Run(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	sOut := string(out)
	if r.debug != nil {
		_, _ = fmt.Fprintf(r.debug, "cmd: %s\n", cmd)
		_, _ = fmt.Fprintf(r.debug, "%s\n", sOut)
	}
	if err != nil {
		return "", fmt.Errorf("cmd (%s) failed: %w - %s", cmd, err, sOut)
	}
	return sOut, nil
}

// Runf operates like 'run' except it takes a Sprintf format and params to generate the command
func (r *runner) Runf(format string, a ...any) (string, error) {
	out, err := r.Run(fmt.Sprintf(format, a...))
	if err != nil {
		return "", err
	}
	return out, err
}

// Runfp operates like 'runf' except it panics on error
func (r *runner) Runfp(format string, a ...any) string {
	out, err := r.Runf(format, a...)
	FatalErr(err)
	return out
}
