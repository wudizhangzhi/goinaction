// 文件夹批量导入文件
// 超过一定大小自动重新建立websocket并且断点续传
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// fileBinaryCh chan []byte
	stopCh   chan interface{}
	helloMsg string = `{"timestamps":true,"content-type":"audio/wav","interim_results":true,"keywords":["IBM","admired","AI","transformations","cognitive","Artificial Intelligence","data","predict","learn"],"keywords_threshold":0.01,"word_alternatives_threshold":0.01,"smart_formatting":true,"speaker_labels":false,"action":"start"}`
	stopMsg  string = `{"action":"stop"}`
)

type TokenResp struct {
	AccessToken string `json:"accessToken"`
	ServiceUrl  string `json:"serviceUrl"`
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

func main() {
	interrupt := make(chan os.Signal, 1)
	done := make(chan interface{}, 1)
	signal.Notify(interrupt, os.Interrupt)

	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(token)

	url := "wss://" + token.ServiceUrl[8:] + "/v1/recognize?model=en-US_BroadbandModel&access_token=" + token.AccessToken
	fmt.Printf("连接： %v\n", url)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	defer c.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// 接受消息
	go func() {
		for {
			messageType, message, err := c.ReadMessage()
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
	c.WriteMessage(websocket.TextMessage, []byte(helloMsg))

	// 发送文件
	go func() {
		f, err := os.Open("./audios/test.wav")
		if err != nil {
			return
		}
		defer f.Close()
		r := bufio.NewReader(f)
		buf := make([]byte, 0, 4*1024)
		for {
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
			fmt.Printf("发送: %d\n", n)
			c.WriteMessage(websocket.BinaryMessage, buf)
			time.Sleep(500 * time.Millisecond)
		}
		c.WriteMessage(websocket.TextMessage, []byte(stopMsg))
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
