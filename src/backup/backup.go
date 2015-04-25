package backup

import (
	"io/ioutil"
	."path/filepath"
	"types"
	"encoding/json"
	"fmt" 
)

const Floors = 4
const (
	filepath = "backup.txt"
)

func WriteBackup(OrderQueue []types.Order) {
	MessageToCode := make([]types.Order,len(OrderQueue))
	MessageToCode = OrderQueue
	buf, err := json.Marshal(MessageToCode)
	path, _ := Abs(filepath)
	ioutil.WriteFile(path, buf, 0644)
	if err != nil {
		fmt.Printf("Error: json.Marshal encoder failed: encoding int to byte\n")
		panic(err)
	}

}

func ReadBackup()  []types.Order { 

	path, _ := Abs(filepath)
	buf, err := ioutil.ReadFile(path)
	

	var ElevOrders []types.Order
	if err == nil {}
	
	// from byte to int
	//fmt.Println(buf)
	if len(buf) > 0 {
		err = json.Unmarshal(buf, &ElevOrders) 
			if err != nil {	
				fmt.Printf("Error: json.Marshal decoder failed: decoding byte to int\n")
				panic(err)
			}
		}
	return ElevOrders
}









