package queue

Floors = 4
Lifts = 3


type order struct {
	floor int
	orderType int // 0 for up, 1 for down and 2 for command
	//elevatorId int //Which elevator the order is from
}

type position struct {
	CurrentFloor int
	DestinationFloor int
}

type direction int8

var OrdersToComplete[Floors*3] order

//Only for master
var LiftPos[Lifts] position
var MasterArrayOfOrders[Lifts][Floors*3] order
var cost = 100000000



func Direction(LiftPos position) int {
	Direction = -1 //0 for up, 1 for down, 2 for idle, -1 for unused 
	if LiftPos.CurrentFloor < LiftPos.DestinationFloor {
		Direction = 0
	} else if LiftPos.CurrentFloor > LiftPos.DestinationFloor {
		Direction = 1
	} else {
		Direction = 2
	}
}



func CostFunction(NewOrder, LiftPos, ArrayOfOrders) {
	for lift = 0; lift < Lifts; lift++ {
		Switch (Direction(LiftPos[lift])) {
		case 0: //Up
		
		case 1: //Down

		case 2: //Idle
		}
			if SearchForBadOrder(NewOrder,ArrayOfOrders,lift) == false {
				cost =  abs(CurrentPos.Floor - NewOrder.Floor)
				for elemtents in ArrayOfOrders != -1 {
					cost += 3 
				}
			}
		}
	}
}


func SearchForBadOrder(NewOrder, ArrayOfOrders, lift) {
	var BadOrder = false
	for floor = 0; floor < NumberOfFloors; floor++ {
		for button = 0; button < 3; button++ {

			if NewOrder.Direction =! ArrayOfOrders.Direction[lift][floor][button]
				BadOrder = true
				break
			}
			else if NewOrder.Direction == up and NewOrder.Floor > ArrayOfOrders.floor[lift][floor][button]
				BadOrder = true	
				break

			else if NewOrder.Direction == down and NewOrder.Floor < ArrayOfOrders.floor[lift][floor][button]
				BadOrder = true	
				break
		}
		if BadOrder == true
		break
	}
	return BadOrder
}