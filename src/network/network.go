package network



import (
	"types"
	"fmt"
	"udp"
	"encoding/json"
	//"time"
)


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
	snd_msg := udp.Udp_message{Raddr:"broadcast", Data:MessageCoded, Length:len(MessageCoded)}
	send_ch <- snd_msg
}

func ReadFromNetwork (receive_ch chan udp.Udp_message,MessageToProcess chan types.NetworkMessage){
	var MessageDecoded types.NetworkMessage
	for {	
		rcv_msg := <- receive_ch
		err := json.Unmarshal(rcv_msg.Data[:rcv_msg.Length],&MessageDecoded) 
		if err != nil {	
			fmt.Printf("Error: json.Marshal decoder failed: ReadFromNetwork\n")
			panic(err)
		}
		MessageToProcess <- MessageDecoded
	}
}
