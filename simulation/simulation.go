package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/status-im/mvds"
)

type Transport struct {
	sync.Mutex

	in  <-chan mvds.Packet
	out map[mvds.PeerId]chan<- mvds.Packet
}

func (t *Transport) Watch() *mvds.Packet {
	select {
		case p := <- t.in:
			return &p
	default:
		return nil
	}
}

func (t *Transport) Send(group mvds.GroupID, sender mvds.PeerId, peer mvds.PeerId, payload mvds.Payload) error {
	t.Lock()
	defer t.Unlock()

	c, ok := t.out[peer]
	if !ok {
		return errors.New("peer unknown")
	}

	c <- mvds.Packet{Group: group, Sender: sender, Payload: payload}
	return nil
}

func main() {

	ain := make(chan mvds.Packet, 10)
	bin := make(chan mvds.Packet, 10)
	cin := make(chan mvds.Packet, 10)

	at := &Transport{in: ain, out: make(map[mvds.PeerId]chan<- mvds.Packet)}
	bt := &Transport{in: bin, out: make(map[mvds.PeerId]chan<- mvds.Packet)}
	ct := &Transport{in: cin, out: make(map[mvds.PeerId]chan<- mvds.Packet)}

	group := groupId("meme kings")

	na := createNode(at, peerId())
	nb := createNode(bt, peerId())
	nc := createNode(ct, peerId())

	at.out[nb.ID] = bin
	at.out[nc.ID] = cin

	bt.out[na.ID] = ain
	bt.out[nc.ID] = cin

	ct.out[na.ID] = ain
	ct.out[nb.ID] = bin

	na.AddPeer(group, nb.ID)
	na.AddPeer(group, nc.ID)

	nb.AddPeer(group, na.ID)
	nb.AddPeer(group, nc.ID)

	nc.AddPeer(group, na.ID)
	nc.AddPeer(group, nb.ID)

	na.Share(group, nb.ID)
	na.Share(group, nc.ID)

	nb.Share(group, na.ID)
	nb.Share(group, nc.ID)

	nc.Share(group, na.ID)
	nc.Share(group, nb.ID)

	go na.Run()
	go nb.Run()
	go nc.Run()

	chat(group, na, nb)
}

func createNode(transport *Transport, id mvds.PeerId) *mvds.Node {
	ds := mvds.NewDummyStore()
	return mvds.NewNode(&ds, transport, Calc, id)
}

func chat(group mvds.GroupID, nodes ...*mvds.Node) {
	for {
		<-time.After(5 * time.Second)

		for _, n := range nodes {
			_, err := n.AppendMessage(group, []byte("test"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func Calc(count uint64, epoch int64) int64 {
	return epoch + int64(count*2)
}

func peerId() mvds.PeerId {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return mvds.PeerId(key.PublicKey)
}

func groupId(n string) mvds.GroupID {
	bytes := []byte(n)
	id := mvds.GroupID{}
	copy(id[:], bytes)
	return id
}
