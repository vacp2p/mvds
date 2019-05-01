package mvds

import (
	"io"
	"reflect"

	"github.com/russolsen/transit"
)

var (
	messageType                = reflect.TypeOf(Payload{})
	defaultMessageValueEncoder = &encoder{}
)

func NewMessageEncoder(w io.Writer) *transit.Encoder {
	encoder := transit.NewEncoder(w, false)
	encoder.AddHandler(messageType, defaultMessageValueEncoder)
	return encoder
}

type encoder struct {}

func (encoder) IsStringable(reflect.Value) bool {
	return true
}

func (encoder) Encode(e transit.Encoder, value reflect.Value, asString bool) error {
	payload := value.Interface().(Payload)

	messages := make([]interface{}, 0)
	for _, msg := range payload.Messages {
		m := []interface{}{
			msg.GroupID,
			msg.Timestamp,
			msg.Body,
		}

		messages = append(messages, m)
	}

	taggedValue := transit.TaggedValue{
		Tag: messageTag,
		Value: []interface{}{
			payload.Ack.Messages,
			payload.Offer.Messages,
			payload.Request.Messages,
			// @todo messages
		},
	}
	return e.EncodeInterface(taggedValue, false)
}
