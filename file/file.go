package file

import (
	"os"
	"path/filepath"
	"strings"
)

// 获取指定目录下指定扩展名的所有文件路径
func GetFiles(filePath, extName string) []string {
	var files []string

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), extName) {
			files = append(files, filePath+"/"+info.Name())
		}
		return nil
	})

	if err != nil {
		return nil
	}
	return files
}
