# Elevator Project TTK4145

To run this program on the elevators or using a simulator, follow the description found [here](https://github.com/TTK4145/elevator-server)

Terminal commands using port :
* go run main.go PORTNMUBER
* ./SimElevatorServer --port PORTNUMBER


## Problem Description

In this project, we have created a program for controlling x elevators with y floors. The elevators communicates over network. A full problem description can be found [here](https://github.com/TTK4145/Project)

### Main requirements
The following requirements should be met:

* Multiple elevators are more efficient than one
* The lights and buttons should function as expected
* No orders are lost
* An individual elevator should behave sensibly and efficiently

## Our solution

We wrote our solution in the google go programming language. We decided to share information between elevators by using UDP boradcast, sending information to all peers connected at a set frequency. A more complementary description of each module is given in the README.md files for the module.
