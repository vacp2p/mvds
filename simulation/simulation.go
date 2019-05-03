package main

import (
	"errors"

	"github.com/status-im/mvds"
)

type packet struct {
	sender  mvds.PeerId
	payload mvds.Payload
}

type Transport struct {
	in  <-chan packet
	out map[mvds.PeerId]chan<- packet
}

func (t *Transport) Watch() (mvds.PeerId, mvds.Payload) {
	p := <-t.in
	return p.sender, p.payload
}

func (t *Transport) Send(sender mvds.PeerId, peer mvds.PeerId, payload mvds.Payload) error {
	c, ok := t.out[peer]
	if !ok {
		return errors.New("peer unknown")
	}

	c <- packet{sender: sender, payload: payload}
	return nil
}

func main() {

	ain := make(chan packet)
	bin := make(chan packet)
	cin := make(chan packet)

	at := &Transport{in: ain, out: make(map[mvds.PeerId]chan<- packet)}
	bt := &Transport{in: bin, out: make(map[mvds.PeerId]chan<- packet)}
	ct := &Transport{in: cin, out: make(map[mvds.PeerId]chan<- packet)}

	group := groupId("meme kings")

	na := createNode(at, peerId("a"), group)
	nb := createNode(bt, peerId("b"), group)
	nc := createNode(ct, peerId("c"), group)

	at.out[nb.ID] = bin
	at.out[nc.ID] = cin

	bt.out[na.ID] = ain
	bt.out[nc.ID] = cin

	ct.out[na.ID] = ain
	ct.out[nb.ID] = bin

	na.AddPeer(nb.ID)
	na.AddPeer(nc.ID)

	nb.AddPeer(na.ID)
	nb.AddPeer(nc.ID)

	nc.AddPeer(na.ID)
	nc.AddPeer(nb.ID)

	na.Share(group, nb.ID)
	na.Share(group, nc.ID)

	nb.Share(group, na.ID)
	nb.Share(group, nc.ID)

	nc.Share(group, na.ID)
	nc.Share(group, nb.ID)

	go na.Run()
	go nb.Run()
	go nc.Run()

	chat(na, nb, nc)
}

func createNode(transport *Transport, id mvds.PeerId, groupID mvds.GroupID) mvds.Node {
	ds := mvds.NewDummyStore()
	return mvds.NewNode(&ds, transport, Calc, id, groupID)
}

func chat(nodes ...mvds.Node) {

}

func Calc(count uint64, lastTime int64) int64 {
	return lastTime + int64(count)*2
}

func peerId(n string) mvds.PeerId {
	bytes := []byte(n)
	id := mvds.PeerId{}
	copy(id[:], bytes)
	return id
}

func groupId(n string) mvds.GroupID {
	bytes := []byte(n)
	id := mvds.GroupID{}
	copy(id[:], bytes)
	return id
}