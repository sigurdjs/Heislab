package main

import (
	"driver"
	"time"
	"fmt"
	"Network"
	"udp"
	"queue"
)


var dest int
var LightArray [3][4] int //row 0 for up, row 1 for down, row 2 for inside
var MasterArray[] queue.order
var LiftPos[3] queue.position


func CheckUpButtons(send_ch chan udp.Udp_message) {	
	for i := 0; i < 3; i++ {
		if driver.GetButtonSignal(0,i) == 1 && LightArray[0][i] == 0{
			LightArray[0][i] = 1
			driver.SetButtonLampOn(0,i)
			Network.SendNewOrderMessage(send_ch,1,0,i)		
		}
	}
}

func CheckDownButtons(send_ch chan udp.Udp_message) {	
	for i := 1; i < 4; i++ {
		if driver.GetButtonSignal(1,i) == 1 && LightArray[1][i] == 0{
			LightArray[1][i] = 1
			driver.SetButtonLampOn(1,i)			
			Network.SendNewOrderMessage(send_ch,1,1,i)
		}
	}
}

func CheckCommandButtons(send_ch chan udp.Udp_message, commandbutton chan int) {
	for i := 0; i < 4; i++ {
		if driver.GetButtonSignal(2,i) == 1 && LightArray[2][i] == 0 {
			LightArray[2][i] = 1			
			driver.SetButtonLampOn(2,i)
			Network.SendNewOrderMessage(send_ch,1,2,i)
		}
	}
}

func ButtonPoller(send_ch chan udp.Udp_message, commandbutton chan int) {
	for {
		CheckDownButtons(send_ch)
		CheckUpButtons(send_ch)
		CheckCommandButtons(send_ch,commandbutton)
		time.Sleep(10)
	}
}

func InitializeLift() {
	driver.Init()
	driver.SetDirection(1)
	if driver.GetFloor() == 0 {
		driver.SetDirection(0)
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
		time.Sleep(100)
	}		
}	
	





func MessageRecieved(MessageToProcess chan Network.NetworkMessage) {
	var newOrder queue.order
	timeout := make(chan bool, 1)
	go func() {
    	time.Sleep(1 * time.Second)
    	timeout <- true
	}()
	for {	
		select {
		case <- MessageToProcess:
			Message := <- MessageToProcess
			switch Message.MessageType {
			case 1:
				fmt.Printf("I'm alive from: %v \n", Message.AliveMessage)
			case 2:
				newOrder.DestinationFloor = Message.DestinationFloor
				newOrder.ButtonType = Message.ButtonType
				MasterArray = append(MasterArray,newOrder)
				cost := queue.CostFunction(MasterArray[0],LiftPos)
				fmt.Printf("Kosten for denne knapp er: %v \n",cost)
			case 3:
				LiftPos[Message.elevatorId] = queue.position {Message.currentFloor,Message.DestinationFloor}	
			}
		case <- timeout:
			fmt.Printf("Timeout occured\n")
		}
	}
}

func Start(commandbutton chan int) {
	for {
		dest = <- commandbutton
		if dest < driver.GetFloor() {
			driver.SetDirection(1)
		} else {
			driver.SetDirection(-1)
		}
	}
}

func Stop(flooreached chan int) {
	for {
		flr := <- flooreached
		if flr == dest {
			driver.SetDirection(0)
		}
	}
}






func PrintOrders(ch chan queue.order) {
	for {
		newOrder := <- ch 
		fmt.Printf("Floor: %v \n", newOrder.floor)	
		fmt.Printf("Button:%v \n", newOrder.orderType)
		time.Sleep(10)
	}
}


func main () {
	driver.Init()
	commandbutton := make(chan int)
	flooreached := make(chan int)
	/*send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	err := udp.Udp_init(20019, 20019, 100, send_ch, receive_ch)	*/
	go ButtonPoller(send_ch,commandbutton)
	go FloorPoller(floor)
	go Start(commandbutton)
	go Stop(flooreached)
	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
	neverReturn := make (chan int)
	<-neverReturn
}

		
	
