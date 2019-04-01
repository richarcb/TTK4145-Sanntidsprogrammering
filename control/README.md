# Control

This module coordinates information between the elevators on the network, and module the elevator state machine module. It is responsible for making sure that new hall orders is shared on the network, and that the elevator who is best suited based on the [cost function](https://github.com/TTK4145-students-2019/project-team-46/blob/master/control/control_func.go), ends up executing the order. This procedure is explained in the subsections below.  

## Cost function
Calculates the cost of each external order, assigning it to the elevator with the the lowest cost. Making an optimal cost function was not the focus of this project. We implemented a quite simple cost function, which calculates a cost based on the elevators states, destination and current position. Since the elevators have all the required information about everyone on the network, the same cost function runs locally on each elevator.

## Order acknowledging
Each elevator sends a 2 by the number of floors acknowledging list, representing an elevators acknowledge status for HallUp and HallDown orders at each floor. The acknowledge status is 0, 1 or -1, and can only be changed in this order:

* 0->1->(-1)->0.

If an elevators acknowledge status is lacking the status of an incoming message, the status is updated. (For example, local elevator status is 0 and it receives a 1). When an elevator receives a hall order, it sets its acknowledge list element corresponding to this order to 1. The elevator who has the lowest cost sets this element to -1 when it receives a 1. When another elevator receives a message with a -1 from the elevator with the lowest cost, it adds it to the order list corresponding to this elevator, and changes its acknowledge status to 0. The elevator with the lowest cost will then receive a 0, and will then send this order to the elevator state machine module.

By using this acknowledging procedure, we makes sure that at least to elevator knows about each order, such that they can handle the situation if the executing elevator gets a problem.     
