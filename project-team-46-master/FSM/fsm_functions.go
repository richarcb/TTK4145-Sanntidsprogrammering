package FSM

import (
	//windows:
	//"container/list"
	//"driver/elevio"
	"fmt"

	"../driver/elevio"
	//"../backup"
	//linux:

	config "../Config"
	"../backup"
	/*
		"fmt"
		"../driver/elevio"
	*/)

//Function that removes middle element from queue:
//Linked list, not queue! Can perhaps remove elements from middel of list!
func remove_elem(index int) {
	for i := index; i < (len(queue) - 1); i++ {
		queue[i] = queue[i+1]
		if queue[i].Floor == Empty_order.Floor {
			break
		}
	}
}
func insert_front(front_elem elevio.ButtonEvent) {

	for i := len(queue) - 1; i > 0; i-- {
		queue[i] = queue[i-1]
	}
	queue[0] = front_elem
}

func push_back(elem elevio.ButtonEvent) {
	for i := 0; i < len(queue); i++ {
		if queue[i].Floor == elem.Floor && queue[i].Button == elem.Button {
			return
		}
		if queue[i].Floor == Empty_order.Floor {
			queue[i] = elem
			break
		}
	}
}

func check_for_extra_stop() {

	extra_stop = Empty_order

	switch elevator.Dir {
	case elevio.MD_Stop:
		if elevator.Last_known_floor < elevator.Destination.Floor { //UP
			for i := elevator.Last_known_floor + 1; i < elevator.Destination.Floor; i++ {
				if intern_order_list[i] == 1 && (i < extra_stop.Floor || extra_stop.Floor == Empty_order.Floor) {
					extra_stop = elevio.ButtonEvent{Floor: i, Button: elevio.BT_Cab}
					break
				}
			}
			for i := 0; i < len(queue); i++ {
				if queue[i].Floor == Empty_order.Floor {
					break
				} else if queue[i].Button == elevio.BT_HallUp && queue[i].Floor > elevator.Last_known_floor && queue[i].Floor < elevator.Destination.Floor {
					if extra_stop.Floor == Empty_order.Floor || queue[i].Floor < extra_stop.Floor {
						extra_stop = queue[i]
					}
				}
			}
		} else if elevator.Last_known_floor > elevator.Destination.Floor { //Down
			for i := elevator.Last_known_floor - 1; i > elevator.Destination.Floor; i-- {
				if intern_order_list[i] == 1 && i > extra_stop.Floor {
					extra_stop = elevio.ButtonEvent{Floor: i, Button: elevio.BT_Cab}
					break
				}
			}
			for i := 0; i < len(queue); i++ {
				if queue[i].Button == elevio.BT_HallDown && queue[i].Floor > extra_stop.Floor && queue[i].Floor < elevator.Last_known_floor && queue[i].Floor > elevator.Destination.Floor {
					extra_stop = queue[i]
				}
			}
		}
	case elevio.MD_Up:
		for i := elevator.Last_known_floor + 1; i < elevator.Destination.Floor; i++ {
			if intern_order_list[i] == 1 && (i < extra_stop.Floor || extra_stop.Floor == Empty_order.Floor) {
				extra_stop = elevio.ButtonEvent{Floor: i, Button: elevio.BT_Cab}
				break
			}
		}
		for i := 0; i < len(queue); i++ {
			if queue[i].Floor == Empty_order.Floor {
				break
			}
			if queue[i].Button == elevio.BT_HallUp && queue[i].Floor > elevator.Last_known_floor && queue[i].Floor < elevator.Destination.Floor {
				if extra_stop.Floor == Empty_order.Floor || queue[i].Floor < extra_stop.Floor {
					extra_stop = queue[i]
				}
			}
		}
	case elevio.MD_Down:
		for i := elevator.Last_known_floor - 1; i > elevator.Destination.Floor; i-- {
			if intern_order_list[i] == 1 && i > extra_stop.Floor {
				extra_stop = elevio.ButtonEvent{Floor: i, Button: elevio.BT_Cab}
				break
			}
		}
		for i := 0; i < len(queue); i++ {
			if queue[i].Button == elevio.BT_HallDown && queue[i].Floor > extra_stop.Floor && queue[i].Floor < elevator.Last_known_floor && queue[i].Floor > elevator.Destination.Floor {
				extra_stop = queue[i]
			}
		}
	}
}

func find_new_destination(priority bool) {
	//Could add priorityvariable to get better features... NOW: Just chosing the lowest intern order first...
	if elevator.State == MOVING {
		return
	}
	new_dest := Empty_order
	if priority {
		for i := len(intern_order_list) - 1; i >= 0; i-- {
			if intern_order_list[i] == 1 {
				new_dest.Floor = i
				new_dest.Button = elevio.BT_Cab
				elevator.Destination = new_dest
				intern_order_list[i] = 0
				return
			}
		}
	} else {
		for i := 0; i < len(intern_order_list); i++ {
			if intern_order_list[i] == 1 {
				new_dest.Floor = i
				new_dest.Button = elevio.BT_Cab
				elevator.Destination = new_dest
				intern_order_list[i] = 0
				return
			}
		}
	}

	if queue[0].Floor != Empty_order.Floor {
		new_dest = queue[0]
		remove_elem(0)
	}
	elevator.Destination = new_dest
}

func close_door() {
	elevio.SetDoorOpenLamp(false)
	elevator.State = IDLE
}

func open_door() { //On floor, doors open
	if elevator.Dir != elevio.MD_Stop {
		return
	}
	elevio.SetDoorOpenLamp(true)
	elevator.State = DOOROPEN
	//start timer
}

func update_direction() {
	//Never called whene
	if extra_stop.Floor == elevator.Last_known_floor || elevator.Destination.Floor == elevator.Last_known_floor {
		elevator.Dir = elevio.MD_Stop
	} else if elevator.Destination.Floor > elevator.Last_known_floor {
		elevator.Dir = elevio.MD_Up
	} else if elevator.Destination.Floor < elevator.Last_known_floor && elevator.Destination.Floor != Empty_order.Floor {
		elevator.Dir = elevio.MD_Down
	} else {
		elevator.Dir = elevio.MD_Stop
	}
}

func drive() { //Drive
	if elevator.State == DOOROPEN {
		return
	}
	elevio.SetMotorDirection(elevator.Dir)
}

func button_event(button_pushed elevio.ButtonEvent, new_order_ch chan<- elevio.ButtonEvent, reset_timer_ch chan<- bool, reset_power_loss_timer_ch chan<- bool) {
	if button_pushed.Button == elevio.BT_Cab {
		switch elevator.State {
		case IDLE:
			fsm_print()
			if button_pushed.Floor == elevator.Last_known_floor {
				elevator.Destination = button_pushed
				open_door()
				reset_timer_ch <- true
			} else {
				elevio.SetButtonLamp(elevio.BT_Cab, button_pushed.Floor, true)
				//Inside order: Add order to intern list, no need for sharing
				intern_order_list[button_pushed.Floor] = 1
				find_new_destination(false)
				update_direction()
				elevio.SetMotorDirection(elevator.Dir)
				if elevator.Dir != elevio.MD_Stop {
					elevator.State = MOVING
					go func() { reset_power_loss_timer_ch <- true }()
				}
			}

		case MOVING:
			fsm_print()
			elevio.SetButtonLamp(elevio.BT_Cab, button_pushed.Floor, true)
			//Inside order: Add order to intern list, no need for sharing
			intern_order_list[button_pushed.Floor] = 1
			check_for_extra_stop()

		case DOOROPEN:
			fsm_print()
			if button_pushed.Floor == elevator.Last_known_floor {
				if elevator.Destination.Floor == Empty_order.Floor { //BUGFIX
					elevator.Destination = button_pushed
				}
				open_door()
				reset_timer_ch <- true
			} else {
				intern_order_list[button_pushed.Floor] = 1
				elevio.SetButtonLamp(elevio.BT_Cab, button_pushed.Floor, true)
			}
		case POWERLOSS:
			elevio.SetButtonLamp(elevio.BT_Cab, button_pushed.Floor, true)
			intern_order_list[button_pushed.Floor] = 1

			/*fmt.Printf("opendoor")
			timerReset <- true*/
		}
		backup.UpdateBackup(intern_order_list, elevator.Destination) //New backup.
	} else if button_pushed.Floor == elevator.Last_known_floor && elevator.State != MOVING {
		open_door()
		reset_timer_ch <- true
	} else { //Send order to other Module
		new_order := elevio.ButtonEvent{Floor: button_pushed.Floor, Button: button_pushed.Button}
		go func() { new_order_ch <- new_order }()
	}
}
func clear_extern_ligts_on_floor(floor int) {
	for i := elevio.BT_HallUp; i <= elevio.BT_HallDown; i++ {
		elevio.SetButtonLamp(i, floor, false)
	}
}
func clear_extern_order_on_floor(floor int) {
	for i := 0; i < len(queue); i++ {
		if queue[i].Floor == floor {
			remove_elem(i)
		}
	}
}

func clear_all_lights_on_floor(floor int) {
	for i := elevio.BT_HallUp; i <= elevio.BT_Cab; i++ {
		elevio.SetButtonLamp(i, floor, false)
	}
}

func clear_all_order_on_floor(floor int) {
	intern_order_list[floor] = 0
	for i := 0; i < len(queue); i++ {
		if queue[i].Floor == floor {
			remove_elem(i)
		}
	}
}

func power_loss_event(stop_power_loss_timer_ch chan<- bool) {
	//QUICK // FIX
	if elevator.State == MOVING{
		elevator.State = POWERLOSS
		for i := 0; i < config.N_floors; i++ {
			clear_all_order_on_floor(i)
		}
	}else{
		stop_power_loss_timer_ch<-true
	}

}

func floor_event(floor int, reset_timer_ch chan<- bool, stop_power_loss_timer_ch chan<- bool, reset_power_loss_timer_ch chan<- bool) {
	elevator.Last_known_floor = floor

	switch elevator.State {
	case IDLE:
		on_floor:=false
		for i:=0;i<len(intern_order_list);i++{
			if intern_order_list[i] == 1 && elevator.Last_known_floor == i{

				//konfigurer lys
				elevator.Destination = elevio.ButtonEvent{Button: elevio.BT_Cab, Floor: i}
				clear_all_lights_on_floor(i)
				clear_all_order_on_floor(i)
				open_door()
				reset_timer_ch <- true
				on_floor = true
			}
		}
		if !on_floor{
			find_new_destination(false)
			check_for_extra_stop()
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			if elevator.Dir != elevio.MD_Stop{
				elevator.State = MOVING
				go func() { reset_power_loss_timer_ch <- true }()
			}
		}

	case POWERLOSS:
		if elevator.Destination.Floor == floor || extra_stop.Floor == floor {
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			//konfigurer lys
			clear_all_lights_on_floor(floor)
			clear_all_order_on_floor(floor)
			open_door()
			reset_timer_ch <- true
		} else {
			elevator.State = MOVING
			go func() { reset_power_loss_timer_ch <- true }()
		}
	case MOVING:
		go func() { reset_power_loss_timer_ch <- true }()
		fsm_print()
		if elevator.Destination.Floor == floor || extra_stop.Floor == floor {
			//Stops_at_floor:
			go func() { stop_power_loss_timer_ch <- true }()
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			//konfigurer lys
			clear_all_lights_on_floor(floor)
			clear_all_order_on_floor(floor)
			open_door()
			reset_timer_ch <- true
		}
	}
}

func door_open_event(reset_power_loss_timer_ch chan<- bool) {
	switch elevator.State {
	case DOOROPEN:
		close_door()
		if elevator.Destination.Floor == elevator.Last_known_floor {
			priorityVariable := false
			if elevator.Destination.Button == elevio.BT_HallUp {
				priorityVariable = true
			}
			elevator.Destination = Empty_order
			find_new_destination(priorityVariable)
		}
		check_for_extra_stop()
		update_direction()
		elevio.SetMotorDirection(elevator.Dir)
		if elevator.Dir != elevio.MD_Stop{
			elevator.State = MOVING
			go func(){reset_power_loss_timer_ch<- true}()
		}
	}
	backup.UpdateBackup(intern_order_list, elevator.Destination)
}

func extern_order_event(order elevio.ButtonEvent, reset_timer_ch chan<- bool, reset_power_loss_timer_ch chan<- bool) {
	switch elevator.State {
	case IDLE:
		if order.Floor == elevator.Last_known_floor {
			elevator.Destination = order
			clear_all_order_on_floor(order.Floor)
			open_door()
			reset_timer_ch <- true
		} else {
			elevio.SetButtonLamp(order.Button, order.Floor, true)
			//Inside order: Add order to intern list, no need for sharing
			push_back(order)
			find_new_destination(false)
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			if elevator.Dir != elevio.MD_Stop {
				elevator.State = MOVING
				go func() { reset_power_loss_timer_ch <- true }()
			}
		}
	case MOVING:
		elevio.SetButtonLamp(order.Button, order.Floor, true)
		push_back(order)
		check_for_extra_stop()
	case DOOROPEN:
		if order.Floor == elevator.Last_known_floor {
			if elevator.Destination.Floor == Empty_order.Floor {
				elevator.Destination = order
			}
			open_door()
			reset_timer_ch <- true
		} else {
			push_back(order)
			elevio.SetButtonLamp(order.Button, order.Floor, true)
			check_for_extra_stop()
		}
	}
}

func fsm_print() {

	fmt.Printf("-----------NEW UPDATE -----------\n")
	fmt.Printf("State: %#v\n", elevator.State)
	fmt.Printf("Floor: %#v\n", elevator.Last_known_floor)
	fmt.Printf("Direction: %#v\n", elevator.Dir)
	fmt.Printf("Extra_stop: %#v\n", extra_stop.Floor)
	fmt.Printf("Destination: %#v\n", elevator.Destination.Floor)
	/*fmt.Printf("Orders: \n")
	/*for i:=0; i<len(queue); i++{
		fmt.Printf("%#v", queue[i].Floorr
	}
	for i:=0; i<len(intern_order_list);i++{
		fmt.Printf("%#v\n", intern_order_list[i])
	}*/
	return
}
