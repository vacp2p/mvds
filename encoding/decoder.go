package encoding

import (
	"errors"
	"fmt"
	"io"

	"github.com/russolsen/transit"
	"github.com/status-im/mvds"
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

	payload := mvds.Payload{}
	for idx, v := range values {
		var ok bool

		switch idx {
		case 0:
			payload.Ack.Messages, ok = v.([]mvds.MessageID)
		case 1:
			payload.Offer.Messages, ok = v.([]mvds.MessageID)
		case 2:
			payload.Request.Messages, ok = v.([]mvds.MessageID)
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
