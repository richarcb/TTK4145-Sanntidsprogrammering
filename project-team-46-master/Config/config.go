//ELEVATOR CONFIGURATION
package config
//import "flag"
/*
import "os"
import "strconv"
*/




const (
	//System variables
	N_floors    = 4
	N_elevators = 3

	//Local variables
)
var ElevatorNumber int
/*
func Init_elevconfig(){
	ElevatorString := os.Args[2]
	ElevatorNumber, _ = strconv.Atoi(ElevatorString)
}*/


/*
type Msg_struct struct {
	//Destination floor
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            FSM.ElevState
	ID               int
	Ack_list				 [N_elevators][N_floors]int
	//IP 							 string
}

type elevator_states struct {
	//Destination floor
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            FSM.ElevState
	ID				 			 int
	queue            [2][ N_floors]int
}
*/
