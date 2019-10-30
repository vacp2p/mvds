package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	math "math/rand"
	"time"

	"github.com/vacp2p/mvds/node"
	"github.com/vacp2p/mvds/peers"
	"github.com/vacp2p/mvds/state"
	"github.com/vacp2p/mvds/store"
	"github.com/vacp2p/mvds/transport"
	"go.uber.org/zap"
)

var (
	offline       int
	nodeCount     int
	communicating int
	sharing       int
	interval      int64
	interactive   int
)

func parseFlags() {
	flag.IntVar(&offline, "offline", 90, "percentage of time a node is offline")
	flag.IntVar(&nodeCount, "nodes", 3, "amount of nodes")
	flag.IntVar(&communicating, "communicating", 2, "amount of nodes sending messages")
	flag.IntVar(&sharing, "sharing", 2, "amount of nodes each node shares with")
	flag.Int64Var(&interval, "interval", 5, "seconds between messages")
	flag.IntVar(&interactive, "interactive", 3, "amount of nodes to use InteractiveMode mode, the rest will be BatchMode") // @todo should probably just be how many nodes are interactive
	flag.Parse()
}

func main() {

	parseFlags()
	// @todo validate flags

	transports := make([]*transport.ChannelTransport, 0)
	input := make([]chan transport.Packet, 0)
	nodes := make([]*node.Node, 0)
	for i := 0; i < nodeCount; i++ {
		in := make(chan transport.Packet)

		t := transport.NewChannelTransport(offline, in)

		input = append(input, in)
		transports = append(transports, t)

		mode := node.InteractiveMode
		if i+1 >= interactive {
			mode = node.BatchMode
		}

		node, err := createNode(t, peerID(), mode)
		if err != nil {
			log.Printf("Could not create node: %+v\n", err)
		}
		nodes = append(
			nodes,
			node,
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
			_ = n.AddPeer(group, peer)

			log.Printf("%x sharing with %x", n.ID[:4], peer[:4])
		}
	}

	for _, n := range nodes {
		n.Start(1 * time.Second)
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

func createNode(transport transport.Transport, id state.PeerID, mode node.Mode) (*node.Node, error) {
	ds := store.NewDummyStore()
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return node.NewNode(
		ds,
		transport,
		state.NewSyncState(),
		Calc,
		0,
		id,
		mode,
		peers.NewMemoryPersistence(),
		logger,
	), nil
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
	_, _ = rand.Read(bytes)

	id := state.PeerID{}
	copy(id[:], bytes)

	return id
}

func groupId() state.GroupID {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)

	id := state.GroupID{}
	copy(id[:], bytes)

	return id
}
