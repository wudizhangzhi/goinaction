package speechToText

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResult(t *testing.T) {
	testData := `{
		"result_index": 0,
		"results": [
		   {
			  "final": false,
			  "alternatives": [
				 {
					"transcript": "I have a general media okay ",
					"timestamps": [
					   [
						  "I",
						  4.01,
						  4.23
					   ],
					   [
						  "have",
						  4.23,
						  4.47
					   ],
					   [
						  "a",
						  4.47,
						  4.64
					   ],
					   [
						  "general",
						  4.64,
						  5.18
					   ],
					   [
						  "media",
						  5.57,
						  5.97
					   ],
					   [
						  "okay",
						  5.97,
						  6.28
					   ]
					]
				 }
			  ]
		   }
		]
	 }`
	r := MsgResponse{}
	err := JsonBytes2Struct([]byte(testData), &r)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("导出的结果： %+v", r)

}

func TestLoadHelloMsg(t *testing.T) {
	helloMsgString := `{"timestamps":true,"content-type":"audio/mp3","interim_results":true,"keywords":["IBM","admired","AI","transformations","cognitive","Artificial Intelligence","data","predict","learn"],"keywords_threshold":0.01,"word_alternatives_threshold":0.01,"smart_formatting":true,"speaker_labels":false,"action":"start"}`
	data, err := json.Marshal(helloMsgString)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
	helloMsgObj := map[string]interface{}{}
	err = json.Unmarshal([]byte(helloMsgString), &helloMsgObj)
	if err != nil {
		t.Error(err)
	}
	helloMsgObj["content-type"] = "audio/wav"
	fmt.Println(helloMsgObj)
	basename := "hello.blah"
	fmt.Println(filepath.Ext(basename)[1:])
}

func TestChunkFile(t *testing.T) {
	filepath := "F:/GoWorkplace/goinaction/speechToText/audios/test.wav"
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	// buf := make([]byte, BufSize)
	var count int8
	bufArray := make([][]byte, 0, 2)
	for count < 2 {
		buf := make([]byte, BufSize)
		// io.ReadFull(r, buf)
		// n, err := f.Read(buf)
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		bufArray = append(bufArray, buf)
		if n == 0 {
			if err != nil {
				t.Error(err)
			}
			if err == io.EOF {
				break
			}
		}
		count++
	}
	for i, buf := range bufArray {
		f, err := os.Create(fmt.Sprint(i) + ".txt")
		if err != nil {
			t.Error(err)
		}
		// f.Write(buf)
		// hex_string_data := hex.EncodeToString(buf)
		f.Write(buf)
	}
}

func TestReadFull(t *testing.T) {
	filepath := "F:/GoWorkplace/goinaction/speechToText/audios/test.wav"
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 2; i++ {
		b := bytes[i*BufSize : (i+1)*BufSize]
		f, err := os.Create(fmt.Sprint(i) + ".txt")
		if err != nil {
			t.Error(err)
		}
		f.Write(b)
	}
}
