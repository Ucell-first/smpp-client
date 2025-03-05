package smpp

import (
	"crypto/tls"
	"net"
)

func DialWithTLS(addr string, config *tls.Config) (net.Conn, error) {
	return tls.Dial("tcp", addr, config)
}
