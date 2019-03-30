// FINITE STATE MACHINE
package fsm

import (
 ."../config"
	"../driver/elevio"
)

type FsmChannels struct{
  Cancel_illuminate_extern_order_ch chan int
	Illuminate_extern_order_ch chan elevio.ButtonEvent
	Extern_order_ch chan elevio.ButtonEvent
	Buttons_ch chan elevio.ButtonEvent
	Floors_ch chan int
	Reached_extern_floor_ch chan elevio.ButtonEvent
	New_order_ch chan elevio.ButtonEvent
	State_ch chan Elevator
}

// STATE MACHINE //
func EventHandler(ch FsmChannels, start_floor int, clear_lights_and_extern_orders_ch chan int) {
	elevator.Last_known_floor = start_floor
	go func() { ch.State_ch <- elevator }()

	door_timer_ch := make(chan bool)
	reset_timer_ch := make(chan bool)
	power_loss_ch := make(chan bool)
	reset_power_loss_timer_ch := make(chan bool)
	stop_power_loss_timer_ch := make(chan bool)
	go DoorTimer(door_timer_ch, reset_timer_ch)
	go Powerloss_timer(power_loss_ch, reset_power_loss_timer_ch, stop_power_loss_timer_ch)

	for {
		select {
		case button_pushed := <-ch.Buttons_ch:
			button_event(button_pushed, ch.New_order_ch, reset_timer_ch, reset_power_loss_timer_ch)
			fsm_print()
			go func() { ch.State_ch <- elevator }()
		case floor := <-ch.Floors_ch:
			elevio.SetFloorIndicator(floor)
			floor_event(floor, reset_timer_ch, stop_power_loss_timer_ch, reset_power_loss_timer_ch)
			fsm_print()
			go func() { ch.State_ch <- elevator }()

		case <-door_timer_ch:
			door_open_event(reset_power_loss_timer_ch)
			fsm_print()
			go func() { ch.State_ch <- elevator }()

		case order := <-ch.Extern_order_ch:
			extern_order_event(order, reset_timer_ch, reset_power_loss_timer_ch)
			go func() { ch.State_ch <- elevator }()

		case order := <-ch.Illuminate_extern_order_ch:
			elevio.SetButtonLamp(order.Button, order.Floor, true)

		case <-power_loss_ch:
			power_loss_event(stop_power_loss_timer_ch)
			go func() { ch.State_ch <- elevator }()

		case floor := <-ch.Cancel_illuminate_extern_order_ch:
			clear_all_lights_on_floor(floor)

		case floor := <-clear_lights_and_extern_orders_ch:
			clear_extern_ligts_on_floor(floor)
			clear_extern_order_on_floor(floor)
		}
	}
}
