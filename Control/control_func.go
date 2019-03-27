package control

import (
	config "../Config"
	"../FSM"
	sync "../Synchronizing"
	"../driver/elevio"
	//"fmt"
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
var elevID string
//var offline_elev_list [config.N_elevators]bool
var single_mode bool

type elevator_states struct {
	//Destination floor
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            FSM.ElevState
	queue           [2][config.N_floors]int
}

type elevator_list map[string]*elevator_states

var elev_list elevator_list

var outgoing_msg sync.Msg_struct

func Init_variables(init_ID_ch <-chan string, init_outgoing_msg_ch chan<- sync.Msg_struct) {
	select{
	case ID_string := <- init_ID_ch:
		elevID = ID_string
		single_mode = true
			var empty_queue [2][config.N_floors]int
			for j := 0; j < 2; j++ {
				for k := 0; k < config.N_floors; k++ {
					empty_queue[j][k] = 0
				}
			}
			var empty_ack_list [2][config.N_floors]int
			for j := 0; j < 2; j++ {
				for k := 0; k < config.N_floors; k++ {
					empty_ack_list[j][k] = 0
				}
			}
			elev_list = make(map[string]*elevator_states)
			elev_list[ID_string] = &elevator_states{Destination: FSM.Empty_order, Last_known_floor: -1, Dir: elevio.MD_Stop, State: FSM.IDLE,  queue: empty_queue}
			outgoing_msg = sync.Msg_struct{Destination: FSM.Empty_order, Last_known_floor: -1, Dir: elevio.MD_Stop, State: FSM.IDLE, Ack_list:empty_ack_list ,ID: ID_string}
			go func(){init_outgoing_msg_ch<-outgoing_msg}()
	}
	/*
	//initialize out_of_ queue
	for k := 0; k < config.N_elevators; k++ {
		if k == config.ID {
			offline_elev_list[k] = false
		}else {
			offline_elev_list[k] = true
		}
	}*/
}
func add_new_peer_to_elevlist(id string){
	var empty_queue [2][config.N_floors]int
	for j := 0; j < 2; j++ {
		for k := 0; k < config.N_floors; k++ {
			empty_queue[j][k] = 0
		}
	}
	new_empty_peer :=elevator_states{Destination: FSM.Empty_order, Last_known_floor: -1, Dir: elevio.MD_Stop, State: FSM.IDLE,  queue: empty_queue}
	elev_list[id] = &new_empty_peer

}


func set_value_in_ack_list(value int, order elevio.ButtonEvent){
	if order.Floor == FSM.Empty_order.Floor{return}
	bt_type:=0
	if order.Button == elevio.BT_HallDown{bt_type = 1}
	outgoing_msg.Ack_list[bt_type][order.Floor] = value

}

func update_local_elevator_struct(elevator FSM.Elevator) {
	//Updates its own elevator_struct

	(*elev_list[elevID]).Destination = elevator.Destination
	(*elev_list[elevID]).Last_known_floor = elevator.Last_known_floor
	(*elev_list[elevID]).State = elevator.State
	(*elev_list[elevID]).Dir = elevator.Dir
}

func update_outgoing_msg (elevator FSM.Elevator){
	outgoing_msg.Destination = elevator.Destination
	outgoing_msg.Last_known_floor = elevator.Last_known_floor
	outgoing_msg.State = elevator.State
	outgoing_msg.Dir = elevator.Dir
}

func update_extern_elevator_struct(elevator sync.Msg_struct) {
	//Update elevator_struct from msg!
	if elev_list[elevator.ID] == nil{return}
	(*elev_list[elevator.ID]).Destination = elevator.Destination
	(*elev_list[elevator.ID]).Last_known_floor = elevator.Last_known_floor
	(*elev_list[elevator.ID]).State = elevator.State
	(*elev_list[elevator.ID]).Dir = elevator.Dir
}



func cost_function(id string, order elevio.ButtonEvent) int {
	cost := 0
	if (elev_list[id].State == FSM.IDLE || elev_list[id].State == FSM.DOOROPEN) && elev_list[id].Last_known_floor == order.Floor{cost-=15}
	//Order already in list...
	if (*elev_list[id]).Destination.Floor == order.Floor{
		cost -= 15
	}
	for i:= 0; i<2 ; i++{
		if (*elev_list[id]).queue[i][order.Floor] == 1{
			cost -= 15
		}
	}
	
	if order.Button == elevio.BT_HallUp { //Order is up
		if elev_list[id].Last_known_floor < order.Floor && elev_list[id].Destination.Floor > order.Floor { //going up and flor is between:
			cost -= 10
		}
	} else { //Order is down                                                                                                              //config.N_floors-2 since the first down button is in the second floor!
		if elev_list[id].Last_known_floor > order.Floor && elev_list[id].Destination.Floor < order.Floor && elev_list[id].State == FSM.MOVING { //going down and floor is between orders:Needs to check MOVING since no destination 0-1<order.Floor
			cost -= 10
		}
	}
	if elev_list[id].State == FSM.IDLE && elev_list[id].Destination.Floor == FSM.Empty_order.Floor { //Nothing to do
		cost -= 5
	}
	//Adding the value of the distance
	if order.Floor > elev_list[id].Last_known_floor {
		cost += order.Floor - elev_list[id].Last_known_floor
	} else {
		cost += elev_list[id].Last_known_floor - order.Floor
	}

	return cost
}

func add_order_to_elevlist (assigned_id string, order elevio.ButtonEvent){
	bt_type := 0
	if order.Button == elevio.BT_HallDown{bt_type = 1}
	elev_list[assigned_id].queue[bt_type][order.Floor] = 1
}

func getLowestCostElevatorID(order elevio.ButtonEvent) string {
	lowestCost := config.N_floors
	assignedID := elevID
	//Get Number_of_Online_elevators! (POWERLOSS???)

	for k := range elev_list{
		//fmt.Println(k)
		cost := cost_function(k, order)
		/*
		fmt.Println("COST:")
		fmt.Println(cost)
		fmt.Println("\n")
		*/
		if cost < lowestCost {
			lowestCost = cost
			assignedID = k
		}
	}
	return assignedID
}
