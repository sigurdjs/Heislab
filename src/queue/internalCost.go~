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



// Lager et midlertidig array
var ElevatorQueue = [][]int {{1,2},{0,1},{2,4},{1,4}}

func InternalCostFunction(ElevatorQueue [][]int, LiftPos Position int) {
	
	
	// Sjekker kosten
	var Cost = make(map[int]int)
	for i := 0; i < len(ElevatorQueue); i++ {
		
		if (FindDirection == 0 && ElevatorQueue[i][0] == 1) || (FindDirection == 1 && ElevatorQueue[i][0] == 0 {				// feil retning i forhold til kjøring		
			Cost[i] += 10 }
		
		if (FindDirection == 0 && (ElevatorQueue[i][1] < LiftPos.CurrentFloor || ElevatorQueue[i][1] > LiftPos.DestinationFloor)) || 
			(FindDirection == 1 && (ElevatorQueue[i][1] > LiftPos.CurrentFloor || ElevatorQueue[i][1] < LiftPos.DestinationFloor)) {	// ikke mellom CF og DF	
			Cost[i] += 5 }			

		if ElevatorQueue[i][1]  {														// ytre knapper
			Cost[i] += 2 }
		
		if ElevatorQueue[i][0] == 2 {														// indre knapper
			Cost[i] += 1 }

		Cost[i] += 0	

		//time.Sleep(10)	
	}
	MinCostPosition := 0
	MinCost := Cost[0]
	for j := 1; j < len(ElevatorQueue); j++ {
		if Cost[j] < MinCost {
			MinCost = Cost[j]
			MinCostPosition = j
		}
	} 
	
	fmt.Println(MinCostPosition)

	fmt.Println("Before =", ElevatorQueue)	

	SortArray(ElevatorQueue,MinCostPosition)
	
	fmt.Println("After =", ElevatorQueue)	

}

func SortArray(X [][]int, Position int) {
	temp := X[Position]
	X[Position] = X[Position+1] 
	for i := len(X)-1; i > 0; i-- {
		X[i] = X[i-1] 
	}
	X[0] = temp
}






















