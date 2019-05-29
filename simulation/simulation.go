package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	math "math/rand"
	"sync"
	"time"

	"github.com/status-im/mvds"
)

var offline int
var nodeCount int

type Transport struct {
	sync.Mutex

	in  <-chan mvds.Packet
	out map[mvds.PeerId]chan<- mvds.Packet
}

func (t *Transport) Watch() mvds.Packet {
	return <- t.in
}

func (t *Transport) Send(group mvds.GroupID, sender mvds.PeerId, peer mvds.PeerId, payload mvds.Payload) error {
	if math.Intn(100) < offline {
		return nil
	}

	c, ok := t.out[peer]
	if !ok {
		return errors.New("peer unknown")
	}

	c <- mvds.Packet{Group: group, Sender: sender, Payload: payload}
	return nil
}

func init() {
	flag.IntVar(&offline, "offline", 90, "percentage of node being offline")
	flag.Parse()

	flag.IntVar(&nodeCount, "nodes", 3, "amount of nodes")
	flag.Parse()
}

func main() {

	transports := make([]*Transport, 0)
	input := make([]chan mvds.Packet, 0)
	nodes := make([]*mvds.Node, 0)
	for i := 0; i < nodeCount; i++ {
		in := make(chan mvds.Packet)

		transport := &Transport{
			in: in,
			out: make(map[mvds.PeerId]chan<- mvds.Packet),
		}

		input = append(input, in)
		transports = append(transports, transport)
		nodes = append(nodes, createNode(transport, peerId()))
	}

	group := groupId("meme kings")
	// @todo add multiple groups
	// @todo maybe dms?

	// @todo allow for not all nodes to be peered and sharing to test how that looks
	for i, n := range nodes {
		for j, p := range nodes {
			if j == i {
				continue
			}

			transports[i].out[p.ID] = input[j]
			n.AddPeer(group, p.ID)
			n.Share(group, p.ID)
		}
	}

	for _, n := range nodes {
		n.Run()
	}

	// @todo configure how many nodes should be chatting
	chat(group, nodes[:len(nodes)-2]...)
}

func createNode(transport *Transport, id mvds.PeerId) *mvds.Node {
	ds := mvds.NewDummyStore()
	return mvds.NewNode(&ds, transport, Calc, id)
}

func chat(group mvds.GroupID, nodes ...*mvds.Node) {
	for {
		time.Sleep(5 * time.Second)

		for _, n := range nodes {
			_, err := n.AppendMessage(group, []byte(fmt.Sprintf("%x testing", n.ID)))
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
