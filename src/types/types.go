package types

type NetworkMessage struct {
	MessageType int //1 for alive, 2 for order to master, 3 for floor update to master, 4 for order from master, 5 for reciept, 6 for set lights on, 7 for set lights off
	AliveMessage string //Message just to send something	
	ButtonType int //0 for up, 1 for down, 2 for inside, -1 for unused
	DestinationFloor int //-1 for unused
	CurrentFloor int //Current floor or destination floor depending on message type, also -1 for unused
	ElevatorID int //Lift number 1,2,3..n
}


type Order struct {
	DestinationFloor int
	ButtonType int // 0 for up, 1 for down and 2 for command
}

type Position struct {
	CurrentFloor int
	DestinationFloor int
}

