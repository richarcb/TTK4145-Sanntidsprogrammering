package fsm

import (
	"../driver/elevio"
. "../config"
	"../backup"
	
)
var extra_stop Order
var Empty_order Order
var queue [N_floors * 2]Order //Creates empty queue (Should be Linked list)
var intern_order_list [N_floors]int
var elevator Elevator

func Init_mem() {
	Empty_order.Floor = -1
	Empty_order.Button = BT_No_call
	for i := 0; i < len(queue); i++ {
		queue[i] = Empty_order
	}
	for i:=0; i<N_floors; i++{
		clear_all_lights_on_floor(i)
	}

	if backup.BackupExists() {
		intern_order_list = backup.ReadFromBackup()
		for i := 0; i < N_floors; i++ {
			if intern_order_list[i] == 1 {
				elevio.SetButtonLamp(BT_Cab, i, true)
			}
		}

	}
	extra_stop = Empty_order

	elevator.Destination = Empty_order
	elevator.Last_known_floor = -1
	elevator.Dir = MD_Stop
	elevator.State = IDLE
}
