package main

import (
	"bufio"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tarm/serial"
	serialV1 "go.bug.st/serial.v1"
	"log"
)

func LoadPorts() tea.Cmd {
	return func() tea.Msg {
		log.Println("Loading ports")
		portsList, err := serialV1.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}

		if len(portsList) == 0 {
			log.Fatal("No serial ports found")
		}

		return portsLoaded{
			Ports: portsList,
		}
	}
}

func ListenForData(stream *serial.Port, sub chan []byte) tea.Cmd {
	return func() tea.Msg {
		log.Println("preparing listen for data")
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			bytes := scanner.Bytes()
			sub <- bytes
			log.Println("Published to the channel.")
		}

		log.Println("no new data")

		return nil
	}
}

func WaitForData(sub chan []byte) tea.Cmd {
	return func() tea.Msg {
		log.Println("Waiting for a message to come from the channel")
		value := <-sub

		log.Println("Got a message from the channel")

		message, err := De
		code(value)
		if err != nil {
			log.Fatalf("We got error: %s. Value is: %x\nIn string: %s\n", err, value, string(value))
		}
		log.Printf("Got new data : %+v\n", message)
		log.Printf("Data is: %s\n", message.DataAsString())
		log.Printf("Actual bytes: %08b\n", value)
		return responseMsg{
			bytes:   value,
			message: *message,
		}
	}
}
