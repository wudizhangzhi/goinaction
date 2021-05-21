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

func AudioSampleRate(f *os.File) (int, error) {
	d, err := mp3.NewDecoder(f)
	return d.SampleRate(), err
}

func saveToFile(output string, m map[int]map[int]*MsgResponse) error {
	f, err := os.OpenFile(output+".txt", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			sentance := strings.TrimRight(m[i][j].Results[0].Alternatives[0].Transcript, " ")
			if len(sentance) == 0 {
				continue
			}
			if !strings.HasSuffix(sentance, ".") {
				sentance += ". "
			}
			sentance = strings.Title(sentance)
			f.WriteString(sentance)
			f.WriteString("\n")
		}
	}
	return nil
}
