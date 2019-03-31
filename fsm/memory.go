package fsm

import (
	"../driver/elevio"
. "../config"
	"../backup"

)
//MEMORY START

var extra_stop Order
var Empty_order Order

//External elevator queue
var queue [N_floors * 2]Order //Creates empty queue (Should be Linked list)

//Local elevator queue
var intern_order_list [N_floors]int

//Local elevator variable
var elevator Elevator

//MEMORY END


//Initializes elevator's memory. Memory is defined above
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
