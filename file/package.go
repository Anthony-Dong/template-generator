package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NewTemplate(dir, new string) error {
	return ReplacePackage(dir, "go-demo", new, []string{"vendor", ".git", ".idea"})
}

func ReplacePackage(dir, old, new string, ignoreFile []string) error {
	files, err := GetAllFiles(dir, func(fileName string) bool {
		var pass = true
		if ignoreFile == nil || len(ignoreFile) == 0 {
			return true
		}
		for _, elem := range ignoreFile {
			if strings.Contains(fileName, elem) {
				pass = false
				break
			}
		}
		return pass
	})

	if err != nil {
		return err
	}
	for _, file := range files {
		err := ModifyFile(old, new, file)
		if err != nil {
			return err
		}
		fmt.Printf("change success : %s\n", file)
	}
	return nil
}

type FilterFile func(fileName string) bool

func GetAllFiles(dirPth string, filter FilterFile) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPth, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filter(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func ModifyFile(old, new string, fileName string) error {

	rfile, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer rfile.Close()
	reader := bufio.NewReader(rfile)
	fileLine := make([]string, 0)
	count := 0
	for {
		lines, isEOF, err := reader.ReadLine()
		if err != nil {
			if strings.Compare(err.Error(), "EOF") == 0 {
				break
			}
			return err
		}
		if isEOF {
			break
		}

		line := string(lines)
		if strings.Contains(line, old) {
			count++
			fileLine = append(fileLine, strings.ReplaceAll(string(lines), old, new))
		} else {
			fileLine = append(fileLine, line)
		}
	}
	if count == 0 {
		return nil
	}
	wfile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer wfile.Close()
	for _, elem := range fileLine {
		_, err := fmt.Fprintln(wfile, elem)
		if err != nil {
			return err
		}
	}
	return nil
}
