package main

import (
	"flag"
	"fmt"
	"os"
	"./network/bcast"
	"./network/localip"
	"./network/peers"
. "./config"
	"./control"
	"./fsm"
	"./synchronise"
	"./driver/elevio"
)

func main() {
	port := os.Args[1]
	//config.Init_elevconfig()
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
//
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	init_outgoing_msg_ch:= make(chan Msg_struct)
	init_ID_ch:= make(chan string,1)

	go control.Init_variables(init_ID_ch, init_outgoing_msg_ch)
	init_ID_ch<-id


//15657
	elevio.Init("localhost:"+port, N_floors)
	fsm.Init_mem()
	start_floor := elevio.InitElev()



	/*drv_buttons := make(chan elevio.ButtonEvent)
	  drv_floors  := make(chan int)
	  drv_obstr   := make(chan bool)
	  drv_stop    := make(chan bool)
	*/
	//Network
	//UDPmsgTx := make(chan FSM.Elevator)
	//UDPmsgRx := make(chan FSM.Elevator)

	//FSM
	fsmCh := fsm.FsmChannels{
				Cancel_illuminate_extern_order_ch: make(chan int),
				Illuminate_extern_order_ch: make(chan elevio.ButtonEvent),
				Extern_order_ch: make(chan elevio.ButtonEvent),
				Buttons_ch: make(chan elevio.ButtonEvent),
				Floors_ch: make(chan int),
				Reached_extern_floor_ch: make(chan elevio.ButtonEvent),
				New_order_ch: make(chan elevio.ButtonEvent),
				State_ch: make(chan Elevator),
				}
	/*
	   outgoing_msg_ch := make(chan sync.Msg_struct)
	   incoming_msg_ch := make(chan sync.Msg_struct)
	   peer_update_ch := make(chan peers.PeerUpdate)
	   peer_trans_en_ch := make(chan bool)
	*/
	//Local FSM
	//reset_timer_pl_ch  := make(chan bool)
	//stop_timer_pl_ch  := make(chan bool)
	//power_loss_ch  := make(chan bool)

	//Distribute_and_control channes:
	ctrlCh := control.ControlChannels{
				Reset_received_order_ch : make(chan bool),
				Update_outgoing_msg_ch : make(chan Msg_struct),
				Update_elev_list_ch : make(chan Msg_struct),
				Lost_peers_ch : make(chan []string),
				New_peer_ch : make(chan string),
				Outgoing_msg_ch : make(chan Msg_struct),
				Incoming_msg_ch : make(chan Msg_struct),
				Peer_trans_en_ch : make(chan bool),
				Peer_update_ch : make(chan peers.PeerUpdate),
				Clear_lights_and_extern_orders_ch :make(chan int),
	}

//ROUTINES

	go elevio.PollButtons(fsmCh.Buttons_ch)
	go elevio.PollFloorSensor(fsmCh.Floors_ch)
	//go elevio.PollObstructionSwitch(drv_obstr)
	//go elevio.PollStopButton(drv_stop)

	go control.Distribute_and_control(fsmCh, ctrlCh)
	go fsm.EventHandler(fsmCh, start_floor, ctrlCh.Clear_lights_and_extern_orders_ch)
	go sync.Synchronizing(init_outgoing_msg_ch, fsmCh, ctrlCh)



	//Synchronize
	go bcast.Transmitter(12345, ctrlCh.Outgoing_msg_ch)
	go bcast.Receiver(12345, ctrlCh.Incoming_msg_ch)
	go peers.Transmitter(15647, id, ctrlCh.Peer_trans_en_ch)
	go peers.Receiver(15647, ctrlCh.Peer_update_ch)
	//go sync.TransmitMsg(outgoing_msg_ch)
	//go sync.ReceiveMsg(incoming_msg_ch)


	select {}
}



/*
  for{

    timer1 := time.NewTimer( 4* time.Second)
    <-timer1.C
    RandomOrderGenerator(buttons_ch)
  }
}

/*

func RandomOrderGenerator(receiver chan elevio.ButtonEvent) {
	var order elevio.ButtonEvent
	BT_type := rand.Intn(3)
	if BT_type == 0 {
		order.Button = elevio.BT_HallUp
		order.Floor = rand.Intn(3)
	} else if BT_type == 1 {
		order.Button = elevio.BT_HallDown
		order.Floor = rand.Intn(3) + 1
	} else {
		order.Button = elevio.BT_Cab
		order.Floor = rand.Intn(4)
	}

	receiver <- order
}*/
