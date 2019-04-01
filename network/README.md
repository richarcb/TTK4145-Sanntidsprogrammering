# Network Module (UDP broadcast)

This module contains code for network and communication between the elevators. Except for the [communication](actual URL to navigate) folder, the code in this module is written by [TTK4145](https://github.com/TTK4145) and can be found [here](https://github.com/TTK4145/Network-go). There has been made some minor changes to the code to make it suit our system,  but it is mostly the equal to the original module.

## Features

Channel-in/channel-out pairs of (almost) any custom or built-in datatype can be supplied to a pair of transmitter/receiver functions. Data sent to the transmitter function is automatically serialized and broadcasted on the specified port. Any messages received on the receiver's port are deserialized (as long as they match any of the receiver's supplied channel datatypes) and sent on the corresponding channel. See [bcast.Transmitter and bcast.Receiver](https://github.com/TTK4145/Network-go/blob/master/network/bcast/bcast.go).

Peers on the local network can be detected by supplying your own ID to a transmitter and receiving peer updates (new, current and lost peers) from the receiver. See [peers.Transmitter and peers.Receiver](https://github.com/TTK4145/Network-go/blob/master/network/peers/peers.go).

Finding your own local IP address can be done with the [LocalIP](https://github.com/TTK4145/Network-go/blob/master/network/localip/localip.go) convenience function, but only when you are connected to the internet.

## Communication
  Sends a message to other elevators at a set frequency, as well as detecting new and lost peers.
