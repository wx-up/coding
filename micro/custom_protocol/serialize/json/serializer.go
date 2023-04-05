package json

import (
	"encoding/json"

	"github.com/wx-up/coding/micro/custom_protocol/serialize"
)

type Serializer struct{}

func New() serialize.Serializer {
	return &Serializer{}
}

func (s *Serializer) Code() serialize.Code {
	return serialize.JsonCode
}

func (s *Serializer) Encode(val any) ([]byte, error) {
	return json.Marshal(val)
}

func (s *Serializer) Decode(data []byte, val any) error {
	return json.Unmarshal(data, val)
}
