package main

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/go-redis/redis/v8"
)

type model struct {
	session  ssh.Session
	rdb      *redis.Client
	db       *DB
	username string
	channel  string
	messages []Message
	input    string
	width    int
	height   int
}

type Message struct {
	Username  string
	Content   string
	Timestamp int64
}

func initialModel(s ssh.Session, rdb *redis.Client, db *DB, width, height int) model {
	return model{
		session: s,
		rdb:     rdb,
		db:      db,
		channel: "general",
		width:   width,
		height:  height,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.input != "" {
				if m.username == "" {
					// Set username on first input
					m.username = m.input
					m.input = ""
				} else {
					// TODO: Send message to Redis
					m.input = ""
				}
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	header := style.Render("[Textale]") + "\n\n"
	
	if m.username == "" {
		return header + "Welcome! Please enter your username:\n> " + m.input
	}

	channelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF"))
	
	chatView := header +
		channelStyle.Render("Channel: #"+m.channel) + "\n\n" +
		"Messages will appear here...\n\n" +
		"─────────────────────────────────────\n" +
		"> " + m.input + "\n\n" +
		lipgloss.NewStyle().Faint(true).Render("Press Ctrl+C to quit")

	return chatView
}
