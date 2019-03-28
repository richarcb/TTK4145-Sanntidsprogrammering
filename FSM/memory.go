package FSM

import (
	//windows:
	"../driver/elevio"
	//linux:
	config "../Config"
	"../backup"
	"fmt"
	//"../driver/elevio"
)

type ElevState int

const (
	IDLE     ElevState = 0
	MOVING             = 1
	DOOROPEN           = 2
	POWERLOSS					 = 3
)

var extra_stop elevio.ButtonEvent

type Elevator struct {
	//Destination floor
	Destination      elevio.ButtonEvent
	Last_known_floor int
	Dir              elevio.MotorDirection
	State            ElevState
}

var Empty_order elevio.ButtonEvent

//Should be linked list!
var queue [config.N_floors * 2]elevio.ButtonEvent //Creates empty queue (Should be Linked list)

var intern_order_list [config.N_floors]int

var elevator Elevator

func Init_mem() {
	Empty_order.Floor = -1
	Empty_order.Button = elevio.BT_No_call
	for i := 0; i < len(queue); i++ {
		queue[i] = Empty_order
	}
	for i:=0; i<config.N_floors; i++{
		clear_all_lights_on_floor(i)
	}


	if backup.BackupExists() {
		fmt.Printf("BACKUP EXISTS")
		intern_order_list = backup.ReadFromBackup()
		for i := 0; i < config.N_floors; i++ {
			if intern_order_list[i] == 1 {
				fmt.Printf("SETTING LIGHTS")
				elevio.SetButtonLamp(elevio.BT_Cab, i, true)
			}
		}

	}
	for i := 0; i < len(intern_order_list); i++ {
		fmt.Println(intern_order_list[i])
	}
	//internal_order_list = backup.ReadFromBackup() // checking for backup
	//OBS! SETT LYS HVIS BACKUP EKSISTERER:
	/*
		for i:= 0; i<len(internal_order_list){
			if intern_order_list[i] == 1{
				SetButtonLamp(BT_Cab, i, true)
				}
			}
	*/
	extra_stop = Empty_order
	elevator.Destination = Empty_order
	elevator.Last_known_floor = -1
	elevator.Dir = elevio.MD_Stop
	elevator.State = IDLE
}
