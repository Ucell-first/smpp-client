package test

import (
	"smpp-client/send"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	err := send.SendMessage("+998900417570")
	require.NoError(t, err)
}
