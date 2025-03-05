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
	systemID := append([]byte(b.SystemID), 0)
	password := append([]byte(b.Password), 0)
	systemType := append([]byte(b.SystemType), 0)
	addressRange := append([]byte(b.AddressRange), 0)

	data := make([]byte, 16)
	binary.BigEndian.PutUint32(data[4:8], b.CommandID)
	binary.BigEndian.PutUint32(data[8:12], b.CommandStatus)
	binary.BigEndian.PutUint32(data[12:16], b.SequenceNumber)

	data = append(data, systemID...)
	data = append(data, password...)
	data = append(data, systemType...)
	data = append(data, b.InterfaceVersion)
	data = append(data, b.AddrTon)
	data = append(data, b.AddrNpi)
	data = append(data, addressRange...)

	binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
	return data, nil
}

// pdu.go fayliga qo'shing
type SubmitSM struct {
	Header
	ServiceType          string
	SourceAddrTon        uint8
	SourceAddrNpi        uint8
	SourceAddr           string
	DestAddrTon          uint8
	DestAddrNpi          uint8
	DestinationAddr      string
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   uint8
	DataCoding           uint8
	ShortMessage         []byte
}

func (s *SubmitSM) MarshalBinary() ([]byte, error) {
	// SubmitSM uchun marshal logikasi
	buf := new(bytes.Buffer)

	// ServiceType (1 octet + NULL)
	buf.WriteString(s.ServiceType)
	buf.WriteByte(0)

	// Source Address TON/NPI
	buf.WriteByte(s.SourceAddrTon)
	buf.WriteByte(s.SourceAddrNpi)

	// Source Address
	buf.WriteString(s.SourceAddr)
	buf.WriteByte(0)

	// Destination Address TON/NPI
	buf.WriteByte(s.DestAddrTon)
	buf.WriteByte(s.DestAddrNpi)

	// Destination Address
	buf.WriteString(s.DestinationAddr)
	buf.WriteByte(0)

	// Protocol ID
	buf.WriteByte(s.ProtocolID)

	// Priority Flag
	buf.WriteByte(s.PriorityFlag)

	// Schedule Delivery Time
	buf.WriteString(s.ScheduleDeliveryTime)
	buf.WriteByte(0)

	// Validity Period
	buf.WriteString(s.ValidityPeriod)
	buf.WriteByte(0)

	// Registered Delivery
	buf.WriteByte(s.RegisteredDelivery)

	// Data Coding
	buf.WriteByte(s.DataCoding)

	// Message Length
	buf.WriteByte(byte(len(s.ShortMessage)))

	// Message Content
	buf.Write(s.ShortMessage)

	return buf.Bytes(), nil
}

// Unbind PDU struktura
type Unbind struct {
	Header
}

func NewUnbind(seq uint32) *Unbind {
	return &Unbind{
		Header: Header{
			CommandID:      0x00000006,
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
