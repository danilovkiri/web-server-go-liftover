// Package liftover provides liftover functionality via running a pyliftover script.
package liftover

import (
	"log"
	"os"
	"os/exec"
)

// Liftover defines an object structure and attributes available for its methods.
type Liftover struct {
	fullpath string
}

// InitLiftover initializes a Liftover service.
func InitLiftover() (*Liftover, error) {
	liftover := &Liftover{
		fullpath: "internal/service/liftover/pyliftover/main.py",
	}
	return liftover, nil
}

// Convert38to19 runs hg38-to-hg19 conversion.
func (l *Liftover) Convert38to19(cwd string, inputFile string, outputFile string) error {
	executableCmd := l.makeCmdStruct(cwd, inputFile, outputFile, "hg38")
	err := executableCmd.Run()
	return err
}

// Convert19to38 runs hg19-to-hg38 conversion.
func (l *Liftover) Convert19to38(cwd string, inputFile string, outputFile string) error {
	executableCmd := l.makeCmdStruct(cwd, inputFile, outputFile, "hg19")
	err := executableCmd.Run()
	return err
}

// makeCmdStruct creates an exec.Cmd object ready to be executed.
func (l *Liftover) makeCmdStruct(cwd string, inputFile string, outputFile string, sourceBuild string) exec.Cmd {
	cmdGo := exec.Cmd{
		Path:   cwd + "/../../" + l.fullpath,
		Args:   []string{cwd + "/../../" + l.fullpath, inputFile, outputFile, sourceBuild},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	log.Println("Compiled shell command:", cmdGo.String())
	return cmdGo
}
