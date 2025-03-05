package send

import (
	"log"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

func SendMessage(phone string) error {
	trans := &smpp.Transmitter{
		Addr:   "smscsim.smpp.org:2775",
		User:   "TytZd4G7K7lcZyk",
		Passwd: "SJCuybzs",
	}

	conn := trans.Bind()
	select {
	case status := <-conn:
		if status.Status() != smpp.Connected {
			log.Fatalf("Ulanishda xatolik: %v", status.Error())
			return status.Error()
		}
		log.Println("SMPP serverga muvaffaqiyatli ulandik.")
	case <-time.After(5 * time.Second):
		log.Fatal("Ulanish vaqti tugadi")
	}
	resp, err := trans.Submit(&smpp.ShortMessage{
		Src:      "12345",
		Dst:      phone,
		Text:     pdutext.UCS2("Hello, this is a test message!"),
		Register: pdufield.NoDeliveryReceipt,
	})
	if err != nil {
		log.Fatalf("SMS jo‘natishda xatolik: %v", err)
		return err
	}
	log.Printf("SMS jo‘natildi, Msg ID: %s", resp.RespID())

	trans.Close()
	return nil
}
