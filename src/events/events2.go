package main

import (
	"driver"
	"time"
	"fmt"
	"network"
	"udp"
	"cost"
	"types"
)

const Floors = 4
const Lifts = 3
var MYID = 0
var LightArray [3][4] int //row 0 for up, row 1 for down, row 2 for inside
var GlobalArrayOfOrders[][] types.Order
var OrderQueue[] types.Order
var LiftPos[Lifts] types.Position

func MakeGlobalArrayOfOrders() {
	GlobalArrayOfOrders = make([][]types.Order,Lifts)
	//LiftPos = make([]queue.Position,Lifts)
}

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
			//driver.SetButtonLampOn(0,i)
			network.SendMessage(send_ch,2,"",0,i,-1,MYID)	
			//OrderQueue = append(OrderQueue,queue.Order{DestinationFloor:i, ButtonType:0})
		}
	}
}

func CheckDownButtons(send_ch chan udp.Udp_message) {	
	for i := 1; i < 4; i++ {
		if driver.GetButtonSignal(1,i) == 1 && LightArray[1][i] == 0{
			LightArray[1][i] = 1
			//driver.SetButtonLampOn(1,i)
			network.SendMessage(send_ch,2,"",1,i,-1,MYID)	
			//OrderQueue = append(OrderQueue,queue.Order{DestinationFloor:i, ButtonType:1})	
		}
	}
}

func CheckCommandButtons(send_ch chan udp.Udp_message) {
	for i := 0; i < 4; i++ {
		if driver.GetButtonSignal(2,i) == 1 && LightArray[2][i] == 0 {
			LightArray[2][i] = 1			
			driver.SetButtonLampOn(2,i)
			OrderQueue = append(OrderQueue,types.Order{DestinationFloor:i, ButtonType:2})
		}
	}
}

func ButtonPoller(send_ch chan udp.Udp_message) {
	for {
		CheckDownButtons(send_ch)
		CheckUpButtons(send_ch)
		CheckCommandButtons(send_ch)
		time.Sleep(25*time.Millisecond)
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
	LiftPos[MYID].DestinationFloor = -1
}

func States(Task string) {
	switch Task {
	case "UP":
		driver.SetDirection(1)
	case "DOWN":
		driver.SetDirection(-1)
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

func RemoveOrder(flr int) {
	if (len(OrderQueue) > 1) {
		for i := len(OrderQueue)-1; i > -1; i-- {
			if (OrderQueue[i].DestinationFloor == flr) && (i == len(OrderQueue)-1) {
				LightArray[OrderQueue[i].ButtonType][OrderQueue[i].DestinationFloor] = 0	
				OrderQueue = OrderQueue[:i]
			} else if (OrderQueue[i].DestinationFloor == flr) && (i != len(OrderQueue)-1) {
				LightArray[OrderQueue[i].ButtonType][OrderQueue[i].DestinationFloor] = 0
				OrderQueue = append(OrderQueue[:i], OrderQueue[i+1:]...)
			}	
		}	
	} else {
		LightArray[OrderQueue[0].ButtonType][OrderQueue[0].DestinationFloor] = 0
		OrderQueue = OrderQueue[:0]
	}
	LightsOff()
}

func MasterReciever(MessageToProcess chan types.NetworkMessage, send_ch chan udp.Udp_message) {
	timeout := make(chan bool, 1)
	go func() {
    	time.Sleep(10 * time.Second)
    	timeout <- true
	}()
	for {	
		//select {
		Message := <- MessageToProcess
			switch Message.MessageType {
			case 1:
				fmt.Printf("I'm alive from: %v \n", Message.ElevatorID)
			case 2:
				NewOrder := types.Order{DestinationFloor:Message.DestinationFloor,ButtonType:Message.ButtonType}
				LiftToUse :=  cost.CostFunction(NewOrder, LiftPos)
				
				network.SendMessage(send_ch,4,"",Message.ButtonType,Message.DestinationFloor,-1,LiftToUse)
				time.Sleep(25*time.Millisecond)
				for i := 0; i < Lifts; i++ {
					network.SendMessage(send_ch,6,"",Message.ButtonType,Message.DestinationFloor,-1,MYID)
					time.Sleep(25*time.Millisecond)
				}
			case 3:
				LiftPos[Message.ElevatorID] = types.Position{Message.CurrentFloor,Message.DestinationFloor}
			case 4:
				if Message.ElevatorID == MYID {
					NewOrder := types.Order{ButtonType:Message.ButtonType,DestinationFloor:Message.DestinationFloor}
					OrderQueue = append(OrderQueue,NewOrder)	
				}	
			case 5:
				for i := 0; i < Lifts; i++ {
					network.SendMessage(send_ch,7,"",Message.ButtonType,Message.DestinationFloor,-1,MYID)
					time.Sleep(25*time.Millisecond)
				}	
			case 6:
				driver.SetButtonLampOn(Message.ButtonType,Message.DestinationFloor)
				LightArray[Message.ButtonType][Message.DestinationFloor] = 1
			case 7:
				LightArray[0][Message.DestinationFloor] = 0
				LightArray[1][Message.DestinationFloor] = 0
				LightsOff()
			}
		//case <- timeout:
		//	fmt.Printf("Timeout occured\n")
	}
}

func SlaveReciever(MessageToProcess chan types.NetworkMessage) {
	timeout := make(chan bool, 1)
	go func() {
    	time.Sleep(10 * time.Second)
    	timeout <- true
	}()
	for {	
	//	select {
			Message := <- MessageToProcess
				switch Message.MessageType {
				//case 1:
				//	fmt.Printf("I'm alive from: %v \n", Message.ElevatorID)
				case 4:
					if Message.ElevatorID == MYID {
						NewOrder := types.Order{ButtonType:Message.ButtonType,DestinationFloor:Message.DestinationFloor}
						OrderQueue = append(OrderQueue,NewOrder)	
					}	
				case 6:
					driver.SetButtonLampOn(Message.ButtonType,Message.DestinationFloor)
					LightArray[Message.ButtonType][Message.DestinationFloor] = 1
				case 7:
					LightArray[0][Message.DestinationFloor] = 0
					LightArray[1][Message.DestinationFloor] = 0
					LightsOff()
				}

			//se <- timeout:
			//t.Printf("Timeout occured\n")
	}
}


func Run(FloorReached chan int,send_ch chan udp.Udp_message) {
	InFloor := true
	go func() {
		for {
			flr := <- FloorReached
			//Network.SendMessage(send_ch,3,"",-1,LiftPos[MYID].DestinationFloor,flr,MYID)
			LiftPos[MYID].CurrentFloor = flr
			if (flr == LiftPos[MYID].DestinationFloor && len(OrderQueue) != 0) && InFloor == false{
				OrderQueue = cost.InternalCostFunction(OrderQueue,LiftPos[MYID])
				States("STOP")
				RemoveOrder(flr)
				network.SendMessage(send_ch,5,"",-1,flr,-1,MYID)
				InFloor = true
			} else if len(OrderQueue) != 0 && InFloor == false{
				OrderQueue = cost.InternalCostFunction(OrderQueue,LiftPos[MYID])
			}
			time.Sleep(250*time.Millisecond)
		}
	}()
	go func(){
		for{
			if len(OrderQueue) != 0  {
				LiftPos[MYID].DestinationFloor = OrderQueue[0].DestinationFloor
				//Network.SendMessage(send_ch,3,"",-1,LiftPos[MYID].CurrentFloor,OrderQueue[0].DestinationFloor,MYID)
				if InFloor == true {
					Dir := cost.FindDirection(LiftPos[MYID])
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
			time.Sleep(250*time.Millisecond)		
		}
	}()
}

func main () {
	
	InitializeLift()
	MakeGlobalArrayOfOrders()
	send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	MessageToProcess := make(chan network.NetworkMessage)
	FloorReached := make(chan int)
	err := udp.Udp_init(20006, 20006, 200, send_ch, receive_ch)	
	go network.ReadFromNetwork (receive_ch, MessageToProcess)
	go MasterReciever(MessageToProcess,send_ch)
	//go SlaveReciever(MessageToProcess)
	go ButtonPoller(send_ch)
	go FloorPoller(FloorReached)
	go Run(FloorReached, send_ch)
	for {
		if driver.GetStopSignal() == 1 {
			States("STOP")
			break
		}
		time.Sleep(100)
	}
		if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}


	//fmt.Println(ElevatorQueue)
	//queue.InternalCostFunction(OrderQueue, LiftPos)
	
	/*
	send_ch := make (chan udp.Udp_message)
	err := udp.Udp_init(20015, 20010, 200, send_ch, receive_ch)	
	go ButtonPoller(send_ch)
	go FloorPoller(FloorReached,send_ch)
	go TestRun(send_ch,FloorReached)
	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
	
	/*neverReturn := make (chan int)
	<-neverReturn*/
}
