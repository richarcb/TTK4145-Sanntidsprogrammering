package control

import (
	config "../Config"
	"../FSM"
	sync "../Synchronizing"
	"../driver/elevio"
)

//var outgoing_msg sync.Msg_struct

/*
type Elevator struct {
	//Destination floor
	destination      elevio.ButtonEvent
	Last_known_floor int
	dir              elevio.MotorDirection
	state            ElevState
	ID 							 int
}
*/
var offline_elevator_list [config.N_elevators]bool
var single_mode bool

type elevator_states struct {
	//Destination floor
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            FSM.ElevState
	ID				 			 int
	queue           [2][config.N_floors]int
}

var elevator_list [config.N_elevators]elevator_states
var outgoing_msg sync.Msg_struct

func Init_variables() {
	single_mode = true



	for i := 0; i < config.N_elevators; i++ {
		var empty_queue [2][config.N_floors]int // Not shallow copy!
		for j := 0; j < 2; j++ {
			for k := 0; k < config.N_floors; k++ {
				empty_queue[j][k] = 0
			}
		}
		elevator_list[i] = elevator_states{Destination: FSM.Empty_order, Last_known_floor: -1, Dir: elevio.MD_Stop, State: FSM.IDLE,  queue: empty_queue, ID: config.ID}
	}


	//initialize out_of_ queue
	for k := 0; k < config.N_elevators; k++ {
		if k == config.ID {
			offline_elevator_list[k] = false
		} else {
			offline_elevator_list[k] = true
		}
	}
}


func set_value_in_ack_list(value int, order elevio.ButtonEvent){
	if order.Floor == FSM.Empty_order.Floor{return}
	bt_type:=0
	if order.Button == elevio.BT_HallDown{bt_type = 1}
	outgoing_msg.Ack_list[bt_type][order.Floor] = value

}

func update_local_elevator_struct(elevator FSM.Elevator) {
	//Updates its own elevator_struct
	elevator_list[config.ID].Destination = elevator.Destination
	elevator_list[config.ID].Last_known_floor = elevator.Last_known_floor
	elevator_list[config.ID].State = elevator.State
	elevator_list[config.ID].Dir = elevator.Dir
}

func update_extern_elevator_struct(elevator sync.Msg_struct) {
	//Update elevator_struct from msg!
	elevator_list[elevator.ID].Destination = elevator.Destination
	elevator_list[elevator.ID].Last_known_floor = elevator.Last_known_floor
	elevator_list[elevator.ID].State = elevator.State
	elevator_list[elevator.ID].Dir = elevator.Dir
}



func cost_function(id int, order elevio.ButtonEvent) int {
	cost := 0
	if order.Button == elevio.BT_HallUp { //Order is up
		if elevator_list[id].Last_known_floor < order.Floor && elevator_list[id].Destination.Floor > order.Floor { //going up and flor is between:
			cost -= 10
		}
	} else { //Order is down
		order.Floor -= config.N_floors - 2                                                                                                                  //config.N_floors-2 since the first down button is in the second floor!
		if elevator_list[id].Last_known_floor > order.Floor && elevator_list[id].Destination.Floor < order.Floor && elevator_list[id].State == FSM.MOVING { //going down and floor is between orders:Needs to check MOVING since no destination 0-1<order.Floor
			cost -= 10
		}
	}
	if elevator_list[id].State == FSM.IDLE && elevator_list[id].Destination.Floor == FSM.Empty_order.Floor { //Nothing to do
		cost -= 5
	}
	//Adding the value of the distance
	if order.Floor > elevator_list[id].Last_known_floor {
		cost += order.Floor - elevator_list[id].Last_known_floor
	} else {
		cost += elevator_list[id].Last_known_floor - order.Floor
	}
	return cost
}

func add_order_to_elevlist (id int, order elevio.ButtonEvent){
	bt_type := 0
	if order.Button == elevio.BT_HallDown{bt_type = 1}
	elevator_list[id].queue[bt_type][order.Floor] = 1
}

func getLowestCostElevatorID(order elevio.ButtonEvent) int {
	lowestCost := config.N_floors
	assignedID := -1
	//Get Number_of_Online_elevators! (POWERLOSS???)

	for i := 0; i < config.N_elevators; i++ {
		if offline_elevator_list[i] {
			continue
		}
		cost := cost_function(i, order)
		if cost < lowestCost {
			lowestCost = cost
			assignedID = i
		}
	}
	return assignedID
}
