// 文件夹批量导入文件(.mp3, .mpeg, .wav, .flac, or .opus only)
// 超过一定大小自动重新建立websocket并且断点续传(最大104857600 byte limit)
package main

import (
	"os"
	stt "speechToText"
)

func main() {
	path := os.Args[1]
	filelist := stt.GetFileList(path)
	for _, filepath := range filelist {
		client := stt.Client{
			Filepath: filepath,
		}
		client.Start()
	}

}
