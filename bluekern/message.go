package message

import (
	"errors"
)

type IO byte

const (
	Input  IO = 0x00
	Output IO = 0x01
)

const (
	STX byte = 0x02
	ETX byte = 0x03
	// First byte masks and shifts
	versionMask  byte = 0xE0 // 11100000: extracts the upper 3 bits
	versionShift      = 5    // shift right by 5 to get the version value
	pinMask      byte = 0x1F // 00011111: extracts the lower 5 bits

	// Second byte masks and shifts
	ioMask     byte = 0x80 // 10000000: extracts the IO flag bit
	ioShift         = 7    // shift right by 7 to get the IO flag value
	rwMask     byte = 0x40 // 01000000: extracts the RW flag bit
	rwShift         = 6    // shift right by 6 to get the RW flag value
	lengthMask byte = 0x3F // 00111111: extracts the lower 6 bits for the length
)

type RW byte

const (
	Read  RW = 0x00
	Write RW = 0x01
)

// Message
// Minimum number of bits (with no data) is 40 bits (5 bytes)
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

func (m *Message) DataAsString() string {
	return string(m.Data)
}

func (m *Message) Encode() []byte {
	header := []byte{
		m.STX,
		// Combine two fields in one byte. the high 3 bits (version) with the low 5 bits (pin).
		// How does it work?
		// Shift version bits to the left by 5 bits, leaving 3 bits. (e.g. Version 1 = 0000 0001, then shift left by 5 -> 0010 0000)
		// The second part will remove any value from the high 3 bits, turning them 0.
		// Combined we will have the 3 first bits from Version and the last 5 bits from Pin.
		//        v- version
		//        0000 0000
		// version -^^- pin
		(m.Version << versionShift) | (m.Pin & pinMask),
		// Combine 3 fields into one byte. First high bit is for IO, second high bit is for RW, the 6 last bits for length.
		// First shift IO bits to the left by 7 bits, so only 1 bit is counted and placed at the 8th position.
		// Second shift RW bits to the left by 6 bits, so only 1 bit is counted and placed at the 7th position.
		(byte(m.IO) << ioShift) | (byte(m.RW) << rwShift) | (m.Length & lengthMask),
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

	firstByte := encoded[1]
	version := (firstByte & versionMask) >> versionShift
	pin := firstByte & pinMask

	secondByte := encoded[2]
	io := (secondByte & ioMask) >> ioShift
	rw := (secondByte & rwMask) >> rwShift
	length := secondByte & lengthMask

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

	message.Checksum = calculatedChecksum

	return message, nil
}
