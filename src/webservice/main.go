package main

import (
	"fmt"
	"net/http"

	"theorylab.com/shopwatch/src/webservice/controllers"
)

func main() {
	port := 3000
	fmt.Println("Starting Server...")
	fmt.Println("Registering Controllers...")
	controllers.RegisterControllers()
	fmt.Println("Server Listening at port", port)
	http.ListenAndServe(":3000", nil)
}
