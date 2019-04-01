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
	"./esm"
	com "./network/communication"
	"./driver/elevio"
)

func initialize_elevator_system() (int,string){
	port := os.Args[1]
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

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

	elevio.Init("localhost:"+port, N_floors)
	esm.Init_mem()
 	return elevio.InitElev(),id
}


func main() {
	start_floor,id := initialize_elevator_system()

	esm_ch := esm.Esm_channels{
		Clear_lights_and_extern_orders_ch: 	make(chan int),
		Illuminate_extern_order_ch:					make(chan Order),
		Extern_order_ch: 										make(chan Order),
		Buttons_ch: 												make(chan Order),
		Floors_ch: 													make(chan int),
		New_order_ch: 											make(chan Order),
		State_ch: 													make(chan Elevator)}

	//Network communication channels
	com_ch := com.Communication_channels{
		Init_outgoing_msg_ch:  							make(chan Msg_struct),
		Update_outgoing_msg_ch:							make(chan Msg_struct),
		Update_control_variables_ch:				make(chan Msg_struct),
		Lost_peers_ch:											make(chan []string),
		New_peer_ch: 												make(chan string),
		Outgoing_msg_ch:										make(chan Msg_struct),
		Incoming_msg_ch:  									make(chan Msg_struct),
		Peer_update_ch: 										make(chan peers.PeerUpdate)}

		peer_trans_en_ch:= make(chan bool)


		//GO ROUTINES
//Driver sensors
	go elevio.PollButtons(esm_ch.Buttons_ch)
	go elevio.PollFloorSensor(esm_ch.Floors_ch)

//Elevator routines
	go control.Control(com_ch.Update_outgoing_msg_ch, com_ch.New_peer_ch, com_ch.Lost_peers_ch,com_ch.Update_control_variables_ch,esm_ch.Extern_order_ch, esm_ch.New_order_ch,esm_ch.State_ch,esm_ch.Illuminate_extern_order_ch,esm_ch.Clear_lights_and_extern_orders_ch)
	go esm.EventHandler(esm_ch, start_floor)
	go com.Communication(com_ch)

//Transmitters/Receivers
	go bcast.Transmitter(12345, com_ch.Outgoing_msg_ch)
	go bcast.Receiver(12345, com_ch.Incoming_msg_ch)
	go peers.Transmitter(15647, id, peer_trans_en_ch)
	go peers.Receiver(15647, com_ch.Peer_update_ch)


	select {}
}
