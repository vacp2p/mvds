package mvds

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"

	"github.com/russolsen/transit"
)

const messageTag = "c5" // @todo

func NewMessageDecoder(r io.Reader) *transit.Decoder {
	decoder := transit.NewDecoder(r)
	decoder.AddHandler(messageTag, payloadHandler)
	return decoder
}

func payloadHandler(_ transit.Decoder, value interface{}) (interface{}, error) {
	taggedValue, ok := value.(transit.TaggedValue)
	if !ok {
		return nil, errors.New("not a tagged value")
	}
	values, ok := taggedValue.Value.([]interface{})
	if !ok {
		return nil, errors.New("tagged value does not contain values")
	}

	payload := Payload{}
	for idx, v := range values {
		ok := true

		switch idx {
		case 0:
			payload.Ack.Messages = convertToMessageIDArray(v.([]interface{}))
		case 1:
			payload.Offer.Messages = convertToMessageIDArray(v.([]interface{}))
		case 2:
			payload.Request.Messages = convertToMessageIDArray(v.([]interface{}))
		default:
			// skip any other values
			ok = true
		}

		if !ok {
			return nil, fmt.Errorf("invalid value for index: %d", idx)
		}
	}

	return payload, nil
}

// @todo some form of error handling
func convertToMessageIDArray(values []interface{}) []MessageID {
	m := make([]MessageID, 0)

	for _, v := range values {

		bytes := make([]byte, 0)

		for _, b := range v.(interface{}).([]interface{}) {
			bytes = append(bytes, b.([]byte)...)
		}

		var id MessageID
		copy(id[:], bytes[:32])

		m = append(m, id)

	}

	return m
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}