package slip

import (
	"encoding/json"
	"github.com/rameshvk/linkinfo"
)

type SlipJSONEncoderDecoder struct{}

func (pp *SlipJSONEncoderDecoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

func (pp *SlipJSONEncoderDecoder) Decode(b []byte) (interface{}, error) {
	if b == nil || len(b) == 0 {
		return nil, nil
	}
	var result linkinfo.Info
	err := json.Unmarshal(b, &result)
	return result, err
}
