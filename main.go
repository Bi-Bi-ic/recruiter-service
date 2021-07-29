package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	app "github.com/rgrs-x/service/api/app"
)

// Color message for port binding
const (
	printColor   = "\033[38;5;%dm%s\033[39;49m\n"
	errColor     = 9
	warningColor = 3
	developColor = 10
)

func main() {

	router := app.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" //localhost
	}

	fmt.Printf(printColor, warningColor, port)

	//@ Launch the app, visit localhost:8000/api
	err := http.ListenAndServe(":"+port, router)

	// Handle to another port if binding error
	if result := catchErrorBinding(err); result {
		port = "8081"
		fmt.Printf(printColor, developColor, "Changed to port: "+port)
		err = http.ListenAndServe(":"+port, router)
		if err != nil {
			fmt.Printf(printColor, errColor, err)
		}
	} else {
		fmt.Printf(printColor, warningColor, err)
	}
}

// For handling second port
func catchErrorBinding(err error) bool {
	if strings.Contains(err.Error(), "bind: address already in use") {
		fmt.Printf(printColor, errColor, err)
		return true
	}
	return false
}
