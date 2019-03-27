package sync

import (
	//"strconv"
	"time"
	"fmt"
	//"../FSM"
	peers "../network/peers"
	//config "../Config"
	//control "../Control"
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



func Synchronizing(init_outgoing_msg_ch <-chan Msg_struct, reset_received_order_ch chan<- bool, update_outgoing_msg_ch <-chan Msg_struct, outgoing_msg_ch chan<- Msg_struct, incoming_msg_ch <-chan Msg_struct, update_elev_list_ch chan<- Msg_struct, peer_update_ch <-chan peers.PeerUpdate, peer_trans_en_ch chan<- bool /*button_pressed_ch <- chan elevio.ButtonEvent*/, lost_peers_ch chan<- []string, new_peer_ch chan<- string) {
	bcastTicker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case outgoing_msg=<-init_outgoing_msg_ch:
		case p := <-peer_update_ch:
			fmt.Printf("Peer update:\n")
	        fmt.Printf("  Peers:    %q\n", p.Peers)
	        fmt.Printf("  New:      %q\n", p.New)
	        fmt.Printf("  Lost:     %q\n", p.Lost)
				if len(p.New) > 0{
					go func(){new_peer_ch<-p.New}()
			}
				if len(p.Lost) > 0 {
					//var lost_peers [len(p.Lost)]int
					go func(){lost_peers_ch<-p.Lost}()
				}


		case incMsg := <-incoming_msg_ch:
			go func(){update_elev_list_ch <- incMsg}()


		case outgoing_msg = <-update_outgoing_msg_ch:

		case <-bcastTicker.C:
			go func(){outgoing_msg_ch <- outgoing_msg}()

			//fmt.Printf("STATE SENT \n")

		}
	}
}
