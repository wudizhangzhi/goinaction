package speechToText

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// json -> struct
func JsonBytes2Struct(jsonBytes []byte, result interface{}) error {
	err := json.Unmarshal(jsonBytes, result)
	return err
}

// 结尾包含
func suffixContains(s []string, e string) bool {
	for _, a := range s {
		if strings.HasSuffix(e, a) {
			return true
		}
	}
	return false
}

// 获取文件列表
func GetFileList(path string) []string {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	var fileList []string
	switch mode := fi.Mode(); {
	case mode.IsDir():
		//
		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if info.Mode().IsRegular() && suffixContains(SuffixList, strings.ToLower(path)) {
					fileList = append(fileList, path)
				}
				return nil
			})
		if err != nil {
			panic(err)
		}
	case mode.IsRegular():
		//
		fileList = append(fileList, path)
	}
	return fileList
}
