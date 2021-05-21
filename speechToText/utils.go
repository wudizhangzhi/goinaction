package speechToText

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	mp3 "github.com/hajimehoshi/go-mp3"
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

func byteOfSecond(sec int, freq int) int {
	return sec * freq
}

func SplitMp3(filepath string, sec int) ([][]byte, error) {
	results := make([][]byte, 0)
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, err
	}
	log.Println("mp3 length: ", d.Length())
	var readed int64 = 0
	size := byteOfSecond(sec, d.SampleRate())
	buf := make([]byte, size)
	for readed < d.Length() {
		n, err := d.Read(buf)
		buf = buf[:n]
		if err != nil {
			return nil, err
		}
		readed += int64(n)
		results = append(results, buf)
	}
	// err = ioutil.WriteFile("output.mp3", buf, os.ModePerm)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	return results, nil
}
