package main


import (
	"fmt"
	"udp"
	"time"
	"strconv"
)

var masterIsAlive bool
var networkMessage string

func print_udp_message(msg udp.Udp_message){ 	
	fmt.Printf("msg:  \n \t raddr = %s \n \t data = %s \n \t length = %v \n", msg.Raddr, msg.Data, msg.Length)
}

func SendMessage(send_ch chan udp.Udp_message, message struct){
	for {
		time.Sleep(1*time.Second)	
		snd_msg := udp.Udp_message{Raddr:"broadcast", Data:message, Length:len(message)}
		//fmt.Printf("Sending------\n")
		send_ch <- snd_msg
		//print_udp_message(snd_msg)		
	}
}


func ReadFromNetwork (receive_ch chan udp.Udp_message){
	for {
		time.Sleep(500)
		fmt.Printf("Receiving----\n")
		rcv_msg:= <- receive_ch
		print_udp_message(rcv_msg)
		networkMessage = rcv_msg.Data		
	}
}

func IsMasterAlive(receive_ch chan udp.Udp_message) {
	time.Sleep(90*time.Second)
	for {
		time.Sleep(500)
		fmt.Printf("Receiving----\n")
		rcv_msg:= <- receive_ch
		print_udp_message(rcv_msg)
		networkMessage = rcv_msg.Data		
	}
}


func main (){
	send_ch := make (chan udp.Udp_message)
	receive_ch := make (chan udp.Udp_message)
	err := udp.Udp_init(20000, 20000, 1024, send_ch, receive_ch)	
	go SendMessage(send_ch, "Im Alive! Id: "+strconv.Itoa(1))
	go SendMessage(send_ch, "Fuck off!")	
	//go ReadFromNetwork(receive_ch)

	if (err != nil){
		fmt.Print("main done. err = %s \n", err)
	}
		neverReturn := make (chan int)
	<-neverReturn

}
