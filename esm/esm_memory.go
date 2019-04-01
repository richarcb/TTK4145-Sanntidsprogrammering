package esm

import (
	"../driver/elevio"
. "../config"


)
//MEMORY START

var extra_stop Order
var empty_order Order

//Hall order elevator queue
var queue [N_floors * 2]Order

//Cab order elevator queue
var intern_order_list [N_floors]int

//elavtor variables
var elevator Elevator

//MEMORY END


//Initializes elevator's memory. Memory is defined above
func Init_mem() {
	empty_order.Floor = -1
	empty_order.Button = BT_No_call
	for i := 0; i < len(queue); i++ {
		queue[i] = empty_order
	}
	for i:=0; i<N_floors; i++{
		clear_all_lights_on_floor(i)
	}

	if backup_exist() {
		intern_order_list = read_from_backup()
		for i := 0; i < N_floors; i++ {
			if intern_order_list[i] == 1 {
				elevio.SetButtonLamp(BT_Cab, i, true)
			}
		}

	}
	extra_stop = empty_order
	elevator.Destination = empty_order
	elevator.Last_known_floor = -1
	elevator.Dir = MD_Stop
	elevator.State = IDLE
}
