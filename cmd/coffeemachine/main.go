package main

import (
	"coffeemachine/handler"
	"fmt"
)

func main() {
	// pass onemachine config and order as params.
	err := handler.InitialSetup("machine1", "orders1")
	if err != nil {
		fmt.Println(err)
	}
}
