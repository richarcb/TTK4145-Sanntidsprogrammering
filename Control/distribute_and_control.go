package control

import (
	//windows:
	//"time"

	config "../Config"
	"../FSM"
	sync "../Synchronizing"
	"../driver/elevio"
	//"config"
	//"fmt"
	//linux:
	//"../driver/elevio"
)

//
//One Elevator System Test:
/*cancel_illuminate_extern_order_ch<-chan elevio.ButtonEvent,illuminate_extern_order_ch<-chan
 elevio.ButtonEvent, door_timer_ch<-chan int,extern_order_ch<-chan elevio.ButtonEvent, buttons_ch <-chan elevio.ButtonEvent,
floors_ch <-chan int, init_ch <-chan int, /*receiveing channels
reached_extern_floor_ch chan<- elevio.ButtonEvent, new_order_ch chan<- elevio.ButtonEvent, state_ch chan<- Current_state*/
func Distribute_and_control(init_outgoing_msg_ch <-chan sync.Msg_struct,cancel_illuminate_extern_order_ch chan<- int,illuminate_extern_order_ch chan<- elevio.ButtonEvent, reset_received_order_ch <-chan bool, update_outgoing_msg_ch chan<- sync.Msg_struct, update_elev_list <-chan sync.Msg_struct, lost_peers_ch <-chan []int, new_peer_ch <-chan int, new_order_ch <-chan elevio.ButtonEvent, state_ch <-chan FSM.Elevator, extern_order_ch chan<- elevio.ButtonEvent) {

	for {
		select {
		case outgoing_msg = <-init_outgoing_msg_ch:
		case inc_msg := <-update_elev_list:
			if inc_msg.ID != config.ID{
				update_extern_elevator_struct(inc_msg)
				for i:=0;i<2;i++{
					for j:=0; j<config.N_elevators;j++{
						switch inc_msg.Ack_list[i][j]{
							case 0:
								if outgoing_msg.Ack_list[i][j] == -1{
									//Add order to list!
									//Set to zero
									//illuminate button
									bt_type := elevio.BT_HallUp
									if i == 1{bt_type = elevio.BT_HallDown}
									order:=elevio.ButtonEvent{Button: bt_type, Floor: j}
									assignedID:= getLowestCostElevatorID(order)
									add_order_to_elevlist(assignedID, order)
									if assignedID == config.ID{
										go func(){extern_order_ch <- order}()
									}else{
										go func(){illuminate_extern_order_ch<-order}()
									}
									outgoing_msg.Ack_list[i][j] = 0
								}

							case 1:
								if outgoing_msg.Ack_list[i][j] == 0{
									outgoing_msg.Ack_list[i][j] = 1
								}else if outgoing_msg.Ack_list[i][j] == 1{
									bt_type := elevio.BT_HallUp
									if i == 1{bt_type = elevio.BT_HallDown}
									order:=elevio.ButtonEvent{Button: bt_type, Floor: j}
									assignedID := getLowestCostElevatorID(order)
									if assignedID == config.ID{
										outgoing_msg.Ack_list[i][j] = -1
									}
								}
							case -1:
								if outgoing_msg.Ack_list[i][j] == 1{
									outgoing_msg.Ack_list[i][j] = -1
								}else if outgoing_msg.Ack_list[i][j] == -1{
									bt_type := elevio.BT_HallUp
									if i == 1{bt_type = elevio.BT_HallDown}
									order:=elevio.ButtonEvent{Button: bt_type, Floor: j}
									assignedID:= getLowestCostElevatorID(order)
									if assignedID == outgoing_msg.ID{
										outgoing_msg.Ack_list[i][j]=0
										add_order_to_elevlist(assignedID, order)
									}
								}
						}
					}
				}
				go func(){update_outgoing_msg_ch <- outgoing_msg}()
			}

		case lost_peers:= <- lost_peers_ch:
			lost_peers_list := make([]int, len(lost_peers))
			for i := 0; i < len(lost_peers); i++ {
				offline_elevator_list[lost_peers_list[i]] = true
			}
			online_elev := 0
			for i := 0; i < len(offline_elevator_list); i++ {
				if offline_elevator_list[i] == false {
					online_elev++
				}
			}
			if online_elev == 1 {
				single_mode = true
			}

		case new_peer:= <-new_peer_ch:
			offline_elevator_list[new_peer] = false
			online_elev := 0
			for i := 0; i < len(offline_elevator_list); i++ {
				if offline_elevator_list[i] == false {
					online_elev++
				}
			}
			if online_elev > 1 {
				single_mode = false
			}

		case order := <-new_order_ch:
			if single_mode{
				extern_order_ch <- order
			}else{
				set_value_in_ack_list(1,order)
				go func(){update_outgoing_msg_ch <- outgoing_msg}()
			}

		case state := <-state_ch:
			update_local_elevator_struct(state)

		case <-reset_received_order_ch:
			update_outgoing_msg_ch <- outgoing_msg
		}
	}
}
