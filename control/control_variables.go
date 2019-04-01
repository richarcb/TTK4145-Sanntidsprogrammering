package control

import (
. "../config"
)

var elevID string
var single_mode bool
type elevator_list map[string]*elevator_states
var elev_list elevator_list
var outgoing_msg Msg_struct
var empty_order Order
type elevator_states struct {
	destination      		Order
	last_known_floor 		int
	dir              		MotorDirection
	state            		ElevState
	queue            		[2][N_floors]int
	ack_list						[2][N_floors]int
}

func Init_variables(init_ID_ch <-chan string, init_outgoing_msg_ch chan<- Msg_struct) {
	select {
	case ID_string := <-init_ID_ch:
		elevID = ID_string
		single_mode = true
		var empty_queue [2][N_floors]int
		for j := 0; j < 2; j++ {
			for k := 0; k < N_floors; k++ {
				empty_queue[j][k] = 0
			}
		}
		var empty_ack_list [2][N_floors]int
		var empty_elev_state_ack_list [2][N_floors]int
		for j := 0; j < 2; j++ {
			for k := 0; k < N_floors; k++ {
				empty_elev_state_ack_list[j][k] = 0
				empty_ack_list[j][k] = 0
			}
		}

		empty_order.Floor = -1
		empty_order.Button = BT_No_call
		elev_list = make(map[string]*elevator_states)
		elev_list[ID_string] = &elevator_states{destination: empty_order, last_known_floor: -1, dir: MD_Stop, state: IDLE, queue: empty_queue, ack_list:empty_elev_state_ack_list}
		outgoing_msg = Msg_struct{Destination: empty_order, Last_known_floor: -1, Dir: MD_Stop, State: IDLE, Ack_list: empty_ack_list, ID: ID_string}
		go func() { init_outgoing_msg_ch <- outgoing_msg }()
	}
}
