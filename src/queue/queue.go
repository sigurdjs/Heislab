package queue

import math

Floors = 4
Lifts = 3


type order struct {
	DestinationFloor int
	OrderType int // 0 for up, 1 for down and 2 for command
	//elevatorId int //Which elevator the order is from
}

type position struct {
	CurrentFloor int
	DestinationFloor int
}

var OrdersToComplete[Floors*3] order

//Only for master
var LiftPos[Lifts] position
var MasterArrayOfOrders[Floors*2] order

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

func CostFunction(NewOrder order, LiftPos[] order) int {
	var Cost[0:Lifts] = 1000000 
	for lift := 0; lift < Lifts; lift++ {
		Switch (Direction(LiftPos[lift])) {
		case 0: //Up
			if NewOrder.DestinationFloor > LiftPos[lift].CurrentFloor && NewOrder.DestinationFloor < LiftPos[lift].DestinationFloor {
				Cost[lift] = 3*(NewOrder.DestinationFloor - LiftPos[lift].CurrentFloor) }	 
		case 1: //Down			
			if NewOrder.DestinationFloor < LiftPos[lift].CurrentFloor && NewOrder.DestinationFloor > LiftPos[lift].DestinationFloor {
				Cost[lift] = 3*(LiftPos[lift].CurrentFloor - NewOrder.DestinationFloor) }	 
		case 2: //Idle
			Cost[lift] = math.abs(LiftPos[lift].CurrentFloor - NewOrder.DestinationFloor)
		}
	}
	var MaxCost = 0, MaxLift = -1
	for i := O; i < Lifts i++ {
		if Cost[i] > MaxCost
			MaxCost = Cost[i]
			MaxLift = i
		}	
	}
	return MaxLift
}



