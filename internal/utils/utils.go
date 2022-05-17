// Package utils provides miscellaneous functionality.
package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// MakeDir creates a directory specified by a path.
func MakeDir(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil && os.IsExist(err) {
		log.Printf("Path %s already exists\n", path)
		return nil
	}
	return err
}

// RemoveFile removes a file specified by a path.
func RemoveFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("File %s successfully deleted\n", path)
}

// GetFileSize provides byte size of a file specified by a path.
func GetFileSize(inputFile string) string {
	fi, err := os.Stat(inputFile)
	if err != nil {
		log.Println(err)
		return "0"
	}
	return strconv.FormatInt(fi.Size(), 10)
}

// GetTempId retrieves temp identifier of a temp file specified by a path.
func GetTempId(inputFile string) string {
	splitName := strings.Split(inputFile, "_upload")
	return splitName[1]
}

// stringInSlice returns True string is contained in a slice of strings.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// CheckUploadedFileConformity checks whether the provided file complies with specified conformity criteria.
func CheckUploadedFileConformity(inputFile string) string {
	plausibleChromosomes := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11",
		"12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "X", "Y", "MT"}
	file, err := os.Open(inputFile)
	defer file.Close()
	if err != nil {
		return "invalid format"
	}
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		if strings.Split(reader.Text(), "\t")[0][0:1] == "#" {
			continue
		}
		oneLine := strings.Split(reader.Text(), "\t")
		if len(oneLine) != 4 {
			return "invalid format"
		}
		chrom := oneLine[1]
		if !stringInSlice(chrom, plausibleChromosomes) {
			return "invalid format"
		}
		pos := oneLine[2]
		_, err := strconv.Atoi(pos)
		if err != nil {
			return "invalid format"
		}

	}
	return "ok"
}
