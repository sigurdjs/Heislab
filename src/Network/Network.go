package network



import (
	"types"
	"fmt"
	"udp"
	"encoding/json"
	//"time"
)








/*func print_udp_message(msg udp.Udp_message){ 	
	fmt.Printf("msg:  \n \t raddr = %s \n \t data = %s \n \t length = %v \n", msg.Raddr, msg.Data, msg.Length)
}*/



//Continously sends an alivemessage to the network
func SendMessage(send_ch chan udp.Udp_message,MessageType int,AliveMessage string, ButtonType int,DestinationFloor int,CurrentFloor int,ElevatorID int) {
	Imalive := &types.NetworkMessage {
		MessageType: MessageType, 
		AliveMessage: AliveMessage, 
		ButtonType: ButtonType, 
		DestinationFloor: DestinationFloor, 
		CurrentFloor: CurrentFloor, 
		ElevatorID: ElevatorID}
	MessageCoded, err := json.Marshal(Imalive)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: AliveMessage\n")
		panic(err)
	}
	//for {
		//time.Sleep(1*time.Second)	
	fmt.Printf("Lenght is %v \n", len(MessageCoded))
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	send_ch <- snd_msg
	//}
}
/*
//Sends a new order once
func SendNewOrderMessage(send_ch chan udp.Udp_message, ElevatorID int, ButtonType int, DestinationFloor int, MessageType int) {
	Order := &NetworkMessage {
		MessageType: MessageType, 
		AliveMessage: "I'm Alive", 
		ButtonType: ButtonType, 
		DestinationFloor: DestinationFloor, 
		CurrentFloor: -1, 
		ElevatorID: ElevatorID}
	MessageCoded, err := json.Marshal(Order)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: NewOrderMessage\n")
		panic(err)
	}
	//fmt.Printf("Lenght is %v \n", len(MessageCoded))
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	//fmt.Printf("Sending------\n")
	send_ch <- snd_msg
}	

func SendSetLights(send_ch chan udp.Udp_message, ElevatorID int, ButtonType int, DestinationFloor int,LightState bool) {
	var MessageType int
	if LightState == true {
		MessageType = 6
	} else {
		MessageType = 7
	}
	Order := &NetworkMessage {
		MessageType: MessageType, 
		AliveMessage: "I'm Alive", 
		ButtonType: ButtonType, 
		DestinationFloor: DestinationFloor, 
		CurrentFloor: -1, 
		ElevatorID: ElevatorID}
	MessageCoded, err := json.Marshal(Order)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: NewOrderMessage\n")
		panic(err)
	}
	//fmt.Printf("Lenght is %v \n", len(MessageCoded))
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	//fmt.Printf("Sending------\n")
	send_ch <- snd_msg
}	
	
//Sends the current floor of the lift once
func SendCurrentFloor(send_ch chan udp.Udp_message, ElevatorID int, Destination int, CurrentFloor int) {
	SendFloor := &NetworkMessage {
		MessageType: 3, 
		AliveMessage: "I'm Alive", 
		ButtonType: -1, 
		DestinationFloor: Destination, 
		CurrentFloor: CurrentFloor, 
		ElevatorID: ElevatorID}
	MessageCoded, err := json.Marshal(SendFloor)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: FloorReachedMessage\n")
		panic(err)
	}
	//fmt.Printf("Lenght is %v \n", len(MessageCoded))
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	//fmt.Printf("Sending------\n")
	send_ch <- snd_msg
}
*/
func ReadFromNetwork (receive_ch chan udp.Udp_message,MessageToProcess chan types.NetworkMessage){
	var MessageDecoded types.NetworkMessage
	for {	
		rcv_msg := <- receive_ch
		//Message := rcv_msg.Data
		//Message = Message[:rcv_msg.Length]
		fmt.Printf("Lenght recieved is %v \n", rcv_msg.Length)
		err := json.Unmarshal(rcv_msg.Data[:rcv_msg.Length],&MessageDecoded) 
		if err != nil {	
			fmt.Printf("Error: json.Marshal decoder failed: ReadFromNetwork\n")
			panic(err)
		}
		MessageToProcess <- MessageDecoded
	}
}



/*func main (){
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
}*/











