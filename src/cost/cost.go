package cost


import (
	"fmt"
	"math/rand"
	"time"
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

func Abs(num int) int {
	if num >= 0 {
		return num
	} else {
		return -num
	}
}


func InternalCostFunction(OrderQueue[] Order, LiftPos Position) []Order{	
	Cost := make(map[int]int)
	for i := 0; i < len(OrderQueue); i++ {

		if (OrderQueue[i].ButtonType == FindDirection(LiftPos)) { //Hvis retningen på bestillingen og heisen er den samme
			if (OrderQueue[i].DestinationFloor >= LiftPos.CurrentFloor && OrderQueue[i].DestinationFloor <= LiftPos.DestinationFloor) && (OrderQueue[i].ButtonType != 2) { 
				Cost[i] += Abs(LiftPos.CurrentFloor - OrderQueue[i].DestinationFloor) //Hvis bestilling oppover er på veien				
				//fmt.Println("Passer oppover med kost ")
			}
			if (OrderQueue[i].DestinationFloor <= LiftPos.CurrentFloor && OrderQueue[i].DestinationFloor >= LiftPos.DestinationFloor) && (OrderQueue[i].ButtonType != 2) { 
				Cost[i] += Abs(LiftPos.CurrentFloor - OrderQueue[i].DestinationFloor) //Hvis bestilling nedover er på veien
				//fmt.Println("Passer nedover med kost ")
			} 
		} else if OrderQueue[i].ButtonType == 2 { //Hvis bestilling kommer innenifra
			Cost[i] += Abs(LiftPos.CurrentFloor - OrderQueue[i].DestinationFloor) 				
				//fmt.Println("Inneknapp med kost")
		} else {
			Cost[i] += 10
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
	
	//fmt.Println(Cost)

	//fmt.Println("Before =", OrderQueue)	SKRIVE BEDRE NAVN PÅ x!
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
	
	//fmt.Println("After =", OrderQueue)	
	return x

}


func CostFunction(NewOrder Order, LiftPos[3] Position) int {

	for lift := 0; lift < Lifts; lift++ { // Går igjennom alle heiskøene og sjekker om de kan ta noen på veien

		switch (FindDirection(LiftPos[lift])) { 

		case 0: //Up

			if (NewOrder.DestinationFloor > LiftPos[lift].CurrentFloor) && (NewOrder.DestinationFloor <= LiftPos[lift].DestinationFloor) && (NewOrder.ButtonType != 1) {
				fmt.Println("case 0")
				if lift == 0 {return lift} // Heis 1	
				if lift == 1 {return lift} // Heis 2
				if lift == 2 {return lift} // Heis 3
			}
	 
		case 1: //Down			

			if (NewOrder.DestinationFloor < LiftPos[lift].CurrentFloor) && (NewOrder.DestinationFloor >= LiftPos[lift].DestinationFloor) && (NewOrder.ButtonType != 0) {
				fmt.Println("case 1")				
				if lift == 0 {return lift} // Heis 1	
				if lift == 1 {return lift} // Heis 2
				if lift == 2 {return lift} // Heis 3
			}	 	
		}
	}
	fmt.Println("Random")
	rand := random(1, 4)
	fmt.Println(rand-1)
	return rand-1
}



func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

/*func main() {
	
	LiftPos := append(LiftPos, Position{CurrentFloor:1, DestinationFloor:3})
	LiftPos = append(LiftPos, Position{CurrentFloor:0, DestinationFloor:3})
	LiftPos = append(LiftPos, Position{CurrentFloor:3, DestinationFloor:1})

	NewOrder.ButtonType = 1
	NewOrder.DestinationFloor = 2
	
	
	CostFunction(NewOrder, LiftPos)

}*/
















