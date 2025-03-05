package smpp

import "errors"

var (
	ErrNotConnected = errors.New("connection not established")
	ErrBindFailed   = errors.New("bind failed")
	ErrSubmitFailed = errors.New("message submission failed")
	ErrInvalidPDU   = errors.New("invalid PDU format")
	ErrTimeout      = errors.New("operation timed out")
)
