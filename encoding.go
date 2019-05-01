package mvds

import (
	"bytes"
	"errors"
)



var (
	// ErrInvalidDecodedValue means that the decoded message is of wrong type.
	// This might mean that the status message serialization tag changed.
	ErrInvalidDecodedValue = errors.New("invalid decoded value type")
)


// DecodeMessage decodes a raw payload to Message struct.
func DecodePayload(data []byte) (message Payload, err error) {
	buf := bytes.NewBuffer(data)
	decoder := NewMessageDecoder(buf)
	value, err := decoder.Decode()
	if err != nil {
		return
	}

	message, ok := value.(Payload)
	if !ok {
		return message, ErrInvalidDecodedValue
	}
	return
}

// EncodeMessage encodes a Message using Transit serialization.
func EncodeMessage(value Payload) ([]byte, error) {
	var buf bytes.Buffer
	encoder := NewMessageEncoder(&buf)

	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
