package fsm

import (
	"../driver/elevio"
. "../config"
	"../backup"
	"fmt"
)
var extra_stop elevio.ButtonEvent
var Empty_order elevio.ButtonEvent
var queue [N_floors * 2]elevio.ButtonEvent //Creates empty queue (Should be Linked list)
var intern_order_list [N_floors]int
var elevator Elevator

func Init_mem() {
	Empty_order.Floor = -1
	Empty_order.Button = elevio.BT_No_call
	for i := 0; i < len(queue); i++ {
		queue[i] = Empty_order
	}
	for i:=0; i<N_floors; i++{
		clear_all_lights_on_floor(i)
	}

	if backup.BackupExists() {
		fmt.Printf("BACKUP EXISTS")
		intern_order_list = backup.ReadFromBackup()
		for i := 0; i < N_floors; i++ {
			if intern_order_list[i] == 1 {
				fmt.Printf("SETTING LIGHTS")
				elevio.SetButtonLamp(elevio.BT_Cab, i, true)
			}
		}

	}
	for i := 0; i < len(intern_order_list); i++ {
		fmt.Println(intern_order_list[i])
	}
	extra_stop = Empty_order

	elevator.Destination = Empty_order
	elevator.Last_known_floor = -1
	elevator.Dir = elevio.MD_Stop
	elevator.State = IDLE
}
