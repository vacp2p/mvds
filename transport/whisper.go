package transport

import whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"

type WhisperNode struct {
	whisper whisper.Whisper
	topic   whisper.TopicType
}

func NewWhisperNode() *WhisperNode {
	return nil
}

func (n *WhisperNode) Watch() {
	panic("implement me")
}

func (n *WhisperNode) SendMessage(senderId []byte, to []byte, message []byte) error {
	msg, err := whisper.NewSentMessage(&whisper.MessageParams{
		TTL:      0,
		Src:      nil,
		Dst:      nil,
		KeySym:   nil,
		Topic:    n.topic,
		WorkTime: 0,
		PoW:      0,
		Payload:  nil,
		Padding:  nil,
	}) // @todo

	if err != nil {
		return err // @todo probably wrap into a new error before bubbling up
	}

	envelope := whisper.NewEnvelope(10, n.topic, msg)

	err = envelope.Seal(nil) // @todo
	if err != nil {
		return err // @todo probably wrap before bubbling up
	}

	err = n.whisper.Send(envelope)
	if err != nil {
		return err // @todo probably wrap before bubbling up
	}

	return nil
}
