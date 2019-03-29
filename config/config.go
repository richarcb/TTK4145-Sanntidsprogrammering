//ELEVATOR CONFIGURATION
package config
//import "flag"
/*
import "os"
import "strconv"
*/

import(

	"../driver/elevio"
)



const (
	//System variables
	N_floors    = 4
	N_elevators = 3

	//Local variables
)
var ElevatorNumber int

type Elevator struct {
	//Destination floor
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            ElevState
}

type ElevState int

const (
	IDLE     ElevState = 0
	MOVING             = 1
	DOOROPEN           = 2
	POWERLOSS					 = 3
)


type Msg_struct struct {
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            ElevState
	ID               string
	Ack_list				 [2][N_floors]int
}
