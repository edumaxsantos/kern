package main

import "errors"

// The whole concept of Pins and the way it is built here
// is inspired by tinyGo library where it contains machine.Pin.
// This way we can also convert the information here instead
// of relying on the whole library.

type Pin uint8

const (
	// Note: start at port B because there is no port A.
	portB Pin = iota * 8
	portC
	portD
)

const (
	PB0 = portB + 0
	PB1 = portB + 1
	PB2 = portB + 2
	PB3 = portB + 3
	PB4 = portB + 4
	PB5 = portB + 5
	PC0 = portC + 0
	PC1 = portC + 1
	PC2 = portC + 2
	PC3 = portC + 3
	PC4 = portC + 4
	PC5 = portC + 5
	PC6 = portC + 6
	PD0 = portD + 0
	PD1 = portD + 1
	PD2 = portD + 2
	PD3 = portD + 3
	PD4 = portD + 4
	PD5 = portD + 5
	PD6 = portD + 6
	PD7 = portD + 7
)

func GetPin(pin string) (Pin, error) {
	value, ok := Pins[pin]

	if !ok {
		return 0, errors.New("Pin not found")
	}

	return value, nil
}

var Pins = map[string]Pin{
	"D0":    PD0,
	"D1":    PD1,
	"D2":    PD2,
	"D3":    PD3,
	"D4":    PD4,
	"D5":    PD5,
	"D6":    PD6,
	"D7":    PD7,
	"D8":    PB0,
	"D9":    PB1,
	"D10":   PB2,
	"D11":   PB3,
	"D12":   PB4,
	"D13":   PB5,
	"D14":   PC0,
	"D15":   PC1,
	"D16":   PC2,
	"D17":   PC3,
	"D18":   PC4,
	"D19":   PC5,
	"RESET": PC6,
	"A1":    PC0,
	"A2":    PC1,
	"A3":    PC2,
	"A4":    PC3,
	"A5":    PC4,
	"A6":    PC5,
}
