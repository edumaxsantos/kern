package main

import "github.com/charmbracelet/lipgloss"

const dot = "⬤"

var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	ledOff = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999")).Render(dot)
	ledOn  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9e01")).Bold(true).Render(dot)
)
