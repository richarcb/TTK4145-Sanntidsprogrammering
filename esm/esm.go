// FINITE STATE MACHINE
package esm

import (
 ."../config"
	"../driver/elevio"
)

type Esm_channels struct{
  Clear_lights_and_extern_orders_ch chan int
	Illuminate_extern_order_ch chan Order
	Extern_order_ch chan Order
	Buttons_ch chan Order
	Floors_ch chan int
	New_order_ch chan Order
	State_ch chan Elevator
}

// STATE MACHINE //
func EventHandler(esm_ch Esm_channels, start_floor int) {
	elevator.Last_known_floor = start_floor
	go func() { esm_ch.State_ch <- elevator }()
  //Local channels
	door_timer_ch := make(chan bool)
	reset_door_timer_ch := make(chan bool)
	power_loss_ch := make(chan bool)
	reset_power_loss_timer_ch := make(chan bool)
	stop_power_loss_timer_ch := make(chan bool)

  //Esm timers
	go door_timer(door_timer_ch, reset_door_timer_ch)
	go powerloss_timer(power_loss_ch, reset_power_loss_timer_ch, stop_power_loss_timer_ch)


	for {
		select {
		case button_pushed := <-esm_ch.Buttons_ch:
			button_event(button_pushed, esm_ch.New_order_ch, reset_door_timer_ch, reset_power_loss_timer_ch)
			go func() { esm_ch.State_ch <- elevator }()
		case floor := <-esm_ch.Floors_ch:
			elevio.SetFloorIndicator(floor)
			floor_event(floor, reset_door_timer_ch, stop_power_loss_timer_ch, reset_power_loss_timer_ch)
			go func() { esm_ch.State_ch <- elevator }()

		case <-door_timer_ch:
			door_timer_event(reset_power_loss_timer_ch)
			go func() { esm_ch.State_ch <- elevator }()

		case order := <-esm_ch.Extern_order_ch:
			extern_order_event(order, reset_door_timer_ch, reset_power_loss_timer_ch)
			go func() { esm_ch.State_ch <- elevator }()

		case order := <-esm_ch.Illuminate_extern_order_ch:
			elevio.SetButtonLamp(order.Button, order.Floor, true)

		case <-power_loss_ch:
			power_loss_event(stop_power_loss_timer_ch)
			go func() { esm_ch.State_ch <- elevator }()

		case floor := <-esm_ch.Clear_lights_and_extern_orders_ch:
			clear_extern_ligts_on_floor(floor)
			clear_extern_order_on_floor(floor)
		}
	}
}
