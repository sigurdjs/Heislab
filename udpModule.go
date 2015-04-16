package main


import (
	"fmt"
	"udp"
	"encoding/json"
	"time"
)

var masterIsAlive bool
var networkMessage string

type AliveMessage struct {
	AliveMessage string //Message just to send something
	ElevatorID int //Lift number 1,2,3..n
}

type NewOrderMessage struct {
	ElevatorID int //Lift number 1,2,3..n
	OrderType int //0 for inside, 1 for up, 2 for down
	DestinationFloor int	
}

type FloorReachedMessage struct {
	ElevatorID int //Lift number 1,2,3..n
	Floor int //Current floor 
	Direction int //1 for up, 2 for down
}


func print_udp_message(msg udp.Udp_message){ 	
	fmt.Printf("msg:  \n \t raddr = %s \n \t data = %s \n \t length = %v \n", msg.Raddr, msg.Data, msg.Length)
}



//Continously sends an alivemessage to the network
func SendAliveMessage(send_ch chan udp.Udp_message, ElevatorID int) {
	Imalive := &AliveMessage {
		ElevatorID: ElevatorID,
		AliveMessage: "I'm alive"}
	MessageCoded, err := json.Marshal(Imalive)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: AliveMessage\n")
		panic(err)
	}
	for {
		time.Sleep(1*time.Second)	
		snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
		fmt.Printf("Sending------\n")
		send_ch <- snd_msg
		//print_udp_message(snd_msg)		
	}
}

//Sends a new order once
func SendNewOrderMessage(send_ch chan udp.Udp_message, ElevatorID int, OrderType int, DestinationFloor int) {
	Order := &NewOrderMessage {
		ElevatorID: ElevatorID,
		OrderType: OrderType,
		DestinationFloor: DestinationFloor}
	MessageCoded, err := json.Marshal(Order)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: NewOrderMessage\n")
		panic(err)
	}
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	//fmt.Printf("Sending------\n")
	send_ch <- snd_msg
}	
	

func SendFloorReached(send_ch chan udp.Udp_message, ElevatorID int, Floor int, Direction int) {
	FloorReached := &FloorReachedMessage {
		ElevatorID: ElevatorID,
		Floor: Floor,
		Direction: Direction}
	MessageCoded, err := json.Marshal(FloorReached)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: FloorReachedMessage\n")
		panic(err)
	}
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	//fmt.Printf("Sending------\n")
	send_ch <- snd_msg
}










func ReadFromNetwork (receive_ch chan udp.Udp_message){
	for {
		fmt.Printf("Receiving----\n")
		rcv_msg:= <- receive_ch
		NetworkMessage := rcv_msg.Data
		var MessageDecoded AliveMessage
		for i := 0; i < len(NetworkMessage); i++ { //Format the message to be decoded
			if NetworkMessage[i] == 0 {
				NetworkMessage = NetworkMessage[:i]
				break
			}		
		}			
		err := json.Unmarshal(NetworkMessage,&MessageDecoded); 
		if err != nil {	
			fmt.Printf("Error: json.Marshal decoder failed: ReadFromNetwork\n")
			panic(err)
		}
		fmt.Println(MessageDecoded)
	}
}



func main (){
	send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	err := udp.Udp_init(20014, 20014, 60, send_ch, receive_ch)	
	go ReadFromNetwork (receive_ch)
	//go SendAliveMessage(send_ch,1)	
	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
	neverReturn := make (chan int)
	<-neverReturn
}
