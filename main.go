package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	math "math/rand"
	"time"

	"github.com/status-im/mvds/node"
	"github.com/status-im/mvds/protobuf"
	"github.com/status-im/mvds/state"
	"github.com/status-im/mvds/store"
	"github.com/status-im/mvds/transport"
)

var (
	offline       int
	nodeCount     int
	communicating int
	sharing       int
	interval      int64
	interactive   int
)

func init() {
	flag.IntVar(&offline, "offline", 90, "percentage of time a node is offline")
	flag.IntVar(&nodeCount, "nodes", 3, "amount of nodes")
	flag.IntVar(&communicating, "communicating", 2, "amount of nodes sending messages")
	flag.IntVar(&sharing, "sharing", 2, "amount of nodes each node shares with")
	flag.Int64Var(&interval, "interval", 5, "seconds between messages")
	flag.IntVar(&interactive, "interactive", 3, "amount of nodes to use INTERACTIVE mode, the rest will be BATCH") // @todo should probably just be how many nodes are interactive
	flag.Parse()
}

func main() {

	// @todo validate flags

	transports := make([]*transport.ChannelTransport, 0)
	input := make([]chan transport.Packet, 0)
	nodes := make([]*node.Node, 0)
	for i := 0; i < nodeCount; i++ {
		in := make(chan transport.Packet)

		t := transport.NewChannelTransport(offline, in)

		input = append(input, in)
		transports = append(transports, t)

		mode := node.INTERACTIVE
		if i+1 >= interactive {
			mode = node.BATCH
		}

		nodes = append(
			nodes,
			createNode(t, peerID(), mode),
		)
	}

	group := groupId()
	// @todo add multiple groups, only one or so nodes in every group so there is overlap
	// @todo maybe dms?

	for i, n := range nodes {
		peers := selectPeers(len(nodes), i, sharing)
		for _, p := range peers {
			peer := nodes[p].ID

			transports[i].AddOutput(peer, input[p])
			n.AddPeer(group, peer)

			log.Printf("%x sharing with %x", n.ID[:4], peer[:4])
		}
	}

	for _, n := range nodes {
		n.Start()
	}

	chat(group, nodes[:communicating-1]...)
}

func selectPeers(nodeCount int, currentNode int, sharing int) []int {
	peers := make([]int, 0)

OUTER:
	for len(peers) != sharing {
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

func createNode(transport transport.Transport, id state.PeerID, mode node.Mode) *node.Node {
	ds := store.NewDummyStore()
	return node.NewNode(
		&ds,
		transport,
		state.NewSyncState(),
		Calc,
		0,
		id,
		mode,
		make(chan protobuf.Message),
	)
}

func chat(group state.GroupID, nodes ...*node.Node) {
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

func Calc(count uint64, epoch int64) int64 {
	return epoch + int64(count*2)
}

func peerID() state.PeerID {
	bytes := make([]byte, 65)
	rand.Read(bytes)

	id := state.PeerID{}
	copy(id[:], bytes)

	return id
}

func groupId() state.GroupID {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	id := state.GroupID{}
	copy(id[:], bytes)

	return id
}
