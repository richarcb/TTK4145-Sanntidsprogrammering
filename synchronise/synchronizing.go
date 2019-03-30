package sync

import (
	"time"
	"fmt"
	"../fsm"
	"../control"
 //	"../network/peers"
. "../config"
)

var outgoing_msg Msg_struct

func Synchronizing(init_outgoing_msg_ch chan Msg_struct, ch1 fsm.FsmChannels, ch2 control.ControlChannels){
//	go sync.Synchronizing(init_outgoing_msg_ch, reset_received_order_ch, update_outgoing_msg_ch, outgoing_msg_ch, incoming_msg_ch, update_elev_list, peer_update_ch, peer_trans_en_ch, lost_peers_ch, new_peer_ch)
	bcastTicker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case outgoing_msg=<-init_outgoing_msg_ch:

		case p := <-ch2.Peer_update_ch:
			fmt.Printf("Peer update:\n")
	        fmt.Printf("  Peers:    %q\n", p.Peers)
	        fmt.Printf("  New:      %q\n", p.New)
	        fmt.Printf("  Lost:     %q\n", p.Lost)
				if len(p.New) > 0{
					go func(){ch2.New_peer_ch<-p.New}()
			}
				if len(p.Lost) > 0 {
					//var lost_peers [len(p.Lost)]int
					go func(){ch2.Lost_peers_ch<-p.Lost}()
				}

		case incMsg := <-ch2.Incoming_msg_ch:
			go func(){ch2.Update_elev_list_ch <- incMsg}()

		case outgoing_msg = <-ch2.Update_outgoing_msg_ch:

		case <-bcastTicker.C:
			go func(){ch2.Outgoing_msg_ch <- outgoing_msg}()

		}
	}
}
