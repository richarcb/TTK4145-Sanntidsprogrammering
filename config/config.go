package config

//Scalable declaration of the number of floors and elevators
const (
	N_floors    = 4
	N_elevators = 3
)


type Order struct {
	Floor  int
	Button ButtonType
}

type MotorDirection int
const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int
const (
	BT_HallUp   ButtonType = 0
	BT_HallDown           = 1
	BT_Cab         			= 2
	BT_No_call 					= -1
)

type Elevator struct {
	Destination      Order
	Last_known_floor int
	Dir              MotorDirection
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
	Destination      Order
	Last_known_floor int
	Dir              MotorDirection
	State            ElevState
	ID               string
	Ack_list				 [2][N_floors]int
}
