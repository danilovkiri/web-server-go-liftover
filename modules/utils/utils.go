package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func MakeCmdStruct(cwd string, inputFile string, outputFile string, sourceBuild string) exec.Cmd {
	cmdGo := exec.Cmd{
		Path:   cwd + "liftover/main.py",
		Args:   []string{cwd + "liftover/main.py", inputFile, outputFile, sourceBuild},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmdGo
}

func GetTempId(inputFile string) string {
	splitName := strings.Split(inputFile, "upload")
	return splitName[1]
}
