package queue


import (
	"math"
	"fmt"
//	"time"
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

func CostFunction(NewOrder Order, LiftPos[] Position) int {
	Cost := []int{100000,100000,100000}
	for lift := 0; lift < Lifts; lift++ {
		switch (FindDirection(LiftPos[lift])) {
		case 0: //Up
			if NewOrder.DestinationFloor > LiftPos[lift].CurrentFloor && NewOrder.DestinationFloor <= LiftPos[lift].DestinationFloor {
				Cost[lift] = 3*(NewOrder.DestinationFloor - LiftPos[lift].CurrentFloor) }	 
		case 1: //Down			
			if NewOrder.DestinationFloor < LiftPos[lift].CurrentFloor && NewOrder.DestinationFloor >= LiftPos[lift].DestinationFloor {
				Cost[lift] = 3*(LiftPos[lift].CurrentFloor - NewOrder.DestinationFloor) }	 
		case 2: //Idle
			temp := float64(LiftPos[lift].CurrentFloor - NewOrder.DestinationFloor)
			Cost[lift] = int(math.Abs(temp))
		}
	}
	MinCost := 100000 
	for i := 0; i < Lifts; i++ {
		if Cost[i] < MinCost {
			MinCost = Cost[i]
		}	
	}
	return MinCost
}



// Lager et midlertidig array
func InternalCostFunction(OrderQueue[] Order, LiftPos Position) []Order{	
	/*ElevatorQueue := make([][]int, len(OrderQueue))
	for a:= 0; a < len(OrderQueue); a++ {
		ElevatorQueue[a] = make([]int, 2)
	}
	for k := 0; k < len(OrderQueue); k++{
		ElevatorQueue[k][0] = OrderQueue[k].ButtonType
		ElevatorQueue[k][1] = OrderQueue[k].DestinationFloor
	}*/
	// Sjekker kosten
	Cost := make(map[int]int)
	for i := 0; i < len(OrderQueue); i++ {
		if (FindDirection(LiftPos) == 0 && OrderQueue[i].ButtonType == 1) || (FindDirection(LiftPos) == 1 && OrderQueue[i].ButtonType == 0) { // feil retning i forhold til kjøring		
			if (FindDirection(LiftPos) == 0 && OrderQueue[i].DestinationFloor == 3) {
				Cost[i] += 3
			}
			if (FindDirection(LiftPos) == 1 && OrderQueue[i].DestinationFloor == 0) {
			Cost[i] += 3
			} else {
				Cost[i] += 10} 
			}	
		if (FindDirection(LiftPos) == 0 && (OrderQueue[i].DestinationFloor < LiftPos.CurrentFloor || OrderQueue[i].DestinationFloor > LiftPos.DestinationFloor)) || 
			(FindDirection(LiftPos) == 1 && (OrderQueue[i].DestinationFloor > LiftPos.CurrentFloor || OrderQueue[i].DestinationFloor < LiftPos.DestinationFloor)) {	// ikke mellom CF og DF	
			Cost[i] += 5 }			

		if OrderQueue[i].ButtonType == 0 || OrderQueue[i].ButtonType == 1   {														// ytre knapper
			Cost[i] += 2 }
			
		if OrderQueue[i].ButtonType == 2 {														// indre knapper
			Cost[i] += 1 }

		Cost[i] += 0	
		//time.Sleep(10)
		for j := LiftPos.CurrentFloor; j < OrderQueue[i].DestinationFloor; j++ {
			Cost[i] += 2
		}	
	}
	MinCostPosition := 0
	MinCost := Cost[0]
	for j := 1; j < len(OrderQueue); j++ {
		if Cost[j] < MinCost {
			MinCost = Cost[j]
			MinCostPosition = j
		}
	} 
	
	fmt.Println(Cost)

	//fmt.Println("Before =", OrderQueue)	
	var x[] Order
	a := OrderQueue[MinCostPosition]
	x = append(x,a)
	//fmt.Println(x)
	if MinCostPosition == len(OrderQueue)-1 {
		OrderQueue = OrderQueue[:MinCostPosition]
	} else {
		OrderQueue = append(OrderQueue[MinCostPosition:],OrderQueue[:MinCostPosition+1]...)
	}
	x = append(x,OrderQueue...)
	//SortArray(OrderQueue,MinCostPosition)
	
	//fmt.Println("After =", OrderQueue)	
	return x

}
/*
func SortArray(ElevatorQueue[] Order, Position int) []Order{
	temp := ElevatorQueue[Position]
	//Array[Position] = Array[Position+1] 
	for i := Position; i > 0; i-- {
		ElevatorQueue[i] = ElevatorQueue[i-1] 
	}
	ElevatorQueue[0] = temp
	return ElevatorQueue
}*/




















