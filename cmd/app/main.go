package main

import (
	"fmt"
	"gitchecker/routes"
)

func main() {

	panel := routes.PanelHandler()

	fmt.Println("Listening on port 3000")
	err := panel.Listen(":3000")
	if err != nil {
		fmt.Println(err)
		return
	}
}
