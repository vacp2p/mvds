package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"log"
	math "math/rand"
	"sync"
	"time"

	"github.com/status-im/mvds"
)

var (
	offline       int
	nodeCount     int
	communicating int
	sharing       int
	interval      int64
)

type Transport struct {
	sync.Mutex

	in  <-chan mvds.Packet
	out map[mvds.PeerId]chan<- mvds.Packet
}

func (t *Transport) Watch() mvds.Packet {
	return <-t.in
}

func (t *Transport) Send(group mvds.GroupID, sender mvds.PeerId, peer mvds.PeerId, payload mvds.Payload) error {
	math.Seed(time.Now().UnixNano())
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
	flag.IntVar(&offline, "offline", 90, "percentage of time a node is offline")
	flag.IntVar(&nodeCount, "nodes", 3, "amount of nodes")
	flag.IntVar(&communicating, "communicating", 2, "amount of nodes sending messages")
	flag.IntVar(&sharing, "sharing", 2, "amount of nodes each node shares with")
	flag.Int64Var(&interval, "interval", 5, "seconds between messages")
	flag.Parse()
}

func main() {

	// @todo validate flags

	transports := make([]*Transport, 0)
	input := make([]chan mvds.Packet, 0)
	nodes := make([]*mvds.Node, 0)
	for i := 0; i < nodeCount; i++ {
		in := make(chan mvds.Packet)

		transport := &Transport{
			in:  in,
			out: make(map[mvds.PeerId]chan<- mvds.Packet),
		}

		input = append(input, in)
		transports = append(transports, transport)
		nodes = append(nodes, createNode(transport, peerId()))
	}

	group := groupId()
	// @todo add multiple groups, only one or so nodes in every group so there is overlap
	// @todo maybe dms?

	for i, n := range nodes {
		peers := selectPeers(len(nodes), i, sharing)
		for _, p := range peers {
			peer := nodes[p].ID

			transports[i].out[peer] = input[p]
			n.AddPeer(group, peer)
			n.Share(group, peer)

			log.Printf("%x sharing with %x", n.ID.ToBytes()[:4], peer.ToBytes()[:4])
		}
	}

	for _, n := range nodes {
		n.Run()
	}

	chat(group, nodes[:communicating-1]...)
}

func selectPeers(nodeCount int, currentNode int, sharing int) []int {
	peers := make([]int, 0)

OUTER:
	for {
		if len(peers) == sharing {
			break
		}

		math.Seed(time.Now().UnixNano())
		i := math.Intn(nodeCount)
		if i == currentNode {
			continue
		}

		for _, p := range peers {
			if i == p {
				continue OUTER
			}
		}

		peers = append(peers, i)
	}

	return peers
}

func createNode(transport *Transport, id mvds.PeerId) *mvds.Node {
	ds := mvds.NewDummyStore()
	return mvds.NewNode(&ds, transport, Calc, id)
}

func chat(group mvds.GroupID, nodes ...*mvds.Node) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		for _, n := range nodes {
			_, err := n.AppendMessage(group, []byte(fmt.Sprintf("%x testing", n.ID)))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func Calc(count uint64, epoch uint64) uint64 {
	return epoch + (count*2)
}

func peerId() mvds.PeerId {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return mvds.PeerId(key.PublicKey)
}

func groupId() mvds.GroupID {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	id := mvds.GroupID{}
	copy(id[:], bytes)

	return id
}
