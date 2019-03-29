# Distribute & Control

This module coordinates information between the Synchronize module and FSM. It generates a new order message for the network module to broadcast and forwards matching order messages to the FSM. 

## Cost calculator
Based on minimal movement. Calculates the cost of each external order, assigning it to the elevator with the the lowest cost.

## Memory

 Keeps track of peers (elevators) in the system (lost peer, new peer...) and their status. 
