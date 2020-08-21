package main

import (
	"coffeemachine/handler"
	"fmt"
)

func main() {
	// pass onemachine config and order as params.
	fmt.Println(handler.InitialSetup("machine1", "orders1"))
}
