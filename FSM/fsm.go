// FINITE STATE MACHINE
package FSM

import (
	//"fmt"
	//windows:
	//"driver/elevio"
	//linux:
	"../driver/elevio"
)

// STATE MACHINE //
func Fsm(clear_lights_and_extern_orders_ch <-chan int ,start_floor int, cancel_illuminate_extern_order_ch <-chan int, illuminate_extern_order_ch <-chan elevio.ButtonEvent, door_timer_ch <-chan bool, extern_order_ch <-chan elevio.ButtonEvent, buttons_ch <-chan elevio.ButtonEvent, floors_ch <-chan int /*receiveing channels*/, reached_extern_floor_ch chan<- elevio.ButtonEvent, new_order_ch chan<- elevio.ButtonEvent, state_ch chan<- Elevator, reset_timer_ch chan<- bool) {
	elevator.Last_known_floor = start_floor
	for {
		select {
		case button_pushed := <-buttons_ch:

			button_event(button_pushed, new_order_ch, reset_timer_ch)
			fsm_print()
			go func() { state_ch <- elevator}()

		case floor := <-floors_ch:
			floor_event(floor, reset_timer_ch)
			fsm_print()
			go func() { state_ch <- elevator }()

		case <-door_timer_ch:
			door_open_event()
			fsm_print()
			go func() { state_ch <- elevator }()

		case order := <-extern_order_ch:
			extern_order_event(order, reset_timer_ch)
			go func() { state_ch <- elevator }()

		case order := <-illuminate_extern_order_ch:
			elevio.SetButtonLamp(order.Button, order.Floor, true)

		case floor := <-cancel_illuminate_extern_order_ch:
			clear_all_lights_on_floor(floor)
		case floor:=<-clear_lights_and_extern_orders_ch:
			clear_extern_ligts_on_floor(floor)
			clear_extern_order_on_floor(floor)
		}
	}
}