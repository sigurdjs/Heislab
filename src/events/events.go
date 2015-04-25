package main

import (
	"driver"
	"time"
	"fmt"
	"network"
	"backup"
	"udp"
	"cost"
	"types"
	"os"
	"os/signal"
)


var MYID = 2
var master = false
var LightArray [3][4] int //row 0 for up, row 1 for down, row 2 for inside
var OrderQueue[] types.Order
var LiftPos[] types.Position
var LiftStatus[] bool
var OffGrid bool

func MakeGlobalArrayOfOrders(Lifts int) {
	//GlobalArrayOfOrders := make([][]types.Order,Lifts)
	LiftPos = make([]types.Position,Lifts)
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
			if OffGrid == true {
				OrderQueue = append(OrderQueue,types.Order{DestinationFloor:i, ButtonType:0})
				driver.SetButtonLampOn(0,i)
			}
			network.SendMessage(send_ch,2,"",0,i,-1,MYID)	
		}
	}
}

func CheckDownButtons(send_ch chan udp.Udp_message) {	
	for i := 1; i < 4; i++ {
		if driver.GetButtonSignal(1,i) == 1 && LightArray[1][i] == 0{
			LightArray[1][i] = 1
			if OffGrid == true {
				OrderQueue = append(OrderQueue,types.Order{DestinationFloor:i, ButtonType:0})
				driver.SetButtonLampOn(0,i)
			}
			network.SendMessage(send_ch,2,"",1,i,-1,MYID)	
		}
	}
}

func CheckCommandButtons(send_ch chan udp.Udp_message) {
	for i := 0; i < 4; i++ {
		if driver.GetButtonSignal(2,i) == 1 && LightArray[2][i] == 0 {
			LightArray[2][i] = 1			
			driver.SetButtonLampOn(2,i)
			OrderQueue = append(OrderQueue,types.Order{DestinationFloor:i, ButtonType:2})
			//network.SendMessage(send_ch,1,"et bajs",-1,-1,-1,-1)
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
		if driver.GetFloor() == 0 || driver.GetFloor() == 1 || driver.GetFloor() == 2 || driver.GetFloor() == 3 {
			driver.SetDirection(0)
			break
		}
		time.Sleep(100*time.Millisecond)
	}
	OrderQueue = backup.ReadBackup() // <---------------------------------------------------------------------------------------------------------------------------------------------------------
	if len(OrderQueue) <= 0 { // <----------------------------------------------------------------------------------------------------------------------------------------------------------------
		LiftPos[MYID].DestinationFloor = -1 // <----------------------------------------------------------------------------------------------------------------------------------------------
	} else { // <---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		for i := 0; i < len(OrderQueue); i++ { // <-------------------------------------------------------------------------------------------------------------------------------------------
			driver.SetButtonLampOn(OrderQueue[i].ButtonType,OrderQueue[i].DestinationFloor) // <------------------------------------------------------------------------------------------
			LightArray[OrderQueue[i].ButtonType][OrderQueue[i].DestinationFloor] = 1 
		}
	}
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
		//fmt.Println(LiftStatus)
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
	backup.WriteBackup(OrderQueue)
}

func MasterReciever(MessageToProcess chan types.NetworkMessage, send_ch chan udp.Udp_message,TimedOut chan int,Lifts int) {
	var LiftGone = -1
	var numberOfLifts = Lifts
	var Message types.NetworkMessage
	go func() {
		for {	
			if master == true {
				Message = <- MessageToProcess
					switch Message.MessageType {
					case 1:
						LiftStatus[Message.ElevatorID] = true
						//fmt.Println("Alive from lift number: \n", Message.ElevatorID)
						//fmt.Println(LiftStatus)
					case 2:
						NewOrder := types.Order{DestinationFloor:Message.DestinationFloor,ButtonType:Message.ButtonType}
						LiftToUse :=  cost.CostFunction(NewOrder, LiftPos,LiftGone ,numberOfLifts)
						fmt.Println("Lift gone is:", LiftGone)
						fmt.Println("Number of lifts:",numberOfLifts)
						network.SendMessage(send_ch,4,"",Message.ButtonType,Message.DestinationFloor,-1,LiftToUse)  
						for i := 0; i < Lifts; i++ {
							network.SendMessage(send_ch,6,"",Message.ButtonType,Message.DestinationFloor,-1,MYID)
						}
					case 3:
						LiftPos[Message.ElevatorID] = types.Position{Message.CurrentFloor,Message.DestinationFloor}
					case 4:
						if Message.ElevatorID == MYID {
							NewOrder := types.Order{ButtonType:Message.ButtonType,DestinationFloor:Message.DestinationFloor}
							OrderQueue = append(OrderQueue,NewOrder)	
							backup.WriteBackup(OrderQueue)
						}
					case 5:
						for i := 0; i < Lifts; i++ {
							network.SendMessage(send_ch,7,"",Message.ButtonType,Message.DestinationFloor,-1,MYID)
						}	
					case 6:
						driver.SetButtonLampOn(Message.ButtonType,Message.DestinationFloor)
						LightArray[Message.ButtonType][Message.DestinationFloor] = 1
					case 7:
						LightArray[0][Message.DestinationFloor] = 0
						LightArray[1][Message.DestinationFloor] = 0
						LightsOff()
				}
			} else {
				Message = <- MessageToProcess
					switch Message.MessageType {
					case 1:
						LiftStatus[Message.ElevatorID] = true
						//fmt.Println(LiftStatus)
						//fmt.Println("Alive from lift number: \n", Message.ElevatorID)
					case 4:
						if Message.ElevatorID == MYID {
							NewOrder := types.Order{ButtonType:Message.ButtonType,DestinationFloor:Message.DestinationFloor}
							OrderQueue = append(OrderQueue,NewOrder)	
							backup.WriteBackup(OrderQueue)
						}
					case 6:
						driver.SetButtonLampOn(Message.ButtonType,Message.DestinationFloor)
						LightArray[Message.ButtonType][Message.DestinationFloor] = 1
					case 7:
						LightArray[0][Message.DestinationFloor] = 0
						LightArray[1][Message.DestinationFloor] = 0
						LightsOff()	
					}
				}
			}
		}()
	go func() {
		for {
			LiftGone = <- TimedOut
			if LiftGone != -1 {
				numberOfLifts = Lifts -1
			} else {
				numberOfLifts = Lifts
			}
		}
	}()
}

/*func SlaveReciever(MessageToProcess chan types.NetworkMessage) {
	for {	
		//select {
			Message := <- MessageToProcess
				switch Message.MessageType {
				case 1:
					LiftStatus[Message.ElevatorID] = true
					//fmt.Println("Alive from lift number: \n", Message.ElevatorID)
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

		//	case <- timeout:
		//		t.Printf("Timeout occured\n")
	}
}*/


func Run(FloorReached chan int,send_ch chan udp.Udp_message) {
	InFloor := true
	go func() {
		for {
			flr := <- FloorReached
			LiftPos[MYID].CurrentFloor = flr
			if (flr == LiftPos[MYID].DestinationFloor && len(OrderQueue) != 0) && InFloor == false{
				OrderQueue = cost.InternalCostFunction(OrderQueue,LiftPos[MYID])
				backup.WriteBackup(OrderQueue)
				States("STOP")
				RemoveOrder(flr)
				network.SendMessage(send_ch,5,"",-1,flr,-1,MYID)
				InFloor = true
			} else if len(OrderQueue) != 0 && InFloor == false{
				OrderQueue = cost.InternalCostFunction(OrderQueue,LiftPos[MYID])
				backup.WriteBackup(OrderQueue)
			}
		}
	}()
	go func(){
		for{
			network.SendMessage(send_ch,1,"Alive",-1,-1,-1,MYID)
			if len(OrderQueue) != 0  {
				LiftPos[MYID].DestinationFloor = OrderQueue[0].DestinationFloor
				network.SendMessage(send_ch,3,"",-1,OrderQueue[0].DestinationFloor,LiftPos[MYID].CurrentFloor,MYID)
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
			time.Sleep(100*time.Millisecond)		
		}
	}()
}




func TimeOut(TimedOut chan int){
	timeout := make(chan bool, 1)
	notimeout := make(chan bool, 1)
	var LiftGone int
	var numberOfGone int
	time.Sleep(10*time.Second)
	go func() {
		for {
	 		time.Sleep(3 * time.Second)
	 		timeout <- true
	 	}
	}()
	go func() {
		for {
			if (LiftStatus[0] == true && LiftStatus[1] == true) && (LiftStatus[2] == true) {
				notimeout <- true
			}
			time.Sleep(1000*time.Millisecond)
		}	
	}()
	go func() {
		for {
			select {
			case <-timeout:
				fmt.Println(LiftStatus)
				for i := 0; i < len(LiftStatus); i++ {
					if LiftStatus[i] == false {
						LiftGone = i
					}
				}
				if LiftStatus[0] == false && MYID == 1 {
					master = true
					if LiftStatus[2] == false {
						OffGrid = true
					}
					fmt.Println("Master is 1")
				} else if LiftStatus[0] == false && LiftStatus[1] == false {
					master = true
					OffGrid = true
					fmt.Println("Master is 2")
				} else if MYID == 0 {
					master = true
					if (LiftStatus[1] == false && LiftStatus[2] == false) {
						OffGrid = true
					}
				} else {
					master = false
				} 
				TimedOut <- LiftGone
				TimedOut <- numberOfGone
			case <- notimeout:
				for i := 0; i < len(LiftStatus); i++ {
					LiftStatus[i] = false
				}
				TimedOut <- -1
			}
		}
	}()
}





const Lifts = 3
const Floors = 4

func main () {
	LiftPos = make([]types.Position,Lifts)
	LiftStatus = make([]bool, Lifts)
	InitializeLift()
	send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	MessageToProcess := make(chan types.NetworkMessage)
	FloorReached := make(chan int)
	Timeout := make(chan int,1)
	err:= udp.Udp_init(20015, 20015, 130, send_ch, receive_ch)
	go network.ReadFromNetwork (receive_ch,MessageToProcess)	
	go MasterReciever(MessageToProcess,send_ch,Timeout,Lifts)
	go ButtonPoller(send_ch)
	go FloorPoller(FloorReached)
	go Run(FloorReached, send_ch)
	go TimeOut(Timeout)
	
	go func(){
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
    	<- sigchan 
        driver.SetDirection(0)
        os.Exit(0)
	}()
	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
	neverReturn := make (chan int)
	<-neverReturn
}
