package smpp

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type Client struct {
	conn      net.Conn
	sequence  uint32
	mu        sync.Mutex
	tlsConfig *tls.Config
	timeout   time.Duration
}

func NewClient(addr string, tlsConfig *tls.Config) (*Client, error) {
	var conn net.Conn
	var err error

	if tlsConfig != nil {
		conn, err = tls.Dial("tcp", addr, tlsConfig)
	} else {
		conn, err = net.Dial("tcp", addr)
	}

	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return &Client{
		conn:      conn,
		tlsConfig: tlsConfig,
		timeout:   10 * time.Second,
	}, nil
}

func (c *Client) Bind(systemID, password string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	bindPDU := &BindTransmitter{
		SystemID:         systemID,
		Password:         password,
		InterfaceVersion: 0x34,
	}

	data, err := bindPDU.MarshalBinary()
	if err != nil {
		return err
	}

	if err := c.conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		return err
	}

	if _, err := c.conn.Write(data); err != nil {
		return ErrBindFailed
	}

	// Javobni o'qish
	resp := make([]byte, 16)
	if _, err := c.conn.Read(resp); err != nil {
		return ErrBindFailed
	}

	status := binary.BigEndian.Uint32(resp[8:12])
	if status != 0 {
		return ErrBindFailed
	}

	return nil
}

func (c *Client) SendSMS(src, dest, text string) error {
	operation := func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		// SMS kodlash
		var msg []byte
		dataCoding := uint8(0)
		if isUnicode(text) {
			msg = encodeUCS2(text)
			dataCoding = 0x08
		} else {
			msg = []byte(text)
		}

		submit := SubmitSM{
			SourceAddr:      src,
			DestinationAddr: dest,
			ShortMessage:    msg,
			DataCoding:      dataCoding,
		}

		data, err := submit.MarshalBinary()
		if err != nil {
			return err
		}

		if _, err := c.conn.Write(data); err != nil {
			return ErrSubmitFailed
		}

		// Javobni tekshirish
		resp := make([]byte, 16)
		if _, err := c.conn.Read(resp); err != nil {
			return ErrSubmitFailed
		}

		return nil
	}

	// Qayta urinish mexanizmi
	retryBackoff := backoff.NewExponentialBackOff()
	retryBackoff.MaxElapsedTime = 2 * time.Minute
	return backoff.Retry(operation, retryBackoff)
}
func (c *Client) nextSequence() uint32 {
	c.sequence++
	return c.sequence
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	unbind := NewUnbind(c.nextSequence())
	data, err := unbind.MarshalBinary()
	if err != nil {
		return err
	}

	if _, err := c.conn.Write(data); err != nil {
		return err
	}

	// Javobni kutish (30 soniya)
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	resp := make([]byte, 16)
	if _, err := c.conn.Read(resp); err == nil {
		status := binary.BigEndian.Uint32(resp[8:12])
		if status != 0 {
			return fmt.Errorf("unbind failed: status=%d", status)
		}
	}

	return c.conn.Close()
}
