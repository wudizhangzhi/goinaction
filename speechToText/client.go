package speechToText

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
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var lastReultIndex int8 = -1

type Client struct {
	WsConn      *websocket.Conn
	ErrorCh     chan interface{}
	StopCh      chan interface{}
	InterruptCh chan os.Signal
	FileBufCh   chan []byte
	Filepath    string
	BytesReaded int32
	Mut         sync.Mutex

	FileResults []string
	Results     []Timestamp
}

func getAccessToken() (*TokenResp, error) {
	resp, err := http.Get(TokenUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tokenResp := TokenResp{}
	json.NewDecoder(resp.Body).Decode(&tokenResp)
	return &tokenResp, err
}

func (c *Client) RefreshConn() {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	// fmt.Println(token)

	url := "wss://" + token.ServiceUrl[8:] + "/v1/recognize?model=en-US_BroadbandModel&access_token=" + token.AccessToken
	log.Println("创建连接")
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	c.Mut.Lock()
	c.WsConn = wsConn
	c.BytesReaded = 0
	c.Mut.Unlock()
	//
	c.hello()
}

func (c *Client) hello() {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(HelloMsg), &data); err != nil {
		panic(err)
	}
	//
	filetype := strings.ToLower(filepath.Ext(c.Filepath)[1:])
	data["content-type"] = "audio/" + filetype
	log.Printf("hello： %v", data)
	helloByte, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	c.WsConn.WriteMessage(websocket.TextMessage, helloByte)
}

func (c *Client) receive() {
	for {
		msgType, message, err := c.WsConn.ReadMessage()
		// log.Printf("接收：%s\n", message)
		if err != nil {
			log.Printf("接收数据报错, 退出: %s", err)
			if c.ErrorCh != nil {
				c.ErrorCh <- err
			}
			break
		}

		txtMsg := message

		switch msgType {
		case websocket.TextMessage:
			//
		case websocket.BinaryMessage:
			// txtMsg, err = o.GzipDecode(message)
		}

		rsp, err := loadResponse(txtMsg)
		if err != nil {
			log.Fatal(err)
			break
		}
		switch msg := rsp.(type) {
		case *MsgResponse:
			c.handleMsgRsp(msg)
		case *ErrorResponse:
			log.Printf("ErrorResponse： %+v", rsp)
			c.ErrorCh <- rsp
		case *StateResponse:
			//
			log.Printf("StateResponse %+v", rsp)
		default:
			log.Fatal("没有匹配到类型！")
			break
		}
	}
}

func (c *Client) readfile() {
	log.Printf("开始上传: %s\n", c.Filepath)
	f, err := os.Open(c.Filepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	buf := make([]byte, 0, BufSize)
	for {
		if c.BytesReaded+BufSize >= MaxBytes {
			c.WsConn.Close()
			c.RefreshConn()
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
		c.BytesReaded += int32(n)
		// log.Printf("发送: %d\n", n)
		// c.WsConn.WriteMessage(websocket.BinaryMessage, buf)
		c.FileBufCh <- buf
		time.Sleep(SleepDuration)
	}

	c.WsConn.WriteMessage(websocket.TextMessage, []byte(StopMsg))
	log.Println("上传完成")
}

func (c *Client) work() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.keepAlive()
		case buf := <-c.FileBufCh:
			// log.Printf("发送: %d kb", len(buf)/1024)
			c.WsConn.WriteMessage(websocket.BinaryMessage, buf)
		case <-c.StopCh:
			return
		case <-c.InterruptCh:
			return
		}
	}

}

func (c *Client) keepAlive() {
	c.WsConn.WriteMessage(websocket.TextMessage, []byte("ping"))
}

func (c *Client) finalize() {
	defer func() {
		log.Println("结束!")
	}()

	select {
	case err := <-c.ErrorCh:
		log.Fatalf("服务器报错: %s", err)
	case <-c.StopCh:
	case <-c.InterruptCh:
	}

	close(c.StopCh)
	close(c.ErrorCh)
	close(c.InterruptCh)
	close(c.FileBufCh)
	if c.WsConn != nil {
		c.WsConn.Close()
	}
}

func (c *Client) Start() {
	c.RefreshConn()
	c.ErrorCh = make(chan interface{})
	c.StopCh = make(chan interface{})
	c.InterruptCh = make(chan os.Signal, 1)
	c.FileBufCh = make(chan []byte)
	c.FileResults = make([]string, 10)
	c.Results = make([]Timestamp, 10)

	signal.Notify(c.InterruptCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go c.work()
	go c.readfile()
	go c.receive()
	defer c.WsConn.Close()
	c.finalize()
}

func (c *Client) handleMsgRsp(rsp *MsgResponse) {
	// TODO
	if rsp.ResultIndex != lastReultIndex {
		c.FileResults = append(c.FileResults, rsp.Results[0].Alternatives[0].Transcript)
		c.Results = append(c.Results, rsp.Results[0].Alternatives[0].Timestamps...)
		lastReultIndex = rsp.ResultIndex
		// f, err := os.OpenFile("output.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer f.Close()
		// f.WriteString(" " + rsp.Results[0].Alternatives[0].Transcript)
		sortedResult := ByTime(c.Results)
		sort.Sort(sortedResult)
		for _, i := range sortedResult {
			fmt.Printf("%s ", i.Word)
		}
	}
	// fmt.Println(rsp.Results[0].Alternatives[0].Transcript)
	fmt.Println(strings.Join(c.FileResults, " "))
}

func (c *Client) handleStateRsp(rsp *StateResponse) {

}
