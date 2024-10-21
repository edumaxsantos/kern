package main

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	data := []byte("1")
	length := byte(len(data))
	checksumMsg := append([]byte{VERSION, BUILTIN_LED, byte(Output), byte(Write), length}, data...)
	checksum := calculateChecksum(checksumMsg)

	//{STX:2 Version:1 Pin:5 IO:1 RW:0 Length:6 Data:[50 102 97 108 115 101] Checksum:0 ETX:3}
	message := Message{
		STX:      STX,
		Version:  VERSION,
		Pin:      BUILTIN_LED,
		IO:       Output,
		RW:       Read,
		Length:   length,
		Data:     data,
		Checksum: checksum,
		ETX:      ETX,
	}

	encoded := message.Encode()

	decoded, err := Decode(encoded)

	if err != nil {
		t.Error(err)
	}

	if decoded.Checksum != message.Checksum {
		t.Errorf("Checksum mismatch. Expected: %d, got: %d", message.Checksum, decoded.Checksum)
	}

	fmt.Println(len(encoded))

	fmt.Printf("Message: %+v\n", message)
	fmt.Printf("Encoded binary: %b\n", encoded)
	fmt.Printf("Encoded hex: %x\n", encoded)
	fmt.Printf("Decoded: %+v\n", decoded)
}
