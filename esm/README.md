# Elevator state machine
This is the system finite state machine, which has 4 states;

IDLE - The elevator is at rest, at a floor with closed doors.
MOVING - The elevator is moving.
DOOROPEN - The elevator is at a floor with the door open.
POWERLOSS - The elevator is trying to move, but it is not capable.

The state machine take appropriate actions corresponding to its current state when different events occur. The events which may occur to the state machine is;
Button Event - Someone pushed a button. Handle Cab orders locally, and send hall orders to control module to assign the order to the suitable elevator.
Floor Event - The elevator has reached a new floor
Door Timer Event - The door has been open for set period and is timed out
Powerloss event - The elevator has been moving for a set period and not reached a new floor.
Messages from the control module - Concrete actions which must be done, such as illuminate and darken lights, new hall orders and clear orders.


This module also contains backup features, which saves all its Cab orders in a txt file in case of the program gets terminated and needs to reinitialize.  
