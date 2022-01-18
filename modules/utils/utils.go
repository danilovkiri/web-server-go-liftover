package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func MakeDir(path string) {
	err := os.Mkdir(path, 0755)
	if os.IsExist(err) {
		fmt.Println("### WARNING: Path already exists")
	} else if err != nil {
		log.Fatal("### ERROR:", err)
	}
}

func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println("### ERROR:", err)
		return
	}
	fmt.Printf("File %s successfully deleted", path)
}

func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic("### ERROR: ", err)
	}
	return dir
}

func GetFileSize(inputFile string) string {
	fi, err := os.Stat(inputFile)
	if err != nil {
		log.Panic("### ERROR: ", err)
	}
	return string(fi.Size())
}

func MakeCmdStruct(cwd string, inputFile string, outputFile string, sourceBuild string) exec.Cmd {
	cmdGo := exec.Cmd{
		Path:   cwd + "/" + "liftover/main.py",
		Args:   []string{cwd + "/" + "liftover/main.py", inputFile, outputFile, sourceBuild},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmdGo
}

func GetTempId(inputFile string) string {
	splitName := strings.Split(inputFile, "_upload")
	return splitName[1]
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CheckUploadedFileConformity(inputFile string) string {
	plausibleChromosomes := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11",
		"12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "X", "Y", "MT"}
	status := "ok"
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("### ERROR: ", err)
	}
	defer file.Close()
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		if strings.Split(reader.Text(), "\t")[0][0:1] == "#" {
			continue
		} else {
			oneLine := strings.Split(reader.Text(), "\t")
			if len(oneLine) != 4 {
				status := "invalid format"
				return status
			} else {
				chrom := oneLine[1]
				if !stringInSlice(chrom, plausibleChromosomes) {
					status := "invalid format"
					return status
				}
				pos := oneLine[2]
				_, err := strconv.Atoi(pos)
				if err != nil {
					status := "invalid format"
					return status
				}
			}
		}
	}
	return status
}
