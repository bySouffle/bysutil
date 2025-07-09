package file

import (
	"io/fs"
	"os"
	"strings"
)

func Exist(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.Size() == 0 || !fileInfo.IsDir() {
		return err
	}
	return nil
}

func CreateFile(fullName string, perm fs.FileMode) (*os.File, error) {
	err := os.MkdirAll(string([]rune(fullName)[0:strings.LastIndex(fullName, "/")]), fs.ModeDir+0o775)
	if err != nil {
		return nil, err
	}
	//_, err = os.Stat(fullName)
	//if err == nil {
	//	// 文件已经存在
	//	// 处理文件已经存在的情况
	//	return os.Open(fullName)
	//} else if os.IsNotExist(err) {
	//	return os.Create(fullName)
	//}

	if err = os.Remove(fullName); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return os.Create(fullName)
}
