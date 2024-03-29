package esm

import (
	"time"
)

//This file contains timers for door open and powerloss

func door_timer(close_door_ch chan<- bool, reset_timer <-chan bool) {
	const doorOpenTime = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-reset_timer:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			close_door_ch <- true
		}
	}
}

func powerloss_timer(power_loss_ch chan<- bool, reset_power_loss_timer_ch <-chan bool, stop_power_loss_timer_ch <-chan bool) {
	const no_floor_time = 5 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-reset_power_loss_timer_ch:
			timer.Reset(no_floor_time)
		case <-timer.C:
			timer.Stop()
			go func() { power_loss_ch <- true }()
		case <-stop_power_loss_timer_ch:
			timer.Stop()
		}
	}
}
