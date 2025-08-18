package protocol

import (
	"errors"
	"fmt"
)

const (
	CmdSet   = 0x01
	CmdGet   = 0x02
	CmdError = 0x03
)

const MessageSizeBytes = 5

type BinaryMessage struct {
	Command       uint8
	CheckboxIndex uint32
	IsChecked     bool
}

func EncodeBinaryMessage(msg BinaryMessage) []byte {
	if msg.CheckboxIndex > 0xFFFFFF {
		// Index too large for 24-bit encoding
		msg.CheckboxIndex = 0xFFFFFF
	}

	data := make([]byte, MessageSizeBytes)
	data[0] = msg.Command
	data[1] = uint8(msg.CheckboxIndex >> 16) // High byte
	data[2] = uint8(msg.CheckboxIndex >> 8)  // Mid byte  
	data[3] = uint8(msg.CheckboxIndex)       // Low byte
	
	if msg.IsChecked {
		data[4] = 0x01
	} else {
		data[4] = 0x00
	}
	
	return data
}

func DecodeBinaryMessage(data []byte) (BinaryMessage, error) {
	if len(data) != MessageSizeBytes {
		return BinaryMessage{}, fmt.Errorf("invalid message size: expected %d bytes, got %d", MessageSizeBytes, len(data))
	}

	command := data[0]
	if command != CmdSet && command != CmdGet && command != CmdError {
		return BinaryMessage{}, errors.New("invalid command byte")
	}

	checkboxIndex := uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	isChecked := data[4] == 0x01

	return BinaryMessage{
		Command:       command,
		CheckboxIndex: checkboxIndex,
		IsChecked:     isChecked,
	}, nil
}