// 文件夹批量导入文件(.mp3, .mpeg, .wav, .flac, or .opus only)
// 超过一定大小自动重新建立websocket并且断点续传(最大104857600 byte limit)
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	bufSize  = 1024 * 6
	maxBytes = 104857600
)

var (
	// fileBinaryCh chan []byte
	stopCh      chan interface{}
	helloMsg    string   = `{"timestamps":true,"content-type":"audio/wav","interim_results":true,"keywords":["IBM","admired","AI","transformations","cognitive","Artificial Intelligence","data","predict","learn"],"keywords_threshold":0.01,"word_alternatives_threshold":0.01,"smart_formatting":true,"speaker_labels":false,"action":"start"}`
	stopMsg     string   = `{"action":"stop"}`
	suffixList  []string = []string{".mp3", ".mpeg", ".wav", ".flac", ".opus"}
	fileList    []string
	wsConn      *websocket.Conn
	bytesReaded int32
)

type TokenResp struct {
	AccessToken string `json:"accessToken"`
	ServiceUrl  string `json:"serviceUrl"`
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

func getAccessToken() (*TokenResp, error) {
	resp, err := http.Get("https://speech-to-text-demo.ng.bluemix.net/api/v1/credentials")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tokenResp := TokenResp{}
	json.NewDecoder(resp.Body).Decode(&tokenResp)
	return &tokenResp, err
}

func getFileList(path string) []string {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		//
		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if info.Mode().IsRegular() && suffixContains(suffixList, path) {
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

func RefreshConn() {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(token)

	url := "wss://" + token.ServiceUrl[8:] + "/v1/recognize?model=en-US_BroadbandModel&access_token=" + token.AccessToken
	fmt.Printf("连接： %v\n", url)
	wsConn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	bytesReaded = 0
}

func main() {
	// load all audios
	path := os.Args[1]
	fileList = getFileList(path)
	log.Printf("文件: %s\n", fileList)

	interrupt := make(chan os.Signal, 1)
	done := make(chan interface{}, 1)
	signal.Notify(interrupt, os.Interrupt)

	RefreshConn()
	defer wsConn.Close()

	// 接受消息
	go func() {
		for {
			messageType, message, err := wsConn.ReadMessage()
			if err != nil {
				stopCh <- err
				break
			}
			switch messageType {
			case websocket.TextMessage:
				fmt.Printf("接收: %s\n", message)
			default:
				fmt.Printf("接收: %s\n", message)
			}
		}
	}()
	// 开始
	wsConn.WriteMessage(websocket.TextMessage, []byte(helloMsg))

	// 发送文件
	go func() {
		for _, filepath := range fileList {
			log.Printf("开始上传: %s\n", filepath)
			f, err := os.Open(filepath)
			if err != nil {
				log.Fatal(err)
				continue
			}
			defer f.Close()
			r := bufio.NewReader(f)
			buf := make([]byte, 0, bufSize)
			for {
				if bytesReaded+bufSize >= maxBytes {
					RefreshConn()
				}
				n, err := r.Read(buf[:cap(buf)])
				buf = buf[:n]
				if n == 0 {
					if err == nil {
						continue
					}
					if err == io.EOF {
						break
					}
				}
				bytesReaded += int32(n)
				fmt.Printf("发送: %d\n", n)
				wsConn.WriteMessage(websocket.BinaryMessage, buf)
				time.Sleep(500 * time.Millisecond)
			}
		}

		wsConn.WriteMessage(websocket.TextMessage, []byte(stopMsg))
	}()

	for {
		select {
		case <-interrupt:
			fmt.Println("interrupt")
			return
		case <-done:
			fmt.Println("done")
			return
		}
	}
}
