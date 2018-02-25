package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type debugresponse struct {
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query,omitempty"`
	Body    interface{}         `json:"body"`
	Time    time.Time           `json:"time,omitempty"`
}

func (d *debugresponse) String() string {
	ba, err := json.Marshal(d)

	if err != nil {
		log.Println("Unable to encode data to JSON", err.Error(), fmt.Sprintf("%#v", d))
		return ""
	}

	return string(ba)
}

func parse(str string) (*debugresponse, error) {
	var dr debugresponse

	if err := json.Unmarshal([]byte(str), &dr); err != nil {
		log.Println("Unable to decode data from JSON", err.Error(), fmt.Sprintf("%q", str))
		return nil, err
	}

	return &dr, nil
}

func tostring(responses []debugresponse) string {
	ba, err := json.Marshal(responses)

	if err != nil {
		log.Println("Unable to convert slice of responses to string:", err.Error(), "\n", "Data:", fmt.Sprintf("%#v", responses))
		return ""
	}

	return string(ba)
}

func fromstring(str string) []debugresponse {
	var dr []debugresponse

	if err := json.Unmarshal([]byte(str), &dr); err != nil {
		log.Println("Unable to decode slice of debug response records:", err.Error(), "\n", "String:", str)
		return nil
	}

	return dr
}
