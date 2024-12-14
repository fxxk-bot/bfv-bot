package des

import (
	"bytes"
	"encoding/json"
)

func StringToStruct(jsonStr string, data interface{}) error {
	return ByteToStruct([]byte(jsonStr), data)
}

func ByteToStruct(jsonStr []byte, data interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(jsonStr))
	return decoder.Decode(data)
}
