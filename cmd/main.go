package main

import (
	"log"
	"smpp-client/send"
)

func main() {
	err := send.SendMessage("+998900417570")
	if err != nil {
		log.Fatal(err)
	}
}
