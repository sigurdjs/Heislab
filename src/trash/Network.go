package Network

import (
	"fmt"
	"udp"
	"encoding/json"
	"time"
)



type NetworkMessage struct {
	MessageType int //1 for alive, 2 for order, 3 for floor update
	AliveMessage string //Message just to send something	
	ButtonType int //0 for up, 1 for down, 2 for inside, -1 for unused
	DestinationFloor int //-1 for unused
	CurrentFloor int //Current floor or destination floor depending on message type, also -1 for unused
	ElevatorID int //Lift number 1,2,3..n
}




/*func print_udp_message(msg udp.Udp_message){ 	
	fmt.Printf("msg:  \n \t raddr = %s \n \t data = %s \n \t length = %v \n", msg.Raddr, msg.Data, msg.Length)
}*/



//Continously sends an alivemessage to the network
func SendAliveMessage(send_ch chan udp.Udp_message, ElevatorID int) {
	Imalive := &NetworkMessage {
		MessageType: 1,
		ElevatorID: ElevatorID,
		AliveMessage: "I'm alive",
		ButtonType : -1,
		DestinationFloor: -1,
		CurrentFloor: -1}
	MessageCoded, err := json.Marshal(Imalive)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: AliveMessage\n")
		panic(err)
	}
	for {
		time.Sleep(1*time.Second)	
		snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
		send_ch <- snd_msg
	}
}

//Sends a new order once
func SendNewOrderMessage(send_ch chan udp.Udp_message, ElevatorID int, ButtonType int, DestinationFloor int) {
	Order := &NetworkMessage {
		MessageType: 2,
		ElevatorID: ElevatorID,
		ButtonType: ButtonType,
		DestinationFloor: DestinationFloor,
		CurrentFloor : -1,
		AliveMessage: ""}
	MessageCoded, err := json.Marshal(Order)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: NewOrderMessage\n")
		panic(err)
	}
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	//fmt.Printf("Sending------\n")
	send_ch <- snd_msg
}	
	
//Sends the current floor of the lift once
func SendCurrentFloor(send_ch chan udp.Udp_message, ElevatorID int, Destination int, CurrentFloor int) {
	SendFloor := &NetworkMessage {
		MessageType: 3,
		ElevatorID: ElevatorID,
		CurrentFloor: CurrentFloor,
		DestinationFloor: Destination,
		ButtonType: -1,
		AliveMessage: ""}
	MessageCoded, err := json.Marshal(SendFloor)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: FloorReachedMessage\n")
		panic(err)
	}
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	//fmt.Printf("Sending------\n")
	send_ch <- snd_msg
}





func ReadFromNetwork (receive_ch chan udp.Udp_message), MessageToProcess chan NetworkMessage){
	for {
		fmt.Printf("Receiving----\n")
		rcv_msg:= <- receive_ch
		Message := rcv_msg.Data
		var MessageDecoded NetworkMessage
		for i := 0; i < len(Message); i++ { //Format the message to be decoded
			if Message[i] == 0 {
				Message = Message[:i]
				break
			}		
		}			
		err := json.Unmarshal(Message,&MessageDecoded); 
		if err != nil {	
			fmt.Printf("Error: json.Marshal decoder failed: ReadFromNetwork\n")
			panic(err)
		}
		fmt.Println(MessageDecoded.AliveMessage, )
	}
}


/*
func main (){
	send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	err := udp.Udp_init(20014, 20014, 100, send_ch, receive_ch)	
	go SendAliveMessage(send_ch, 1)
	go ReadFromNetwork (receive_ch)
		
	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
	neverReturn := make (chan int)
	<-neverReturn
}

*/









