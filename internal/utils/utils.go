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

func MakeDir(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil && os.IsExist(err) {
		log.Printf("Path %s already exists\n", path)
		return nil
	}
	return err
}

func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("File %s successfully deleted\n", path)
}

func GetFileSize(inputFile string) string {
	fi, err := os.Stat(inputFile)
	if err != nil {
		log.Println(err)
		return "0"
	}
	return strconv.FormatInt(fi.Size(), 10)
}

func MakeCmdStruct(cwd string, inputFile string, outputFile string, sourceBuild string) exec.Cmd {
	cmdGo := exec.Cmd{
		Path:   cwd + "/../../" + "liftover/main.py",
		Args:   []string{cwd + "/../../" + "liftover/main.py", inputFile, outputFile, sourceBuild},
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
	defer file.Close()
	if err != nil {
		log.Fatal("### ERROR: ", err)
	}
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		if strings.Split(reader.Text(), "\t")[0][0:1] == "#" {
			continue
		}
		oneLine := strings.Split(reader.Text(), "\t")
		if len(oneLine) != 4 {
			status := "invalid format"
			return status
		}
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
	return status
}
