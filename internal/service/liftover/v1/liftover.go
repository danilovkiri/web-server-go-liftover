// Package liftover provides liftover functionality via running a pyliftover script.
package liftover

import (
	"fmt"
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/exec"
)

// Liftover defines an object structure and attributes available for its methods.
type Liftover struct {
	fullpath string
	logger   *zerolog.Logger
}

// InitLiftover initializes a Liftover service.
func InitLiftover(logger *zerolog.Logger) *Liftover {
	liftover := &Liftover{
		fullpath: "internal/service/liftover/pyliftover/main.py",
		logger:   logger,
	}
	return liftover
}

// Convert38to19 runs hg38-to-hg19 conversion.
func (l *Liftover) Convert38to19(wd string, inputFile string, outputFile string) error {
	executableCmd := l.makeCmdStruct(wd, inputFile, outputFile, "hg38")
	err := executableCmd.Run()
	return err
}

// Convert19to38 runs hg19-to-hg38 conversion.
func (l *Liftover) Convert19to38(wd string, inputFile string, outputFile string) error {
	executableCmd := l.makeCmdStruct(wd, inputFile, outputFile, "hg19")
	err := executableCmd.Run()
	return err
}

// makeCmdStruct creates an exec.Cmd object ready to be executed.
func (l *Liftover) makeCmdStruct(wd string, inputFile string, outputFile string, sourceBuild string) *exec.Cmd {
	cmdGo := &exec.Cmd{
		Path:   wd + "/" + l.fullpath,
		Args:   []string{wd + "/" + l.fullpath, inputFile, outputFile, sourceBuild},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	l.logger.Info().Msg(fmt.Sprintf("Compiled shell command: %s", cmdGo.String()))
	log.Println()
	return cmdGo
}
