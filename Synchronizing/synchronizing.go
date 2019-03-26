package sync

import (
	"strconv"
	"time"
	"fmt"
	//"../FSM"
	peers "../network/peers"
	//config "../Config"
)

/*
button_pressed_ch := make(chan, elevio.ButtonEvent)
outgoing_msg_ch := make(chan sync.Msg_struct)
incoming_msg_ch := make(chan sync.Msg_struct)
peer_update_ch := make(chan peers.PeerUpdate)
peer_trans_en_ch := make(chan bool)
lost_peer_ch:= make(chan string[])
new_peer_ch:= make(chan string[])


var message Msg_struct
*/



func Synchronizing(reset_received_order_ch chan<- bool, update_outgoing_msg_ch <-chan Msg_struct, outgoing_msg_ch chan<- Msg_struct, incoming_msg_ch <-chan Msg_struct, update_elev_list chan<- Msg_struct, peer_update_ch <-chan peers.PeerUpdate, peer_trans_en_ch chan<- bool /*button_pressed_ch <- chan elevio.ButtonEvent*/, lost_peers_ch chan<- []int, new_peer_ch chan<- int) {
	bcastTicker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case p := <-peer_update_ch:
			fmt.Printf("Peer update:\n")
	        fmt.Printf("  Peers:    %q\n", p.Peers)
	        fmt.Printf("  New:      %q\n", p.New)
	        fmt.Printf("  Lost:     %q\n", p.Lost)
				if len(p.New) > 0 {
					new_id_int, _ := strconv.Atoi(p.New)
					go func(){new_peer_ch<-new_id_int}()
			}
				if len(p.Lost) > 0 {
					//var lost_peers [len(p.Lost)]int
					lost_peers := make([]int, len(p.Lost))
					for i := 0; i < len(p.Lost); i++ {
						lost_peers[i], _ = strconv.Atoi(p.Lost[i])
					}
					go func(){lost_peers_ch<-lost_peers}()
				}


		case incMsg := <-incoming_msg_ch:

			update_elev_list <- incMsg


		case outgoing_msg = <-update_outgoing_msg_ch:

		case <-bcastTicker.C:
			outgoing_msg_ch <- outgoing_msg
			fmt.Printf("STATE SENT \n")

		}
	}
}
