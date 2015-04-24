package main

import "driver"
import "time"



func main() {
	var floor int
	driver.Init()
	driver.SetDirection(1)
	for {
		floor = driver.GetFloor() 
		if floor != -1 {
			driver.SetFloorLamp(floor)		
		}
		if floor == 3 {
			driver.SetButtonLampOn(2,1)
			driver.SetDirection(-1)
		} else if floor == 0 {
			driver.SetButtonLampOff(2,1)
			driver.SetDirection(1)
		} else if driver.GetStopSignal() == 1 {
			driver.SetDirection(0)
			break
		}
		time.Sleep(100)
	}
}
		
		

	
