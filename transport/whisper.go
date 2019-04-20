package transport

import "github.com/ethereum/go-ethereum/whisper/whisperv6"

type WhisperNode struct {
	whisper whisperv6.Whisper
}

func NewWhisperNode() *WhisperNode {
	return nil
}

func (n *WhisperNode) Tick() {
	panic("implement me")
}

func (n *WhisperNode) SendMessage(senderId []byte, to []byte, message []byte) {
	// n.whisper.Send()
}
