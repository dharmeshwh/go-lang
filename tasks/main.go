package main

import (
	"fmt"
	"log"
	"net/http"
	"tasks/router"
)

func main() {
	fmt.Println("Server is getting started...")

	log.Fatal(http.ListenAndServe(":4001", router.Router()))

	fmt.Println("Listening at port 4001 ...")
}
