package main

import "github.com/status-im/mvds"

func main() {

	ain := make(chan packet)
	bin := make(chan packet)
	cin := make(chan packet)

	at := Transport{in: ain, out: make(map[mvds.PeerId]chan<- packet)}
	bt := Transport{in: bin, out: make(map[mvds.PeerId]chan<- packet)}
	ct := Transport{in: cin, out: make(map[mvds.PeerId]chan<- packet)}
}

func chat(nodes ...mvds.Node) {

}
