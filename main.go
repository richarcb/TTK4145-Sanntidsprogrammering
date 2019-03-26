package main

import (
	//"control"
	//linux
	/*
	   "./driver/elevio"
	   "./FSM"
	   "./Control"
	*/
	//Windows
	//"elevator_project/network/peers"
	"flag"
	"fmt"
	"os"

	"./network/bcast"
	"./network/localip"
	"./network/peers"

	config "./Config"
	control "./Control"
	"./FSM"
	//"os"
	sync "./Synchronizing"
	"./driver/elevio"
	//"math/rand"
	//"sync"
	// "Network-go-master/network/bcast"
	//"Network-go-master/network/peers"
	//"fmt"
	//"time"
	//"strconv"
)

func main() {
	//port := os.Args[1]

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
//15657
	elevio.Init("localhost:15657", config.N_floors)
	FSM.Init_mem()
	control.Init_variables()
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
	cancel_illuminate_extern_order_ch := make(chan int)
	illuminate_extern_order_ch := make(chan elevio.ButtonEvent)
	extern_order_ch := make(chan elevio.ButtonEvent)
	buttons_ch := make(chan elevio.ButtonEvent)
	floors_ch := make(chan int)
	reached_extern_floor_ch := make(chan elevio.ButtonEvent)
	new_order_ch := make(chan elevio.ButtonEvent)
	state_ch := make(chan FSM.Elevator)

	/*
	   outgoing_msg_ch := make(chan sync.Msg_struct)
	   incoming_msg_ch := make(chan sync.Msg_struct)
	   peer_update_ch := make(chan peers.PeerUpdate)
	   peer_trans_en_ch := make(chan bool)
	*/
	//Local FSM:
	door_timer_ch := make(chan bool)
	reset_timer_ch := make(chan bool)
	//reset_timer_pl_ch  := make(chan bool)
	//stop_timer_pl_ch  := make(chan bool)
	//power_loss_ch  := make(chan bool)

	//Distribute_and_control channes:
	reset_received_order_ch := make(chan bool)
	update_outgoing_msg_ch := make(chan sync.Msg_struct)
	update_elev_list := make(chan sync.Msg_struct)
	lost_peers_ch := make(chan []int)
	new_peer_ch := make(chan int)
	outgoing_msg_ch := make(chan sync.Msg_struct)
	incoming_msg_ch := make(chan sync.Msg_struct)
	peer_trans_en_ch := make(chan bool)
	peer_update_ch := make(chan peers.PeerUpdate)
	init_outgoing_msg_ch:=make(chan sync.Msg_struct)

	sync.Init_sync(init_outgoing_msg_ch)

	go elevio.PollButtons(buttons_ch)
	go elevio.PollFloorSensor(floors_ch)
	//go elevio.PollObstructionSwitch(drv_obstr)
	//go elevio.PollStopButton(drv_stop)
	go FSM.DoorTimer(door_timer_ch, reset_timer_ch)
	go control.Distribute_and_control(init_outgoing_msg_ch,cancel_illuminate_extern_order_ch, illuminate_extern_order_ch,reset_received_order_ch, update_outgoing_msg_ch, update_elev_list, lost_peers_ch, new_peer_ch, new_order_ch, state_ch, extern_order_ch)
	go FSM.Fsm(start_floor, cancel_illuminate_extern_order_ch, illuminate_extern_order_ch, door_timer_ch, extern_order_ch, buttons_ch, floors_ch, reached_extern_floor_ch, new_order_ch, state_ch, reset_timer_ch)
	go sync.Synchronizing(reset_received_order_ch, update_outgoing_msg_ch, outgoing_msg_ch, incoming_msg_ch, update_elev_list, peer_update_ch, peer_trans_en_ch, lost_peers_ch, new_peer_ch)



	//Synchronize
	go bcast.Transmitter(12345, outgoing_msg_ch)
	go bcast.Receiver(12345, incoming_msg_ch)
	go peers.Transmitter(15647, id, peer_trans_en_ch)
	go peers.Receiver(15647, peer_update_ch)
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
