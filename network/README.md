#Network Module

This module contains code regarding network and communication between the elevators.

bcast
----
Makes it possible to send custom datatypes over the network. This is done by creating a pair of channels and pairing them up with bcast.Transmitter and bcast.Receiver functions. By transmitting the custom datastructure, it is possible to pick the message up on the specified port. In our project, we chose to only utilize UDP broadcasting.

peers
----
Contains help-functions used to detect if any elevators are on the network, if any new elevators are connected and if any elevators are disconnected.

localip
------
Contains functions used to manage IP addresses such as finding your current IP address.

conn
----
Ensures connection between the elevator driver and our written code

Communication
-------------
Takes care of all network communication between the elevators. It makes sure messages are sent out to the other elevators, and manages incoming ones. The communication module sews together the local fsm and the provided network modules above, making the elevators communicate with eachother.
