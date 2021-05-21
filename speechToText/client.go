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
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	WsConn      *websocket.Conn
	ErrorCh     chan interface{}
	StopCh      chan interface{}
	InterruptCh chan os.Signal
	FileBufCh   chan []byte
	Filepath    string
	BytesReaded int32
	Mut         sync.Mutex

	RespMap map[int]*MsgResponse
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

func (c *Client) RefreshConn() error {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	url := "wss://" + token.ServiceUrl[8:] + "/v1/recognize?model=en-US_BroadbandModel&access_token=" + token.AccessToken
	log.Println("创建连接")
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("websocket连接失败: %v", err)
		return err
	}
	c.Mut.Lock()
	c.closeWsConn()
	c.WsConn = wsConn
	c.BytesReaded = 0
	c.Mut.Unlock()

	c.hello()
	return nil
}

func (c *Client) hello() error {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(HelloMsg), &data); err != nil {
		log.Fatal(err)
		return err
	}
	//
	filetype := strings.ToLower(filepath.Ext(c.Filepath)[1:])
	data["content-type"] = "audio/" + filetype
	log.Printf("hello： %v", data)
	helloByte, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	c.WsConn.WriteMessage(websocket.TextMessage, helloByte)
	return nil
}

func (c *Client) receive() {
OuterLoop:
	for {
		c.Mut.Lock()
		msgType, message, err := c.WsConn.ReadMessage()
		c.Mut.Unlock()
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
			break OuterLoop
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
		c.WsConn.WriteMessage(websocket.BinaryMessage, buf)
		// c.FileBufCh <- buf
		time.Sleep(SleepDuration)
	}
	// allBytes, err := ioutil.ReadAll(f)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// total := 0
	// chunk := 0
	// // var buf []byte
	// var lastBuf []byte
	// for total < len(allBytes) {

	// 	// if (chunk+1)*BufSize > len(bytes) {
	// 	// 	buf = bytes[chunk*BufSize:]
	// 	// } else {
	// 	// 	buf = bytes[chunk*BufSize : (chunk+1)*BufSize]
	// 	// }
	// 	buf := allBytes[chunk*BufSize : (chunk+1)*BufSize]
	// 	if lastBuf != nil && bytes.Compare(lastBuf, buf) == 0 {
	// 		log.Println("一样！")
	// 	}
	// 	lastBuf = buf
	// 	// log.Printf("发送: %d kb\n", len(buf)/1024)
	// 	c.WsConn.WriteMessage(websocket.BinaryMessage, buf)
	// 	c.BytesReaded += int32(len(buf))
	// 	if c.BytesReaded+BufSize >= MaxBytes {
	// 		err := c.RefreshConn()
	// 		if err != nil {
	// 			break
	// 		}
	// 	}
	// 	chunk++
	// 	total += len(buf)
	// 	time.Sleep(SleepDuration)
	// }
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
	c.closeWsConn()
}

func (c *Client) Start() error {
	if err := c.RefreshConn(); err != nil {
		return err
	}
	c.ErrorCh = make(chan interface{})
	c.StopCh = make(chan interface{})
	c.InterruptCh = make(chan os.Signal, 1)
	c.FileBufCh = make(chan []byte)
	// c.FileResults = make([]string, 10)
	// c.Results = make([]Timestamp, 10)
	c.RespMap = make(map[int]*MsgResponse)

	signal.Notify(c.InterruptCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// go c.work()
	go c.readfile()
	go c.receive()
	defer c.WsConn.Close()
	c.finalize()
	return nil
}

func (c *Client) closeWsConn() error {
	if c.WsConn != nil {
		err := c.WsConn.Close()
		if err != nil && !websocket.IsCloseError(err) {
			log.Fatalf("关闭websocket失败: %v", err)
			return err
		}
	}
	return nil
}

func (c *Client) handleMsgRsp(rsp *MsgResponse) {
	// TODO
	// if rsp.ResultIndex != lastReultIndex {
	// 	if lastMsgResponse == nil {
	// 		lastMsgResponse = rsp
	// 	}
	// c.FileResults = append(c.FileResults, lastMsgResponse.Results[0].Alternatives[0].Transcript)
	// c.Results = append(c.Results, lastMsgResponse.Results[0].Alternatives[0].Timestamps...)
	// lastReultIndex = rsp.ResultIndex
	// f, err := os.OpenFile("output.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()
	// f.WriteString(" " + rsp.Results[0].Alternatives[0].Transcript)
	// sortedResult := ByTime(c.Results)
	// sort.Sort(sortedResult)
	// for _, i := range sortedResult {
	// 	fmt.Printf("%s ", i.Word)
	// }
	// fmt.Println(rsp.ResultIndex, lastReultIndex, rsp.Results[0].Alternatives[0].Transcript, c.FileResults)
	// lastMsgResponse = rsp
	// fmt.Println(c.FileResults)
	// }
	c.RespMap[rsp.ResultIndex] = rsp
	for i := 0; i < len(c.RespMap); i++ {
		fmt.Print(strings.TrimRight(c.RespMap[i].Results[0].Alternatives[0].Transcript, " ") + ". ")
	}
	fmt.Println("")
	// fmt.Println(rsp.ResultIndex, lastReultIndex, rsp.Results[0].Alternatives[0].Transcript, c.FileResults)
	// fmt.Println(strings.Join(c.FileResults, " "))
}

func (c *Client) handleStateRsp(rsp *StateResponse) {

}
