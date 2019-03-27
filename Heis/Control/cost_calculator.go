package control
/*package control
import (
	//windows:
	//"time"

	config "../Config"
	"../FSM"
	//sync "../Synchronizing"
	"../driver/elevio"
	//"config"
	//"fmt"
	//linux:
	//"../driver/elevio"
)
func cost_function (id int, order elevio.ButtonEvent) int{
  cost := 0
  if order.Button == elevio.BT_HallUp{//Order is up
   if elevator_list[id].last_known_floor<order.Floor && elevator_list[id].destination.Floor>order.Floor{//going up and flor is between:
      cost -= 10
    }
  }else{//Order is down
    order.Floor -= config.N_floors-2 //config.N_floors-2 since the first down button is in the second floor!
      if elevator_list[id].last_known_floor>order.Floor && elevator_list[id].destination.Floor<order.Floor && elevator_list[id].state == FSM.MOVING{//going down and floor is between orders:Needs to check MOVING since no destination 0-1<order.Floor
        cost -= 10
    }
  }
    if elevator_list[id].state == FSM.IDLE && elevator_list[id].destination.Floor == FSM.Empty_order.Floor{ //Nothing to do
      cost -= 5
    }
  //Adding the value of the distance
  if order.Floor > elevator_list[id].last_known_floor{
    cost += order.Floor-elevator_list[id].last_known_floor
  }else{
    cost += elevator_list[id].last_known_floor-order.Floor
  }
  return cost
}

func getLowestCostElevatorID(order elevio.ButtonEvent) int{
	lowestCost := config.N_floors
	assignedID := -1
  //Get Number_of_Online_elevators! (POWERLOSS???)

	for i := 0; i <= config.N_elevators; i++{
    if offline_elevator_list[i]{
      continue
      }
		cost := cost_function(i, order)
		if cost < lowestCost{
			lowestCost = cost
			assignedID = i
		}
	}
	return assignedID
}
*/
