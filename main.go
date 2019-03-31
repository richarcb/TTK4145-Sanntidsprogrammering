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
	com "./network/communication"
	"./driver/elevio"
)
//Initializes elevator's IP config
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


//15657 --> default simulator port
	elevio.Init("localhost:"+port, N_floors)
	fsm.Init_mem()
 	return elevio.InitElev(),id
}


func main() {
	start_floor,id := initialize_elevator_system()

	//Local fsm channels
	fsm_ch := fsm.Fsm_channels{
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
	go elevio.PollButtons(fsm_ch.Buttons_ch)
	go elevio.PollFloorSensor(fsm_ch.Floors_ch)

//Elevator routines
	go control.Control(com_ch.Update_outgoing_msg_ch, com_ch.New_peer_ch, com_ch.Lost_peers_ch,com_ch.Update_control_variables_ch,fsm_ch.Extern_order_ch, fsm_ch.New_order_ch,fsm_ch.State_ch,fsm_ch.Illuminate_extern_order_ch,fsm_ch.Clear_lights_and_extern_orders_ch)
	go fsm.EventHandler(fsm_ch, start_floor)
	go com.Communication(com_ch)

//Transmitters/Receivers
	go bcast.Transmitter(12345, com_ch.Outgoing_msg_ch)
	go bcast.Receiver(12345, com_ch.Incoming_msg_ch)
	go peers.Transmitter(15647, id, peer_trans_en_ch)
	go peers.Receiver(15647, com_ch.Peer_update_ch)


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

func RandomOrderGenerator(receiver chan Order) {
	var order Order
	BT_type := rand.Intn(3)
	if BT_type == 0 {
		order.Button = BT_HallUp
		order.Floor = rand.Intn(3)
	} else if BT_type == 1 {
		order.Button = BT_HallDown
		order.Floor = rand.Intn(3) + 1
	} else {
		order.Button = BT_Cab
		order.Floor = rand.Intn(4)
	}

	receiver <- order
}*/
