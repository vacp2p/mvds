package transport

type WhisperNode struct {

}

func NewWhisperNode() *WhisperNode {
	return nil
}

func (n *WhisperNode) Tick() {
	panic("implement me")
}

func (n *WhisperNode) SendMessage(senderId []byte, to []byte, message []byte) {
	panic("implement me")
}
