package main

import (
	"driver"
	"time"
)

type order struct {
	floor int
	orderType int // 0 for up, 1 for down and 2 for command
	elevatorId int //Which elevator the order is from
}




func CheckUpButtons() {	
	for i := 0; i < 3; i++ {
		if driver.GetButtonSignal(0,i) == 1 {
			driver.SetButtonLampOn(0,i)
			//add new order somehow
		}
	}
}

func CheckDownButtons() {	
	for i := 1; i < 4; i++ {
		if driver.GetButtonSignal(1,i) == 1 {
			driver.SetButtonLampOn(1,i)
			//add new order somehow
		}
	}
}

func CheckCommandButtons() {
	for j := 0; j < 4; j++ {
		if driver.GetButtonSignal(2,j) == 1 {
			driver.SetButtonLampOn(2,j)
			//add new order somehow
		}
	}
}

func ButtonPoller() {
	for {
		CheckDownButtons()
		CheckUpButtons()
		CheckCommandButtons()
		time.Sleep(10)
	}
}


func main () {
	driver.Init()
	go ButtonPoller()
	neverReturn := make (chan int)
	<-neverReturn
}
				
		
	
