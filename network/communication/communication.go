package communication

import (
	"time"
	"fmt"
. "../../config"
	"../peers"
)

var outgoing_msg Msg_struct

type Communication_channels struct{
	Init_outgoing_msg_ch chan Msg_struct
	Update_outgoing_msg_ch chan Msg_struct
	Update_control_variables_ch chan Msg_struct
	Lost_peers_ch chan []string
	New_peer_ch chan string
	Outgoing_msg_ch chan Msg_struct
	Incoming_msg_ch chan Msg_struct
	Peer_update_ch chan peers.PeerUpdate
}


func Communication(ch Communication_channels){
	bcastTicker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case outgoing_msg=<-ch.Init_outgoing_msg_ch:

		case p := <-ch.Peer_update_ch:
			fmt.Printf("Peer update:\n")
	        fmt.Printf("  Peers:    %q\n", p.Peers)
	        fmt.Printf("  New:      %q\n", p.New)
	        fmt.Printf("  Lost:     %q\n", p.Lost)
				if len(p.New) > 0{
					go func(){ch.New_peer_ch<-p.New}()
			}
				if len(p.Lost) > 0 {
					//var lost_peers [len(p.Lost)]int
					go func(){ch.Lost_peers_ch<-p.Lost}()
				}

		case incMsg := <-ch.Incoming_msg_ch:
			go func(){ch.Update_control_variables_ch <- incMsg}()

		case outgoing_msg = <-ch.Update_outgoing_msg_ch:

		case <-bcastTicker.C:
			go func(){ch.Outgoing_msg_ch <- outgoing_msg}()

		}
	}
}
