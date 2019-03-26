package sync

import ( /*"fmt"
	  "../network/peers"
	  "../driver/elevio"
	  "../FSM"
	  "../Config"*/

	config "../Config"
	"../FSM"
	"../driver/elevio"

	"../network/localip"
)



//elevator dead or power loss (Maybe no networkconnection):


type Msg_struct struct {
	//Destination floor
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            FSM.ElevState
	ID               int
	Ack_list				 [2][config.N_floors]int
	IP 							 string
}

var outgoing_msg Msg_struct


func Init_sync(init_outgoing_msg_ch chan<- Msg_struct) {

	var empty_ack_list [2][config.N_floors]int // Not shallow copy!
	for i := 0; i < 2; i++ {
		for j := 0; j < config.N_floors; j++ {
			empty_ack_list[i][j] = 0
		}
	}
	elevIP,_ := localip.LocalIP()
	empty_msg := Msg_struct{Destination: FSM.Empty_order, Last_known_floor: -1, Dir: elevio.MD_Stop, State: FSM.IDLE,  Ack_list: empty_ack_list, ID: config.ID, IP: elevIP}
	outgoing_msg = empty_msg
	go func(){init_outgoing_msg_ch<-outgoing_msg}()


}



/*
func Update_Msg_struct(){

}

func TransmitMsg(incoming_msg_ch chan<- Msg_struct){
	bcastTicker := time.NewTicker(300*time.Millisecond)
	  for{
		select{
		case <- bcastTicker.C:
		  //UpdateElevatorStruct()
		  incoming_msg_ch <- outgoing_msg
		  //fmt.Printf("STATE SENT \n")
		  //fmt.Printf("STATE SENT \n")

		}
	  }
  }

  func ReceiveMsg(outgoing_msg_ch <- chan Msg_struct){
	  for{
		select{
		case msg := <- outgoing_msg_ch:

		  fmt.Printf("STATE RECEIVED FROM: %v\n", msg.IP)
		}
		}
	  }
*/
