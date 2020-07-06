package file

import (
	"os"
	"path/filepath"
)

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func GetFilfAbsPath(fileName string) (string, error) {
	return filepath.Abs(fileName)
}


func WriteFile(saveFilePath string, body []byte) error {
	dir := filepath.Dir(saveFilePath)
	e := mkdir(dir)
	if e != nil {
		return e
	}
	file, err := os.OpenFile(saveFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func mkdir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}