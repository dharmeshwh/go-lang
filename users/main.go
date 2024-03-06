package main

import (
	"fmt"
	"log"
	"net/http"
	"users/router"
)

func main() {
	fmt.Println("Server is getting started...")

	log.Fatal(http.ListenAndServe(":4000", router.Router()))

	fmt.Println("Listening at port 4000 ...")
}
