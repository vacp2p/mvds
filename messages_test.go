package mvds

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"
)

func TestDecodePayload(t *testing.T) {
	id := [32]byte{}
	rand.Read(id[:])

	b := make([]MessageID, 0)
	b = append(b, MessageID(id))

	payload := Payload{
		Ack: Ack{ Messages: b},
		Offer: Offer{ Messages: b},
		Request: Request{ Messages: b},
		Messages: []Message{},
	}

	bytes, _ := EncodeMessage(payload)

	p, _ := DecodePayload(bytes)

	fmt.Println(p)
	fmt.Println(payload)

	if !reflect.DeepEqual(payload, p) {
		t.Errorf("Payloads do not match")
	}

	//fmt.Print(err)
	//print(bytes)
	//print(err)
}
