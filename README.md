# Minimal Viable Data Sync

:warning: This code is not production ready, race conditions are likely to occur :warning:

Experimental implementation of the [minimal viable data sync protocol specification](https://notes.status.im/O7Xgij1GS3uREKNtzs7Dyw?view).

# Usage

## Simulation

In order to run a very naive simulation, use the simulation command. The simulation is configurable using various CLI flags.

```
Usage of simulation/simulation.go:
  -communicating int
    	amount of nodes sending messages (default 2)
  -interval int
    	seconds between messages (default 5)
  -nodes int
    	amount of nodes (default 3)
  -offline int
    	percentage of time a node is offline (default 90)
  -sharing int
    	amount of nodes each node shares with (default 2)
```

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
