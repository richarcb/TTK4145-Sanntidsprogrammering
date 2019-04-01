package control

import (
. "../config"
	"sort"
)

func add_new_peer_to_elevlist(id string) {
	var empty_queue [2][N_floors]int
	var empty_ack_list [2][N_floors]int
	for j := 0; j < 2; j++ {
		for k := 0; k < N_floors; k++ {
			empty_queue[j][k] = 0
			empty_ack_list[j][k] = 0
		}
	}

	new_empty_peer := elevator_states{destination: empty_order, last_known_floor: -1, dir: MD_Stop, state: IDLE, queue: empty_queue, ack_list: empty_ack_list}
	elev_list[id] = &new_empty_peer
	if len(elev_list) > 1 {
		single_mode = false
	}

}

func set_value_in_ack_list(value int, order Order) {
	if order.Floor == empty_order.Floor {
		return
	}
	bt_type := 0
	if order.Button == BT_HallDown {
		bt_type = 1
	}
	outgoing_msg.Ack_list[bt_type][order.Floor] = value

}

func update_local_elevator_struct(elevator Elevator) {
	//Updates its own elevator_struct in elevator_list
	(*elev_list[elevID]).destination = elevator.Destination
	(*elev_list[elevID]).last_known_floor = elevator.Last_known_floor
	(*elev_list[elevID]).state = elevator.State
	(*elev_list[elevID]).dir = elevator.Dir
	(*elev_list[elevID]).ack_list = outgoing_msg.Ack_list
}

func update_outgoing_msg(elevator Elevator) {
	outgoing_msg.Destination = elevator.Destination
	outgoing_msg.Last_known_floor = elevator.Last_known_floor
	outgoing_msg.State = elevator.State
	outgoing_msg.Dir = elevator.Dir
}

func update_extern_elevator_struct(elevator Msg_struct) {
	//Update elevator_list from msg
	if elev_list[elevator.ID] == nil {
		return
	}
	(*elev_list[elevator.ID]).destination = elevator.Destination
	(*elev_list[elevator.ID]).last_known_floor = elevator.Last_known_floor
	(*elev_list[elevator.ID]).state = elevator.State
	(*elev_list[elevator.ID]).dir = elevator.Dir
	(*elev_list[elevator.ID]).ack_list = elevator.Ack_list
}
func cost_function(id string, order Order) int {
	cost := 0

	//Make sure no elev with powerloss gets assigned
	if (*elev_list[id]).state == POWERLOSS {
		cost += 100
	}

	if ((*elev_list[id]).state == IDLE || (*elev_list[id]).state == DOOROPEN) && (*elev_list[id]).last_known_floor == order.Floor {
		cost -= 15
	}
	//Order already in list
	if (*elev_list[id]).destination.Floor == order.Floor {
		cost -= 15
	}
	for i := 0; i < 2; i++ {
		if (*elev_list[id]).queue[i][order.Floor] == 1 {
			cost -= 15
		}
	}
	if order.Button == BT_HallUp { //Order is up
		if (*elev_list[id]).last_known_floor < order.Floor && (*elev_list[id]).destination.Floor > order.Floor { //going up and flor is between:
			cost -= 10
		}
	} else { //Order is down
		if (*elev_list[id]).last_known_floor > order.Floor && (*elev_list[id]).destination.Floor < order.Floor && (*elev_list[id]).destination.Floor != empty_order.Floor{ //going down and floor is between orders
			cost -= 10
		}
	}
	if (*elev_list[id]).state == IDLE && (*elev_list[id]).destination.Floor == empty_order.Floor { //Nothing to do
		cost -= 5
	}
	//Adding the value of the distance
	if order.Floor > (*elev_list[id]).last_known_floor {
		cost += order.Floor - (*elev_list[id]).last_known_floor
	} else {
		cost += (*elev_list[id]).last_known_floor - order.Floor
	}
	return cost
}

func add_order_to_elevlist(assigned_id string, order Order) {
	bt_type := 0
	if order.Button == BT_HallDown {
		bt_type = 1
	}
	(*elev_list[assigned_id]).queue[bt_type][order.Floor] = 1
}

func getLowestCostElevatorID(order Order) string {
	lowestCost := N_floors
	assignedID := ""

	//Sorting the IDs to make sure that the cost function give the same result when runned on different elevators
	var keys []string
	for k:= range elev_list{
		keys = append(keys,k)
	}
	sort.Strings(keys)

	for i:=0; i<len(keys);i++ {
		cost := cost_function(keys[i], order)

		if cost < lowestCost {
			lowestCost = cost
			assignedID = keys[i]
		}
	}
	return assignedID
}

//The order acknowledging system is implemented in this function
func synchronize (inc_msg Msg_struct, illuminate_extern_order_ch chan<- Order, extern_order_ch chan<- Order){
	for i := 0; i < 2; i++ {
		for j := 0; j < N_floors; j++ {
			switch inc_msg.Ack_list[i][j] {
			case 0:
				if outgoing_msg.Ack_list[i][j] == -1 && (*elev_list[inc_msg.ID]).ack_list[i][j] == -1{
					bt_type := BT_HallUp
					if i == 1 {
						bt_type = BT_HallDown
					}
					order := Order{Button: bt_type, Floor: j}
					assignedID := getLowestCostElevatorID(order)
					add_order_to_elevlist(assignedID, order)
					if assignedID == elevID {
						go func() { extern_order_ch <- order }()
					} else {
						go func() { illuminate_extern_order_ch <- order }()
					}
					outgoing_msg.Ack_list[i][j] = 0
				}

			case 1:
				if outgoing_msg.Ack_list[i][j] == 0 {
					outgoing_msg.Ack_list[i][j] = 1
				} else if outgoing_msg.Ack_list[i][j] == 1 {
					bt_type := BT_HallUp
					if i == 1 {
						bt_type = BT_HallDown
					}
					order := Order{Button: bt_type, Floor: j}
					assignedID := getLowestCostElevatorID(order)
					if assignedID == elevID {
						outgoing_msg.Ack_list[i][j] = -1
					}
				}
			case -1:
				if outgoing_msg.Ack_list[i][j] == 1 {
					outgoing_msg.Ack_list[i][j] = -1
				} else if outgoing_msg.Ack_list[i][j] == -1 {
					bt_type := BT_HallUp
					if i == 1 {
						bt_type = BT_HallDown
					}
					order := Order{Button: bt_type, Floor: j}
					assignedID := getLowestCostElevatorID(order)
					if assignedID == inc_msg.ID {
						outgoing_msg.Ack_list[i][j] = 0
						add_order_to_elevlist(assignedID, order)
						go func() { illuminate_extern_order_ch <- order }()
					}
				}
			}
		}
	}
}

func handling_powerloss(inc_msg Msg_struct){
	if inc_msg.State == POWERLOSS {
		for i := 0; i < 2; i++ {
			for j := 0; j < N_floors; j++ {
				//Delete order in elevators queue if elevator has powerloss
				if (*elev_list[inc_msg.ID]).queue[i][j] == 1 && outgoing_msg.Ack_list[i][j] != -1 {
					outgoing_msg.Ack_list[i][j] = 1
					(*elev_list[inc_msg.ID]).queue[i][j] = 0
				}
			}
		}
	}
	if (*elev_list[elevID]).state == POWERLOSS {
		for i := 0; i < 2; i++ {
			for j := 0; j < N_floors; j++ {
				//Delete orders in own queue if poerloss
				if (*elev_list[elevID]).queue[i][j] == 1 && (inc_msg.Ack_list[i][j] == 1 || inc_msg.Ack_list[i][j] == -1) {
					(*elev_list[elevID]).queue[i][j] = 0
				}
			}
		}
	}
}

func delete_order_if_handled(id string, clear_lights_and_extern_orders_ch chan<- int) {
	if (*elev_list[id]).state == DOOROPEN {
		if id != elevID{
			go func() { clear_lights_and_extern_orders_ch <- (*elev_list[id]).last_known_floor}()
		}
		(*elev_list[id]).queue[0][(*elev_list[id]).last_known_floor] = 0
		(*elev_list[id]).queue[1][(*elev_list[id]).last_known_floor] = 0
	}
}

func lost_peer_event (lost_peers []string){
	for i := 0; i < len(lost_peers); i++ {
		//Copy orders before deleting elevator
		for j := 0; j < 2; j++ {
			for k := 0; k < N_floors; k++ {
				if (*elev_list[lost_peers[i]]).queue[j][k] == 1 && outgoing_msg.Ack_list[j][k] != -1 {
					outgoing_msg.Ack_list[j][k] = 1
					(*elev_list[lost_peers[i]]).queue[j][k] = 0
				}
			}
		}
		if lost_peers[i] != elevID{
			delete(elev_list, lost_peers[i])
		}

	}
	if len(elev_list) == 1 {
		single_mode = true
	}

}
