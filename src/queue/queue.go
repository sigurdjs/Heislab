package queue


import (
	"math"
)

const Floors = 4
const Lifts = 3


type Order struct {
	DestinationFloor int
	ButtonType int // 0 for up, 1 for down and 2 for command
	//elevatorId int //Which elevator the order is from
}

type Position struct {
	CurrentFloor int
	DestinationFloor int
}







func FindDirection(LiftPos Position) int {
	dir := -1 //0 for up, 1 for down, 2 for idle, -1 for unused 
	if LiftPos.CurrentFloor < LiftPos.DestinationFloor {
		dir = 0
	} else if LiftPos.CurrentFloor > LiftPos.DestinationFloor {
		dir = 1
	} else {
		dir = 2
	}
	return dir
}

func CostFunction(NewOrder Order, LiftPos[3] Position) int {
	Cost := []int{100000,100000,1000000}
	for lift := 0; lift < Lifts; lift++ {
		switch (FindDirection(LiftPos[lift])) {
		case 0: //Up
			if NewOrder.DestinationFloor > LiftPos[lift].CurrentFloor && NewOrder.DestinationFloor < LiftPos[lift].DestinationFloor {
				Cost[lift] = 3*(NewOrder.DestinationFloor - LiftPos[lift].CurrentFloor) }	 
		case 1: //Down			
			if NewOrder.DestinationFloor < LiftPos[lift].CurrentFloor && NewOrder.DestinationFloor > LiftPos[lift].DestinationFloor {
				Cost[lift] = 3*(LiftPos[lift].CurrentFloor - NewOrder.DestinationFloor) }	 
		case 2: //Idle
			temp := float64(LiftPos[lift].CurrentFloor - NewOrder.DestinationFloor)
			Cost[lift] = int(math.Abs(temp))
		}
	}
	MaxCost := 0 
//	MaxLift := -1
	for i := 0; i < Lifts; i++ {
		if Cost[i] > MaxCost {
			MaxCost = Cost[i]
			//MaxLift = i
		}	
	}
	return MaxCost
}



