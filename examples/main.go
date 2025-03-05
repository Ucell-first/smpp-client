package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"smpp-client/smpp"
)

func main() {
	config := &tls.Config{InsecureSkipVerify: true}

	client, err := smpp.NewClient("smsc-sim.melroselabs.net:2775", config)
	if err != nil {
		log.Fatal("Ulanish xatosi:", err)
	}
	defer client.Close()

	if err := client.Bind("username", "password"); err != nil {
		log.Fatal("Bind xatosi:", err)
	}

	err = client.SendSMS(
		"1616",
		"+998931727570",
		"Salom, SMPP orqali jo'natilgan xabar! 🚀",
	)

	if err != nil {
		log.Fatal("Xabar jo'natish xatosi:", err)
	}

	fmt.Println("Xabar muvaffaqiyatli jo'natildi!")
}

func isUnicode(text string) bool {
	for _, r := range text {
		if r > 127 {
			return true
		}
	}
	return false
}

// UCS2 kodlash
func encodeUCS2(text string) []byte {
	// Unicode kodlash logikasi
	return []byte(text) // Soddalik uchun
}
