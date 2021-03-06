package main

import (
	"fmt"
	"net/http"
	"log"
	"kalaxia-game-api/security"
	"kalaxia-game-api/route"
)

func main() {
	if security.InitializeRsaVault() {
		fmt.Println("The RSA keys were generated")
	} else {
		fmt.Println("The RSA keys are already generated")
	}
  	log.Fatal(http.ListenAndServe(":80", route.NewRouter()))
}
