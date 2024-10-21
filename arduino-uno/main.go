package main

import (
	"machine"
	"time"
)

const (
	ON  byte = 0x31 // 1
	OFF byte = 0x30 // 0
)

func sendInitialMessage(msg string) {
	initialMessage := Message{
		STX:      STX,
		Version:  0x01,
		Pin:      0x00,
		IO:       IO(0x00),
		RW:       RW(0x00),
		Length:   0,
		Data:     []byte(msg),
		Checksum: 0,
		ETX:      ETX,
	}

	sendMessage(initialMessage, msg)
}

func blink(n int) {
	pin := machine.LED

	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for i := 0; i < n; i++ {
		pin.High()
		time.Sleep(time.Millisecond * 300)
		pin.Low()
		time.Sleep(time.Millisecond * 300)
	}
}

func errorFastBlink() {
	blink(10)
}

func sendMessage(original Message, msg string) {
	original.Data = []byte(msg)
	original.Length = byte(len(original.Data))
	original.Checksum = calculateChecksum(original.checksumData())

	bytes := original.Encode()

	for _, by := range bytes {
		err := machine.Serial.WriteByte(by)
		if err != nil {
			errorFastBlink()
			return
		}
	}

	machine.Serial.WriteByte('\n')
}

func readMessage() *Message {
	bytes := make([]byte, 0)

	for {
		c, err := machine.Serial.ReadByte()

		// no received byte
		if err != nil {
			continue
		}

		// while we have nothing stored, check if the byte
		// we got is STX
		if len(bytes) == 0 {
			if c != STX {
				continue
			}
		}

		bytes = append(bytes, c)

		if c == ETX {
			break
		}
	}

	message, err := Decode(bytes)
	if err != nil {
		errorFastBlink()
		sendInitialMessage(err.Error())
	}

	return message
}

func main() {

	machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 9600,
	})
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Low()

	for {
		msg := readMessage()

		pin := setPin(msg.Pin, msg.IO)

		switch msg.RW {
		case Read:
			if pin.Get() {
				msg.sendOn()
			} else {
				msg.sendOff()
			}
		case Write:
			if msg.Data != nil {
				if msg.Data[0] == ON {
					pin.High()
					msg.sendOn()
				} else if msg.Data[0] == OFF {
					pin.Low()
					msg.sendOff()
				} else {
					sendMessage(*msg, "invalid data")
					errorFastBlink()
				}
			} else {
				sendMessage(*msg, "no data")
			}
		}
	}
}

func setPin(pin byte, io IO) *machine.Pin {
	var pinMode machine.PinMode
	if io == Input {
		pinMode = machine.PinInput
	} else if io == Output {
		pinMode = machine.PinOutput
	}

	machinePin := machine.Pin(pin)

	machinePin.Configure(machine.PinConfig{Mode: pinMode})

	return &machinePin
}
