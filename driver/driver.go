package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "C/io.c"
#include "C/elev.c"
*/
import "C"

func Init(){
	C.elev_init()
} 

func SetDirection(direction int) {
	C.elev_set_motor_direction(direction)
}

func SetDoorOpen(value int) {
	C.elev_set_door_open_lamp(value)
}

func SetStopLamp(value int) {
	C.elev_set_stop_lamp(value)
}

func SetFloorLamp(floor int) {
	C.elev_set_floor_indicator(floor)
}

func GetFloor() int {
	return (C.elev_get_floow_sensor_signal())
}

func GetObstruction() int {
	return (C.elev_get_obstruction())
}

func GetStopSignal() int {
	return (C.elev_get_stop_signal())
}

