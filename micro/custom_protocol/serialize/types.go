package serialize

type Serializer interface {
	Code() Code
	Encode(val any) ([]byte, error)
	// Decode 中的 val 是结构体指针
	Decode(data []byte, val any) error
}

type Code uint8

const (
	JsonCode Code = iota + 1
)
