package cost


import (
	"fmt"
	"math/rand"
	"time"
	"types"
)





func FindDirection(LiftPos types.Position) int {
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


func InternalCostFunction(OrderQueue[] types.Order, LiftPos types.Position) []types.Order{	
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
	var BestCost[] types.Order
	MinCostValue := OrderQueue[MinCostPosition]
	BestCost = append(BestCost,MinCostValue)
	//fmt.Println(BestCost)
	if MinCostPosition == len(OrderQueue)-1 {
		OrderQueue = OrderQueue[:MinCostPosition]
	} else {
		OrderQueue = append(OrderQueue[MinCostPosition:],OrderQueue[:MinCostPosition+1]...)
	}
	BestCost = append(BestCost,OrderQueue...)
	
	//fmt.Println("After =", OrderQueue)	
	return BestCost

}


func CostFunction(NewOrder types.Order, LiftPos[] types.Position, TimedOut int ,Lifts int) int {

	for lift := 0; lift < Lifts; lift++ { // Går igjennom alle heiskøene og sjekker om de kan ta noen på veien

		switch (FindDirection(LiftPos[lift])) { 

		case 0: //Up

			if (NewOrder.DestinationFloor > LiftPos[lift].CurrentFloor) && (NewOrder.DestinationFloor <= LiftPos[lift].DestinationFloor) && (NewOrder.ButtonType != 1) {
				fmt.Println("case 0")
				for i := 0; i < Lifts; i++ {
					if TimedOut == 0 {
						if lift == i {return lift +1} // Heis n
					} else if TimedOut == 1 {
						if lift == i {return lift*2}
					} else {
						if lift == i {return lift}
					}
				}
			}
	 
		case 1: //Down			

			if (NewOrder.DestinationFloor < LiftPos[lift].CurrentFloor) && (NewOrder.DestinationFloor >= LiftPos[lift].DestinationFloor) && (NewOrder.ButtonType != 0) {
				fmt.Println("case 1")				
				for i := 0; i < Lifts; i++ {
					if TimedOut == 0 {
						if lift == i {return lift +1} // Heis n
					} else if TimedOut == 1 {
						if lift == i {return lift*2}
					} else {
						if lift == i {return lift}
					}
				}
			}	

		case 2: //Idle
			if (NewOrder.DestinationFloor == LiftPos[lift].CurrentFloor) {
				fmt.Println("case 2")				
				for i := 0; i < Lifts; i++ {
					if TimedOut == 0 {
						if lift == i {return lift +1} // Heis n
					} else if TimedOut == 1 {
						if lift == i {return lift*2}
					} else {
						if lift == i {return lift}
					}
				}
			} 	
		}
	}
	var rand int
	fmt.Println("Random", Lifts)
	if TimedOut == 0 {
		rand = random(0, Lifts)+1
	} else if TimedOut == 1 {
		rand = 2*random(0, Lifts)
	} else {
		rand = random(0, Lifts)
	}
	fmt.Println(rand)
	return rand
}



func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}
