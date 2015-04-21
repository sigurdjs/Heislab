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
var MasterArray[] queue.Order
var LiftPos[] queue.Position


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
			Network.SendCurrentFloor(send_ch,1,4,1)
			time.Sleep(100)
			Network.SendNewOrderMessage(send_ch,1,2,i)
			commandbutton <- i
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

func FloorPoller(flooreached chan int) {
	for {
		currentFloor := driver.GetFloor()
		switch  currentFloor {
		case 0:
			driver.SetFloorLamp(0)
			flooreached <- 0
		case 1:
			driver.SetFloorLamp(1)
			flooreached <- 1
		case 2:
			driver.SetFloorLamp(2)
			flooreached <- 2
		case 3:
			driver.SetFloorLamp(3)
			flooreached <- 3
		}		
		time.Sleep(100)
	}		
}	
	
func MessageRecieved(MessageToProcess chan Network.NetworkMessage) {
	var newOrder queue.Order
	/*timeout := make(chan bool, 1)
	go func() {
    		time.Sleep(1 * time.Second)
    		timeout <- true
	}()*/
	for {	
		//select {
		//case <- MessageToProcess:
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
			LiftPos[Message.ElevatorID] = queue.Position {Message.CurrentFloor,Message.DestinationFloor}	
		}
		//case <- timeout:
		//	fmt.Printf("Timeout occured\n")
	}
}

func Start(commandbutton chan int) {
	for {
		dest = <- commandbutton
		if dest < driver.GetFloor() {
			driver.SetDirection(-1)
		} else {
			driver.SetDirection(1)
		}
		time.Sleep(100)
	}
}

func Stop(flooreached chan int) {
	for {
		flr := <- flooreached
		if flr == dest {
			driver.SetDirection(0)
		}
		time.Sleep(100)
	}
}






/*func PrintOrders(ch chan queue.Order) {
	for {
		newOrder := <- ch 
		fmt.Printf("Floor: %v \n", newOrder.Floor)	
		fmt.Printf("Button:%v \n", newOrder.OrderType)
		time.Sleep(10)
	}
}*/


func main () {
	driver.Init()
	driver.SetDirection(0)
	commandbutton := make(chan int)
	flooreached := make(chan int)
	send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	MessageToProcess := make (chan Network.NetworkMessage)
	
	err := udp.Udp_init(20025, 20025, 200, send_ch, receive_ch)
	go Network.ReadFromNetwork(receive_ch, MessageToProcess)
	go ButtonPoller(send_ch,commandbutton)
	go FloorPoller(flooreached)
	go MessageRecieved(MessageToProcess)
	go Start(commandbutton)
	go Stop(flooreached)
	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
	neverReturn := make (chan int)
	<-neverReturn
}

		
	
