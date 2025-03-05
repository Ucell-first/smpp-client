package smpp

import (
	"bytes"
	"encoding/binary"
)

// Umumiy PDU header
type Header struct {
	CommandLength  uint32
	CommandID      uint32
	CommandStatus  uint32
	SequenceNumber uint32
}

// BindTransmitter PDU
type BindTransmitter struct {
	Header
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion uint8
	AddrTon          uint8
	AddrNpi          uint8
	AddressRange     string
}

func (b *BindTransmitter) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// Header (avvaliga CommandLength joy tashlab qo'yamiz)
	buf.Write(make([]byte, 4)) // CommandLength uchun joy
	binary.Write(&buf, binary.BigEndian, b.CommandID)
	binary.Write(&buf, binary.BigEndian, b.CommandStatus)
	binary.Write(&buf, binary.BigEndian, b.SequenceNumber)

	// Body
	writeCString(&buf, b.SystemID)
	writeCString(&buf, b.Password)
	writeCString(&buf, b.SystemType)
	buf.WriteByte(b.InterfaceVersion)
	buf.WriteByte(b.AddrTon)
	buf.WriteByte(b.AddrNpi)
	writeCString(&buf, b.AddressRange)

	// CommandLength ni hisoblash
	data := buf.Bytes()
	binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))

	return data, nil
}

func writeCString(buf *bytes.Buffer, s string) {
	buf.WriteString(s)
	buf.WriteByte(0)
}

// pdu.go fayliga qo'shing
type SubmitSM struct {
	Header
	ServiceType     string
	SourceAddrTon   uint8
	SourceAddrNpi   uint8
	SourceAddr      string
	DestAddrTon     uint8
	DestAddrNpi     uint8
	DestinationAddr string
	DataCoding      uint8
	ShortMessage    []byte
}

func (s *SubmitSM) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// Header (CommandLength keyin yoziladi)
	buf.Write(make([]byte, 4))
	binary.Write(&buf, binary.BigEndian, s.CommandID)
	binary.Write(&buf, binary.BigEndian, s.CommandStatus)
	binary.Write(&buf, binary.BigEndian, s.SequenceNumber)

	// Body
	writeCString(&buf, s.ServiceType)
	buf.WriteByte(s.SourceAddrTon)
	buf.WriteByte(s.SourceAddrNpi)
	writeCString(&buf, s.SourceAddr)
	buf.WriteByte(s.DestAddrTon)
	buf.WriteByte(s.DestAddrNpi)
	writeCString(&buf, s.DestinationAddr)
	buf.WriteByte(0)       // ProtocolID
	buf.WriteByte(0)       // PriorityFlag
	writeCString(&buf, "") // ScheduleDeliveryTime
	writeCString(&buf, "") // ValidityPeriod
	buf.WriteByte(0)       // RegisteredDelivery
	buf.WriteByte(0)       // ReplaceIfPresentFlag
	buf.WriteByte(s.DataCoding)
	buf.WriteByte(0) // SMDefaultMsgID
	buf.WriteByte(uint8(len(s.ShortMessage)))
	buf.Write(s.ShortMessage)

	// CommandLength ni hisoblash
	data := buf.Bytes()
	binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))

	return data, nil
}

// Unbind PDU struktura
type Unbind struct {
	Header
}

func NewUnbind(seq uint32) *Unbind {
	return &Unbind{
		Header: Header{
			CommandLength:  16, // Faqat header uchun
			CommandID:      0x00000006,
			CommandStatus:  0,
			SequenceNumber: seq,
		},
	}
}

func (u *Unbind) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u.CommandLength)
	binary.Write(buf, binary.BigEndian, u.CommandID)
	binary.Write(buf, binary.BigEndian, u.CommandStatus)
	binary.Write(buf, binary.BigEndian, u.SequenceNumber)
	return buf.Bytes(), nil
}
