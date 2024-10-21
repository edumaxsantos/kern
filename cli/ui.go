package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tarm/serial"
	"log"
	"time"
)

type item struct {
	title, desc string
	value       string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
func (i item) Value() string       { return i.value }

type model struct {
	list          list.Model
	stream        *serial.Port
	sub           chan []byte
	currentScreen screen
	portsList     list.Model
	textInput     textinput.Model
	currentPin    string
	currentView   view
}

type responseMsg struct {
	bytes   []byte
	message Message
}

type portsLoaded struct {
	Ports []string
}

func (m model) Init() tea.Cmd {
	return LoadPorts()
}

func portSelectionUpdate(m model) (model, tea.Cmd) {
	selectedItem := m.portsList.SelectedItem().(item)

	m.stream = openPort(selectedItem.value)

	m.currentScreen = Main
	m.currentView = TextView

	return m, tea.Batch(m.textInput.Focus(), textinput.Blink, ListenForData(m.stream, m.sub), WaitForData(m.sub))
}

func mainUpdate(m model) (model, tea.Cmd) {
	m.list.StatusMessageLifetime = time.Second
	selectedItem := m.list.SelectedItem().(item)

	var statusCmd tea.Cmd

	pin, err := GetPin(m.currentPin)

	if err != nil {
		log.Fatal(err)
	}

	message := Message{
		STX:     STX,
		Version: VERSION,
		IO:      Output,
		Pin:     byte(pin),
		RW:      Read,
		Length:  0x03,
		Data:    []byte(selectedItem.Value()),
		ETX:     ETX,
	}

	message.Length = byte(len(message.Data))

	if selectedItem.value == "1" || selectedItem.value == "0" {
		message.RW = Write
	}

	message.Checksum = calculateChecksum(message.checksumData())

	log.Printf("Bytes to be sent %+v\n", message)
	bytes := message.Encode()
	log.Printf("Actual bytes: %08b\n", bytes)
	_, err = m.stream.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return m, tea.Batch(WaitForData(m.sub), statusCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.currentScreen == PortSelection {
				return portSelectionUpdate(m)
			}

			if m.currentView == TextView {
				m.currentPin = m.textInput.Value()
				m.currentView = ListView
				return m, nil
			}

			return mainUpdate(m)
		case tea.KeyTab, tea.KeyShiftTab:
			if m.currentView == TextView {
				m.currentView = ListView
			} else {
				m.currentView = TextView
			}
		case tea.KeyRunes:
			if m.currentView == ListView {
				switch msg.String() {
				case "1":
					m.currentView = TextView
					return m, nil
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.portsList.SetSize(msg.Width-h, msg.Height-v)
	case responseMsg:
		log.Printf("Got response msg")
		m.list.StatusMessageLifetime = time.Second * 5
		ledStatus := GetLedStatus(msg.message.Data)
		statusCmd := m.list.NewStatusMessage(statusMessageStyle("LED is currently " + ledStatus))
		if ledStatus == "ON" {
			m.list.Title = "Control the LED " + ledOn
		} else {
			m.list.Title = "Control the LED " + ledOff
		}
		return m, tea.Batch(statusCmd, WaitForData(m.sub))
	case portsLoaded:
		l := make([]list.Item, 0)
		for _, port := range msg.Ports {
			l = append(l, item{title: port, desc: port, value: port})
		}
		m.portsList = list.New(l, list.NewDefaultDelegate(), m.portsList.Width(), m.portsList.Height())
	}

	var cmd tea.Cmd

	if m.currentScreen == Main {
		if m.currentView == ListView {
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
		if m.currentView == TextView {
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}
	m.portsList, cmd = m.portsList.Update(msg)

	return m, cmd
}

func (m model) View() string {

	if m.currentScreen == PortSelection {
		return docStyle.Render(m.portsList.View())
	}

	if m.currentView == TextView {
		return fmt.Sprintf("What is the port you want to access now?\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	} else if m.currentView == ListView {
		return fmt.Sprintf("Current port: %s\n\n%s",
			m.currentPin,
			m.list.View())
	} else if m.currentView == NotSetView {
		return "Error no view set"
	}

	return ""
}

func openPort(name string) *serial.Port {
	config := &serial.Config{Name: name, Baud: 9600}

	stream, err := serial.OpenPort(config)

	if err != nil {
		log.Fatal(err)
	}

	return stream
}

func initModel() model {
	items := []list.Item{
		item{title: "ON", desc: "Turn LED ON", value: "1"},
		item{title: "OFF", desc: "Turn LED OFF", value: "0"},
		item{title: "Retrieve", desc: "Get current LED bytes", value: "2"},
	}
	ti := textinput.New()
	ti.Placeholder = "Port to access"
	ti.CharLimit = 3
	ti.Width = 20

	delegate := list.NewDefaultDelegate()
	delegate.ShortHelpFunc = func() []key.Binding {
		return keys.ShortHelp()
	}

	delegate.FullHelpFunc = func() [][]key.Binding {
		return keys.FullHelp()
	}

	l := list.New(items, delegate, 0, 0)

	m := model{
		list:          l,
		sub:           make(chan []byte),
		currentScreen: PortSelection,
		portsList:     list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		currentView:   NotSetView,
		textInput:     ti,
	}
	m.list.Title = "Control the LED"

	return m
}
