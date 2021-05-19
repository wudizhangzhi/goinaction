package speechToText

import (
	"fmt"
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
