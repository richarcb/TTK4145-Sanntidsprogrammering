package esm

import (
	"../driver/elevio"
. "../config"
)

//Function that remove element from queue:
func remove_elem(index int) {
	for i := index; i < (len(queue) - 1); i++ {
		queue[i] = queue[i+1]
		if queue[i].Floor == empty_order.Floor {
			break
		}
	}
}

//Inserts element in the front of queue
func insert_front(front_elem Order) {

	for i := len(queue) - 1; i > 0; i-- {
		queue[i] = queue[i-1]
	}
	queue[0] = front_elem
}

//Adds order in queue
func push_back(elem Order) {
	for i := 0; i < len(queue); i++ {
		if queue[i].Floor == elem.Floor && queue[i].Button == elem.Button {
			return
		}
		if queue[i].Floor == empty_order.Floor {
			queue[i] = elem
			break
		}
	}
}

//Checks if elevator have an order on its way to destination
func check_for_extra_stop() {
	extra_stop = empty_order

	switch elevator.Dir {
	case MD_Stop:
		if elevator.Last_known_floor < elevator.Destination.Floor { //going up
			for i := elevator.Last_known_floor + 1; i < elevator.Destination.Floor; i++ {
				if intern_order_list[i] == 1 && (i < extra_stop.Floor || extra_stop.Floor == empty_order.Floor) {
					extra_stop = Order{Floor: i, Button: BT_Cab}
					break
				}
			}
			for i := 0; i < len(queue); i++ {
				if queue[i].Floor == empty_order.Floor {
					break
				} else if queue[i].Button == BT_HallUp && queue[i].Floor > elevator.Last_known_floor && queue[i].Floor < elevator.Destination.Floor {
					if extra_stop.Floor == empty_order.Floor || queue[i].Floor < extra_stop.Floor {
						extra_stop = queue[i]
					}
				}
			}
		} else if elevator.Last_known_floor > elevator.Destination.Floor { //going down
			for i := elevator.Last_known_floor - 1; i > elevator.Destination.Floor; i-- {
				if intern_order_list[i] == 1 && i > extra_stop.Floor {
					extra_stop = Order{Floor: i, Button: BT_Cab}
					break
				}
			}
			for i := 0; i < len(queue); i++ {
				if queue[i].Button == BT_HallDown && queue[i].Floor > extra_stop.Floor && queue[i].Floor < elevator.Last_known_floor && queue[i].Floor > elevator.Destination.Floor {
					extra_stop = queue[i]
				}
			}
		}
	case MD_Up:
		for i := elevator.Last_known_floor + 1; i < elevator.Destination.Floor; i++ {
			if intern_order_list[i] == 1 && (i < extra_stop.Floor || extra_stop.Floor == empty_order.Floor) {
				extra_stop = Order{Floor: i, Button: BT_Cab}
				break
			}
		}
		for i := 0; i < len(queue); i++ {
			if queue[i].Floor == empty_order.Floor {
				break
			}
			if queue[i].Button == BT_HallUp && queue[i].Floor > elevator.Last_known_floor && queue[i].Floor < elevator.Destination.Floor {
				if extra_stop.Floor == empty_order.Floor || queue[i].Floor < extra_stop.Floor {
					extra_stop = queue[i]
				}
			}
		}
	case MD_Down:
		for i := elevator.Last_known_floor - 1; i > elevator.Destination.Floor; i-- {
			if intern_order_list[i] == 1 && i > extra_stop.Floor {
				extra_stop = Order{Floor: i, Button: BT_Cab}
				break
			}
		}
		for i := 0; i < len(queue); i++ {
			if queue[i].Button == BT_HallDown && queue[i].Floor > extra_stop.Floor && queue[i].Floor < elevator.Last_known_floor && queue[i].Floor > elevator.Destination.Floor {
				extra_stop = queue[i]
			}
		}
	}
}

func find_new_destination(priority bool) {
	if elevator.State == MOVING {
		return
	}
	new_dest := empty_order
	if priority {
		for i := len(intern_order_list) - 1; i >= 0; i-- {
			if intern_order_list[i] == 1 {
				new_dest.Floor = i
				new_dest.Button = BT_Cab
				elevator.Destination = new_dest
				intern_order_list[i] = 0
				return
			}
		}
	} else {
		for i := 0; i < len(intern_order_list); i++ {
			if intern_order_list[i] == 1 {
				new_dest.Floor = i
				new_dest.Button = BT_Cab
				elevator.Destination = new_dest
				intern_order_list[i] = 0
				return
			}
		}
	}

	if queue[0].Floor != empty_order.Floor {
		new_dest = queue[0]
		remove_elem(0)
	}
	elevator.Destination = new_dest
}

func close_door() {
	elevio.SetDoorOpenLamp(false)
	elevator.State = IDLE
}

func open_door() {
	if elevator.Dir != MD_Stop {
		return
	}
	elevio.SetDoorOpenLamp(true)
	elevator.State = DOOROPEN
}

//Checks if direction should be updated depending on current orders
func update_direction() {
	if extra_stop.Floor == elevator.Last_known_floor || elevator.Destination.Floor == elevator.Last_known_floor {
		elevator.Dir = MD_Stop
	} else if elevator.Destination.Floor > elevator.Last_known_floor {
		elevator.Dir = MD_Up
	} else if elevator.Destination.Floor < elevator.Last_known_floor && elevator.Destination.Floor != empty_order.Floor {
		elevator.Dir = MD_Down
	} else {
		elevator.Dir = MD_Stop
	}
}

//Finite state machine button event
func button_event(button_pushed Order, new_order_ch chan<- Order, reset_door_timer_ch chan<- bool, reset_power_loss_timer_ch chan<- bool) {
	if button_pushed.Button == BT_Cab {
		switch elevator.State {
		case IDLE:
			if button_pushed.Floor == elevator.Last_known_floor {
				elevator.Destination = button_pushed
				open_door()
				reset_door_timer_ch <- true
			} else {
				elevio.SetButtonLamp(BT_Cab, button_pushed.Floor, true)
				//Inside order: Add order to intern list, no need for sharing
				intern_order_list[button_pushed.Floor] = 1
				find_new_destination(false)
				update_direction()
				elevio.SetMotorDirection(elevator.Dir)
				if elevator.Dir != MD_Stop {
					elevator.State = MOVING
					go func() { reset_power_loss_timer_ch <- true }()
				}
			}

		case MOVING:
			elevio.SetButtonLamp(BT_Cab, button_pushed.Floor, true)
			//Inside order: Add order to intern list, no need for sharing
			intern_order_list[button_pushed.Floor] = 1
			check_for_extra_stop()

		case DOOROPEN:
			if button_pushed.Floor == elevator.Last_known_floor {
				if elevator.Destination.Floor == empty_order.Floor {
					elevator.Destination = button_pushed
				}
				open_door()
				reset_door_timer_ch <- true
			} else {
				//Inside order: Add order to intern list, no need for sharing
				intern_order_list[button_pushed.Floor] = 1
				elevio.SetButtonLamp(BT_Cab, button_pushed.Floor, true)
			}
		case POWERLOSS:
			elevio.SetButtonLamp(BT_Cab, button_pushed.Floor, true)
			//Inside order: Add order to intern list, no need for sharing
			intern_order_list[button_pushed.Floor] = 1
		}
		update_backup(intern_order_list, elevator.Destination)
	} else if button_pushed.Floor == elevator.Last_known_floor && elevator.State != MOVING {
		open_door()
		reset_door_timer_ch <- true
	} else { //Send order to other Control module
		new_order := Order{Floor: button_pushed.Floor, Button: button_pushed.Button}
		go func() { new_order_ch <- new_order }()
	}
}

//Clear lights in hall buttons on floor
func clear_extern_ligts_on_floor(floor int) {
	for i := BT_HallUp; i <= BT_HallDown; i++ {
		elevio.SetButtonLamp(i, floor, false)
	}
}

//Clear hall orders on florr
func clear_extern_order_on_floor(floor int) {
	for i := 0; i < len(queue); i++ {
		if queue[i].Floor == floor {
			remove_elem(i)
		}
	}
}

func clear_all_lights_on_floor(floor int) {
	for i := BT_HallUp; i <= BT_Cab; i++ {
		elevio.SetButtonLamp(i, floor, false)
	}
}

func clear_all_order_on_floor(floor int) {
	intern_order_list[floor] = 0
	for i := 0; i < len(queue); i++ {
		if queue[i].Floor == empty_order.Floor{break}
		if queue[i].Floor == floor {
			remove_elem(i)
			i--;
		}
	}
}

//Finite state machine powerloss event
func power_loss_event(stop_power_loss_timer_ch chan<- bool) {
	if elevator.State == MOVING{
		elevator.State = POWERLOSS
		for i := 0; i < N_floors; i++ {
			clear_extern_order_on_floor(i)
			clear_extern_ligts_on_floor(i)
		}
	}else{
		stop_power_loss_timer_ch<-true
	}

}

//Finite state machine floor sensor event
func floor_event(floor int, reset_door_timer_ch chan<- bool, stop_power_loss_timer_ch chan<- bool, reset_power_loss_timer_ch chan<- bool) {
	elevator.Last_known_floor = floor

	switch elevator.State {
	case IDLE://This case will only occur when initializing elevator
		on_floor:=false
		for i:=0;i<len(intern_order_list);i++{
			if intern_order_list[i] == 1 && elevator.Last_known_floor == i{
				elevator.Destination = Order{Button: BT_Cab, Floor: i}
				clear_all_lights_on_floor(i)
				clear_all_order_on_floor(i)
				open_door()
				reset_door_timer_ch <- true
				on_floor = true
			}
		}
		if !on_floor{
			find_new_destination(false)
			check_for_extra_stop()
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			if elevator.Dir != MD_Stop{
				elevator.State = MOVING
				go func() { reset_power_loss_timer_ch <- true }()
			}
		}

	case POWERLOSS:
		if elevator.Destination.Floor == floor || extra_stop.Floor == floor {
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			clear_all_lights_on_floor(floor)
			clear_all_order_on_floor(floor)
			open_door()
			reset_door_timer_ch <- true
		} else {
			elevator.State = MOVING
			go func() { reset_power_loss_timer_ch <- true }()
		}
	case MOVING:
		go func() { reset_power_loss_timer_ch <- true }()
		if elevator.Destination.Floor == floor || extra_stop.Floor == floor {
			go func() { stop_power_loss_timer_ch <- true }()
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			clear_all_lights_on_floor(floor)
			clear_all_order_on_floor(floor)
			open_door()
			reset_door_timer_ch <- true
		}
	}
}

//Finite state machine event for opening door event
func door_timer_event(reset_power_loss_timer_ch chan<- bool) {
	switch elevator.State {
	case DOOROPEN:
		close_door()
		if elevator.Destination.Floor == elevator.Last_known_floor {
			priorityVariable := false
			if elevator.Destination.Button == BT_HallUp {
				priorityVariable = true
			}
			elevator.Destination = empty_order
			find_new_destination(priorityVariable)
		}
		check_for_extra_stop()
		update_direction()
		elevio.SetMotorDirection(elevator.Dir)
		if elevator.Dir != MD_Stop{
			elevator.State = MOVING
			go func(){reset_power_loss_timer_ch<- true}()
		}
	}
	update_backup(intern_order_list, elevator.Destination)
}

//Finite state machine for getting an external order from the control module
func extern_order_event(order Order, reset_door_timer_ch chan<- bool, reset_power_loss_timer_ch chan<- bool) {
	switch elevator.State {
	case IDLE:
		if order.Floor == elevator.Last_known_floor {
			elevator.Destination = order
			clear_all_order_on_floor(order.Floor)
			open_door()
			reset_door_timer_ch <- true
		} else {
			elevio.SetButtonLamp(order.Button, order.Floor, true)
			push_back(order)
			find_new_destination(false)
			update_direction()
			elevio.SetMotorDirection(elevator.Dir)
			if elevator.Dir != MD_Stop {
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
			if elevator.Destination.Floor == empty_order.Floor {
				elevator.Destination = order
			}
			open_door()
			reset_door_timer_ch <- true
		} else {
			push_back(order)
			elevio.SetButtonLamp(order.Button, order.Floor, true)
			check_for_extra_stop()
		}
	}
}
