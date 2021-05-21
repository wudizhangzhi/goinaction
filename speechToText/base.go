package speechToText

import (
	"encoding/json"
	"time"
)

// 静态值
const (
	BufSize  = 1024 * 16
	MaxBytes = 104857600
	TokenUrl = "https://speech-to-text-demo.ng.bluemix.net/api/v1/credentials"
)

// 固定值
var (
	SleepDuration          = time.Microsecond * 200
	HelloMsg      string   = `{"timestamps":true,"content-type":"audio/mp3","interim_results":true,"keywords":["IBM","admired","AI","transformations","cognitive","Artificial Intelligence","data","predict","learn"],"keywords_threshold":0.01,"word_alternatives_threshold":0.01,"smart_formatting":true,"speaker_labels":false,"action":"start"}`
	StopMsg       string   = `{"action":"stop"}`
	SuffixList    []string = []string{".mp3", ".mpeg", ".wav", ".flac", ".opus"}
)

type TokenResp struct {
	AccessToken string `json:"accessToken"`
	ServiceUrl  string `json:"serviceUrl"`
}

// 具体词
type Timestamp struct {
	Word  string
	Start float32
	End   float32
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	a := []interface{}{&t.Word, &t.Start, &t.End}
	return json.Unmarshal(data, &a)
}

type Alternative struct {
	Transcript string      `json:"transcript"`
	Timestamps []Timestamp `json:"timestamps"`
}

type Result struct {
	Final        bool          `json:"final"`
	Alternatives []Alternative `json:"alternatives"`
}

type ByTime []Timestamp

func (r ByTime) Len() int { return len(r) }

func (r ByTime) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r ByTime) Less(i, j int) bool { return r[i].Start < r[j].Start }

type MsgResponse struct {
	ResultIndex int      `json:"result_index"`
	Results     []Result `json:"results"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StateResponse struct {
	State string `json:"state"`
}

func (e *MsgResponse) Valid() bool {
	return len(e.Results) > 0
}

func (e *ErrorResponse) Valid() bool {
	return e.Error != ""
}

func (e *StateResponse) Valid() bool {
	return e.State != ""
}

func loadResponse(respMsg []byte) (interface{}, error) {
	msgR := MsgResponse{}
	err := JsonBytes2Struct(respMsg, &msgR)
	if err == nil && msgR.Valid() {
		return &msgR, nil
	}
	errR := ErrorResponse{}
	err = JsonBytes2Struct(respMsg, &errR)
	if err == nil && errR.Valid() {
		return &errR, nil
	}
	stateR := StateResponse{}
	err = JsonBytes2Struct(respMsg, &stateR)
	if err == nil && stateR.Valid() {
		return &stateR, nil
	}
	return string(respMsg), err
}
