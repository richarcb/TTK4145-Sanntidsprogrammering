package control

import (
	."../config"
)

func Control(update_outgoing_msg_ch chan<- Msg_struct, new_peer_ch <-chan string, lost_peers_ch <-chan []string, update_control_variables_ch <-chan Msg_struct, extern_order_ch chan<- Order, new_order_ch <-chan Order, state_ch<-chan Elevator, illuminate_extern_order_ch chan<- Order, clear_lights_and_extern_orders_ch chan<- int) {
	for {
		select {
case inc_msg := <-update_control_variables_ch:
	if inc_msg.ID != elevID {
		handling_powerloss(inc_msg)
		synchronize(inc_msg, illuminate_extern_order_ch, extern_order_ch)
		update_extern_elevator_struct(inc_msg)
		go func() { update_outgoing_msg_ch <- outgoing_msg}()
		delete_order_if_handled(inc_msg.ID, clear_lights_and_extern_orders_ch)
	}

case lost_peers := <-lost_peers_ch:
			lost_peer_event(lost_peers)

case new_peer := <-new_peer_ch:
			add_new_peer_to_elevlist(new_peer)

case order := <-new_order_ch:
			if single_mode && (*elev_list[elevID]).state != POWERLOSS {
				go func() { extern_order_ch <- order }()
			} else {
				set_value_in_ack_list(1, order)
				go func() { update_outgoing_msg_ch <- outgoing_msg }()
			}
case state := <-state_ch:
			update_local_elevator_struct(state)
			update_outgoing_msg(state)
			delete_order_if_handled(elevID, clear_lights_and_extern_orders_ch)
			update_outgoing_msg_ch <- outgoing_msg
		}
	}
}
