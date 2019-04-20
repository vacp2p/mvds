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

func (n *WhisperNode) SendMessage(senderId []byte, to []byte, message []byte) {
	//n.whisper.Send()
}
