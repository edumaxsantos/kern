package main

import "errors"

type IO byte

const (
	Input  IO   = 0x00
	Output IO   = 0x01
	STX    byte = 0x02
	ETX    byte = 0x03
)

type RW byte

const (
	Read  RW = 0x00
	Write RW = 0x01
)

// Message
// +-----+-----+-----+-----+-----+-----+------+----------+-----+
// | STX | Ver | Pin | I/O | R/W | Len | Data | Checksum | ETX |
// +-----+-----+-----+-----+-----+-----+------+----------+-----+
// |   8 |   3 |   5 |   1 |   1 |   6 | n    |        8 |   8 |
// +-----+-----+-----+-----+-----+-----+------+----------+-----+
type Message struct {
	STX      byte
	Version  byte
	Pin      byte
	IO       IO
	RW       RW
	Length   byte
	Data     []byte
	Checksum byte
	ETX      byte
}

// calculateChecksum use the following fields to generate checksum:
// Message.Version, Message.Pin, Message.IO, Message.RW,
// Message.Length, Message.Data
func calculateChecksum(msg []byte) byte {
	var sum byte
	for _, b := range msg {
		sum += b
	}
	c := 256
	return byte(int(sum) % c)
}

func (m *Message) Encode() []byte {
	header := []byte{
		m.STX,
		(m.Version << 5) | (m.Pin & 0x1F),
		(byte(m.IO) << 7) | (byte(m.RW) << 6) | (m.Length & 0x3F),
	}

	msg := append(header, m.Data...)

	checksum := calculateChecksum(m.checksumData())
	msg = append(msg, checksum, m.ETX)

	return msg
}

func (m *Message) checksumData() []byte {
	return append([]byte{m.Version, m.Pin, byte(m.IO), byte(m.RW), m.Length}, m.Data...)
}

func Decode(encoded []byte) (*Message, error) {
	if len(encoded) < 4 {
		return nil, errors.New("message too short")
	}

	if encoded[0] != STX {
		return nil, errors.New("invalid STX")
	}

	if encoded[len(encoded)-1] != ETX {
		return nil, errors.New("invalid ETX")
	}

	header1 := encoded[1]
	version := (header1 & 0xE0) >> 5
	pin := header1 & 0x1F

	header2 := encoded[2]
	io := (header2 & 0x80) >> 7
	rw := (header2 & 0x40) >> 6
	length := header2 & 0x3F

	data := encoded[3 : 3+length]

	message := &Message{
		STX:     STX,
		Version: version,
		Pin:     pin,
		IO:      IO(io),
		RW:      RW(rw),
		Length:  length,
		Data:    data,
		ETX:     ETX,
	}

	checkSum := encoded[3+length]

	calculatedChecksum := calculateChecksum(message.checksumData())

	if checkSum != calculatedChecksum {
		return nil, errors.New("checksum mismatch")
	}

	return message, nil
}

func (m *Message) sendOn() {
	sendMessage(*m, string(ON))
}

func (m *Message) sendOff() {
	sendMessage(*m, string(OFF))
}
