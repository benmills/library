package httpclient

import (
	"encoding/json"
)

type JsonData map[string]interface{}

func (data JsonData) Encode() string {
	encodedData, _ := json.Marshal(data)
	return string(encodedData);
}
