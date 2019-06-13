# Minimal Viable Data Sync

![Version](https://img.shields.io/github/tag/status-im/mvds.svg)
[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/status-im/mvds) [![Go Report Card](https://goreportcard.com/badge/github.com/status-im/mvds)](https://goreportcard.com/report/github.com/status-im/mvds)

Experimental implementation of the [minimal viable data sync protocol specification](https://notes.status.im/bZHk_BNkSAe8-TY7DxdNEg?view).

## Usage

### Prerequisites

Ensure you have `protoc` (Protobuf) and Golang installed. Then run `make`.

### Simulation

In order to run a very naive simulation, use the simulation command. The simulation is configurable using various CLI flags.

```
Usage of main.go:
  -communicating int
    	amount of nodes sending messages (default 2)
  -interactive int
    	amount of nodes to use INTERACTIVE mode, the rest will be BATCH (default 3)
  -interval int
    	seconds between messages (default 5)
  -nodes int
    	amount of nodes (default 3)
  -offline int
    	percentage of time a node is offline (default 90)
  -sharing int
    	amount of nodes each node shares with (default 2)
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
