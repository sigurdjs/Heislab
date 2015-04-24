package main
import "encoding/json"
import "fmt"


type MessageType1 struct {
    ElevatorID   int
    Message string
}


func main() {
	Imalive := &MessageType1{
		ElevatorID: 1,
		Message: "I'm alive"} 
	MesCoded, _ := json.Marshal(Imalive)
	fmt.Println(MesCoded)
	var MesDecoded MessageType1 
	if err := json.Unmarshal(MesCoded,&MesDecoded); err != nil {
		panic(err)
	}
	fmt.Println(MesDecoded)

}
