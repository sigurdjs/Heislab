package main

import (
	"driver"
	"time"
	"fmt"
	"Network"
	"udp"
	"queue"
)


var Dst int
var LightArray [3][4] int //row 0 for up, row 1 for down, row 2 for inside
var MasterArray[] queue.Order
var OrderQueue[] queue.Order
var LiftPos queue.Position

func LightsOff() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if i == 0 && j == 3 {
			} else if i == 1 && j == 0 {
			} else {
				if LightArray[i][j] == 0 {
					driver.SetButtonLampOff(i,j)
				}
			}
		}
	}
}

func CheckUpButtons(send_ch chan udp.Udp_message) {	
	for i := 0; i < 3; i++ {
		if driver.GetButtonSignal(0,i) == 1 && LightArray[0][i] == 0{
			LightArray[0][i] = 1
			driver.SetButtonLampOn(0,i)
			//Network.SendNewOrderMessage(send_ch,1,0,i) 	
			OrderQueue = append(OrderQueue,queue.Order{DestinationFloor:i, ButtonType:0})
		}
	}
}

func CheckDownButtons(send_ch chan udp.Udp_message) {	
	for i := 1; i < 4; i++ {
		if driver.GetButtonSignal(1,i) == 1 && LightArray[1][i] == 0{
			LightArray[1][i] = 1
			driver.SetButtonLampOn(1,i)
			//Network.SendNewOrderMessage(send_ch,1,1,i)	
			OrderQueue = append(OrderQueue,queue.Order{DestinationFloor:i, ButtonType:1})	
		}
	}
}

func CheckCommandButtons(send_ch chan udp.Udp_message) {
	for i := 0; i < 4; i++ {
		if driver.GetButtonSignal(2,i) == 1 && LightArray[2][i] == 0 {
			LightArray[2][i] = 1			
			driver.SetButtonLampOn(2,i)
			//Network.SendNewOrderMessage(send_ch,1,2,i)
			OrderQueue = append(OrderQueue,queue.Order{DestinationFloor:i, ButtonType:2})
		}
	}
}

func ButtonPoller(send_ch chan udp.Udp_message) {
	for {
		CheckDownButtons(send_ch)
		CheckUpButtons(send_ch)
		CheckCommandButtons(send_ch)
		time.Sleep(10*time.Millisecond)
	}
}

func InitializeLift() {
	driver.Init()
	driver.SetDirection(-1)
	for {
		if driver.GetFloor() == 0 {
			driver.SetDirection(0)
			break
		}
		time.Sleep(100*time.Millisecond)
	}
	LiftPos.DestinationFloor = -1
}

func States(Task string) {
	switch Task {
	case "UP":
		driver.SetDirection(1)
		Dst = 3
	case "DOWN":
		driver.SetDirection(-1)
		Dst = 0
	case "STOP":
		driver.SetDirection(0)
		driver.SetDoorOpen(1)
		time.Sleep(3*time.Second)
		driver.SetDoorOpen(0)
	}
}



func FloorPoller(FloorReached chan int) {
	for {
		//fmt.Println(OrderQueue)
		currentFloor := driver.GetFloor()
		switch  currentFloor {
		case 0:
			driver.SetFloorLamp(0)
			FloorReached <- 0
		case 1:
			driver.SetFloorLamp(1)
			FloorReached <- 1
		case 2:
			driver.SetFloorLamp(2)
			FloorReached <- 2
		case 3:
			driver.SetFloorLamp(3)
			FloorReached <- 3
		}		
		time.Sleep(100*time.Millisecond)
	}		
}	
	





/*func MessageRecieved(MessageToProcess chan Network.NetworkMessage) {
	timeout := make(chan bool, 1)
	go func() {
    	time.Sleep(10 * time.Second)
    	timeout <- true
	}()
	for {	
		select {
		case Message := <- MessageToProcess:
			switch Message.MessageType {
			case 1:
				fmt.Printf("I'm alive from: %v \n", Message.ElevatorID)
			case 2:
				//MasterArray = append(MasterArray,newOrder)
				//fmt.Println(MasterArray)
				//cost := queue.CostFunction(MasterArray[0],LiftPos)
				//fmt.Printf("Kosten for denne knapp er: %v \n",cost)
				ElevatorQueue = append(ElevatorQueue,queue.Order{DestinationFloor:Message.DestinationFloor,ButtonType:Message.ButtonType})
			case 3:
				LiftPos[Message.ElevatorID] = queue.Position {Message.CurrentFloor,Message.DestinationFloor}	
			}
		case <- timeout:
			fmt.Printf("Timeout occured\n")
		}
	}
}*/


func TestRun(send_ch chan udp.Udp_message,FloorReached chan int) {
	States("UP")
	for {
		flr := <- FloorReached
		Network.SendCurrentFloor(send_ch,1,Dst,flr)
		if flr == 3 {
			States("DOWN")
		} else if flr == 0 {
			States("UP")
		}
	}
}

func RemoveOrder(flr int) {
	if (len(OrderQueue) > 1) {
		for i := len(OrderQueue)-1; i > -1; i-- {
			if (OrderQueue[i].DestinationFloor == OrderQueue[0].DestinationFloor) && (i == len(OrderQueue)-1) {
				LightArray[OrderQueue[i].ButtonType][OrderQueue[i].DestinationFloor] = 0	
				OrderQueue = OrderQueue[:i]
			} else if (OrderQueue[i].DestinationFloor == OrderQueue[0].DestinationFloor) && (i != len(OrderQueue)-1) {
				LightArray[OrderQueue[i].ButtonType][OrderQueue[i].DestinationFloor] = 0
				OrderQueue = append(OrderQueue[:i], OrderQueue[i+1:]...)
			}	
		}	
	} else {
		LightArray[OrderQueue[0].ButtonType][OrderQueue[0].DestinationFloor] = 0
		OrderQueue = OrderQueue[:0]
	}
	LightsOff()



		/*if OrderQueue[len(OrderQueue)-1].DestinationFloor == OrderQueue[0].DestinationFloor {
			if len(OrderQueue) >= 2 {
				OrderQueue = OrderQueue[:len(OrderQueue)-2]}
		} 	else {
				OrderQueue = OrderQueue[:0]
		}
	} else if len(OrderQueue) == 2 {
		if OrderQueue[1].DestinationFloor == OrderQueue[0].DestinationFloor {
			OrderQueue = OrderQueue[:0]
		} else {
			OrderQueue = OrderQueue[1:]
		}

	} else if len(OrderQueue) == 1 {
		OrderQueue = OrderQueue[:0]
	}*/
}

func Run(FloorReached chan int) {
	InFloor := true
	go func() {
		for {
			flr := <- FloorReached
			LiftPos.CurrentFloor = flr
			if (flr == LiftPos.DestinationFloor && len(OrderQueue) != 0) && InFloor == false{
				OrderQueue = queue.InternalCostFunction(OrderQueue,LiftPos)
				States("STOP")
				RemoveOrder(flr)
				InFloor = true
			} else if len(OrderQueue) != 0 && InFloor == false{
				OrderQueue = queue.InternalCostFunction(OrderQueue,LiftPos)
			}
			time.Sleep(25*time.Millisecond)
		}
	}()
	go func(){
		for{
			if len(OrderQueue) != 0  {
				LiftPos.DestinationFloor = OrderQueue[0].DestinationFloor
				if InFloor == true {
					Dir := queue.FindDirection(LiftPos)
					if Dir == 0 {
						States("UP")
						InFloor = false
					} else if Dir == 1 {
						States("DOWN")
						InFloor = false
					} else if Dir == 2 {
						InFloor = false
					}
				}
			}
			time.Sleep(25*time.Millisecond)		
		}
	}()
}



func PrintOrders(ch chan queue.Order) {
	
}

func PrintMessage(ch chan udp.Udp_message) {
	for {
		newOrder := <- ch 
		fmt.Printf("Length: %v \n", len(newOrder.Data))
		time.Sleep(10)
	}
}

func main () {
	
	InitializeLift()
	send_ch := make (chan udp.Udp_message)
	FloorReached := make(chan int)
	go ButtonPoller(send_ch)
	go FloorPoller(FloorReached)
	go Run(FloorReached)
	for {
		if driver.GetStopSignal() == 1 {
			States("STOP")
			break
		}
		time.Sleep(100)
	}

	//fmt.Println(ElevatorQueue)
	//queue.InternalCostFunction(OrderQueue, LiftPos)
	
	/*
	send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	MessageToProcess := make(chan Network.NetworkMessage)
	err := udp.Udp_init(20015, 20010, 200, send_ch, receive_ch)	
	go ButtonPoller(send_ch)
	go Network.ReadFromNetwork (receive_ch, MessageToProcess)
	go MessageRecieved(MessageToProcess)
	go FloorPoller(FloorReached,send_ch)
	go TestRun(send_ch,FloorReached)
	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
	
	/*neverReturn := make (chan int)
	<-neverReturn*/
}

		
	
