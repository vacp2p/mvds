package mvds

type Node struct {
	Send     <-chan []byte
	Received chan<- []byte // @todo message type
}
