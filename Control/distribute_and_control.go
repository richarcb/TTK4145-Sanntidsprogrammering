package control

import (
	//windows:
	//"time"

	config "../Config"
	"../FSM"
	sync "../Synchronizing"
	"../driver/elevio"
	//"config"
	"fmt"
	//linux:
	//"../driver/elevio"
)

//
//One Elevator System Test:
/*cancel_illuminate_extern_order_ch<-chan elevio.ButtonEvent,illuminate_extern_order_ch<-chan
 elevio.ButtonEvent, door_timer_ch<-chan int,extern_order_ch<-chan elevio.ButtonEvent, buttons_ch <-chan elevio.ButtonEvent,
floors_ch <-chan int, init_ch <-chan int, /*receiveing channels
reached_extern_floor_ch chan<- elevio.ButtonEvent, new_order_ch chan<- elevio.ButtonEvent, state_ch chan<- Current_state*/
func Distribute_and_control(clear_lights_and_extern_orders_ch chan<- int, cancel_illuminate_extern_order_ch chan<- int, illuminate_extern_order_ch chan<- elevio.ButtonEvent, reset_received_order_ch <-chan bool, update_outgoing_msg_ch chan<- sync.Msg_struct, update_elev_list <-chan sync.Msg_struct, lost_peers_ch <-chan []string, new_peer_ch <-chan string, new_order_ch <-chan elevio.ButtonEvent, state_ch <-chan FSM.Elevator, extern_order_ch chan<- elevio.ButtonEvent) {

	for {
		select {
//case outgoing_msg = <-init_outgoing_msg_ch:
case inc_msg := <-update_elev_list:
	//fmt.Println((*elev_list[elevID]).State)
	//fmt.Println((*elev_list[inc_msg.ID]).State)

	if inc_msg.ID != elevID {
		if inc_msg.State == FSM.POWERLOSS {
			for i := 0; i < 2; i++ {
				for j := 0; j < config.N_floors; j++ {
					if (*elev_list[inc_msg.ID]).queue[i][j] == 1 && outgoing_msg.Ack_list[i][j] != -1 {
						outgoing_msg.Ack_list[i][j] = 1
						(*elev_list[inc_msg.ID]).queue[i][j] = 0
					}
				}
			}
		}
		if (*elev_list[elevID]).State == FSM.POWERLOSS {
			for i := 0; i < 2; i++ {
				for j := 0; j < config.N_floors; j++ {
					if (*elev_list[elevID]).queue[i][j] == 1 && (inc_msg.Ack_list[i][j] == 1 || inc_msg.Ack_list[i][j] == -1) {
						(*elev_list[elevID]).queue[i][j] = 0
					}
				}
			}
		}
		update_extern_elevator_struct(inc_msg)
		for i := 0; i < 2; i++ {
			for j := 0; j < config.N_floors; j++ {
				switch inc_msg.Ack_list[i][j] {
				case 0:
					if outgoing_msg.Ack_list[i][j] == -1 {
						//Add order to list!
						//Set to zero
						//illuminate button
						//fmt.Println("HALLELUJA")
						bt_type := elevio.BT_HallUp
						if i == 1 {
							bt_type = elevio.BT_HallDown
						}
						order := elevio.ButtonEvent{Button: bt_type, Floor: j}
						assignedID := getLowestCostElevatorID(order)
						//fmt.Println(assignedID)
						add_order_to_elevlist(assignedID, order)
						//fmt.Println(assignedID)
						//fmt.Println(elevID)
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
						bt_type := elevio.BT_HallUp
						if i == 1 {
							bt_type = elevio.BT_HallDown
						}
						order := elevio.ButtonEvent{Button: bt_type, Floor: j}
						assignedID := getLowestCostElevatorID(order)
						//fmt.Println(assignedID)
						//fmt.Println(assignedID)
						/*fmt.Println((*elev_list[elevID]).queue)
						fmt.Println((*elev_list[inc_msg.ID]).queue)*/
						//fmt.Println(elevID)
						//fmt.Println(assignedID)
						if assignedID == elevID {
							outgoing_msg.Ack_list[i][j] = -1
						}
					}
				case -1:
					if outgoing_msg.Ack_list[i][j] == 1 {
						outgoing_msg.Ack_list[i][j] = -1
					} else if outgoing_msg.Ack_list[i][j] == -1 {
						bt_type := elevio.BT_HallUp
						if i == 1 {
							bt_type = elevio.BT_HallDown
						}
						order := elevio.ButtonEvent{Button: bt_type, Floor: j}
						assignedID := getLowestCostElevatorID(order)
						//fmt.Println(assignedID)
						if assignedID == inc_msg.ID {
							outgoing_msg.Ack_list[i][j] = 0
							add_order_to_elevlist(assignedID, order)
							go func() { illuminate_extern_order_ch <- order }()
						}
					}
				}
			}
		}
		go func() { update_outgoing_msg_ch <- outgoing_msg }()
		if inc_msg.State == FSM.DOOROPEN {
			go func() { clear_lights_and_extern_orders_ch <- inc_msg.Last_known_floor }()
			(*elev_list[inc_msg.ID]).queue[0][inc_msg.Last_known_floor] = 0
			(*elev_list[inc_msg.ID]).queue[1][inc_msg.Last_known_floor] = 0
		}

		if elev_list[elevID].State == FSM.DOOROPEN {
			(*elev_list[elevID]).queue[0][elev_list[elevID].Last_known_floor] = 0
			(*elev_list[elevID]).queue[1][elev_list[elevID].Last_known_floor] = 0
		}
	}


		case lost_peers := <-lost_peers_ch:
			for i := 0; i < len(lost_peers); i++ {
				//Take Orders first!
				for j := 0; j < 2; j++ {
					for k := 0; k < config.N_floors; k++ {
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

		case new_peer := <-new_peer_ch:
			add_new_peer_to_elevlist(new_peer)
			if len(elev_list) > 1 {

				single_mode = false
			}

		case order := <-new_order_ch:
			if single_mode && (*elev_list[elevID]).State != FSM.POWERLOSS {
				go func() { extern_order_ch <- order }()
			} else {
				set_value_in_ack_list(1, order)
				go func() { update_outgoing_msg_ch <- outgoing_msg }()
			}
		case state := <-state_ch:
			update_local_elevator_struct(state)
			update_outgoing_msg(state)
			fmt.Println(single_mode)
		case <-reset_received_order_ch:
			update_outgoing_msg_ch <- outgoing_msg
		}
	}
}
