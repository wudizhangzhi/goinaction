package speechToText

import (
	"encoding/json"
	"fmt"
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
