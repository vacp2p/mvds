package securetransport

type Node interface {
	Send(data []byte) error
}
