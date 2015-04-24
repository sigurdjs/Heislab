package main

import "fmt"

type Dog struct {
	name string
}

func (d Dog) Say() {
    fmt.Println("Woof!")
    d.name = "fido"
    fmt.Println(d)
}

func main() {
    d := &Dog{}
    fmt.Println(d.name)
    d.Say()
    fmt.Println(d.name)
}