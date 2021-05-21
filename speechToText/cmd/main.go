// 文件夹批量导入文件(.mp3, .mpeg, .wav, .flac, or .opus only)
// 超过一定大小自动重新建立websocket并且断点续传(最大104857600 byte limit)
package main

import (
	"fmt"
	"log"
	"os"
	stt "speechToText"
)

func main() {
	args := os.Args
	var path string
	if len(args) == 1 { // 默认当前目录
		path = "."
	} else {
		path = os.Args[1]
	}
	filelist := stt.GetFileList(path)
	defer func() {
		err := recover()
		log.Fatal(err)
	}()
	if len(filelist) == 0 {
		log.Fatalf("没有满足后缀的文件, 允许的后缀有: %s", stt.SuffixList)
	}
	for _, filepath := range filelist {
		client := stt.Client{
			Filepath: filepath,
		}
		err := client.Start()
		if err != nil {
			log.Fatalf("执行文件：%s 报错: %v", filepath, err)
		}

	}
	fmt.Println("Press the Enter Key to quit")
	fmt.Scanln()
}
