package main

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	incoming chan Message
	cancel   context.CancelFunc
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

func (m *model) Init() tea.Cmd {
	// Load recent history for the starting channel
	msgs, _ := fetchRecentMessages(m.rdb, m.channel)
	m.messages = msgs

	// Start pubsub subscription for realtime updates
	m.incoming = make(chan Message, 32)
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	_ = subscribeChannel(ctx, m.rdb, m.channel, m.incoming)

	// Begin waiting for incoming messages
	return waitForIncoming(m.incoming)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		case "enter":
			if m.input != "" {
				if m.username == "" {
					// Set username on first input
					m.username = m.input
					m.input = ""
				} else {
					// Send message to Redis
					_ = sendMessage(m.rdb, m.channel, Message{
						Username:  m.username,
						Content:   m.input,
						Timestamp: time.Now().Unix(),
					})
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
	case incomingMsg:
		// Append new realtime message
		m.messages = append(m.messages, msg.Message)
		// Keep waiting for more messages
		return m, waitForIncoming(m.incoming)
	}
	return m, nil
}

	func (m *model) View() string {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	header := style.Render("[Textale]") + "\n\n"
	
	if m.username == "" {
		return header + "Welcome! Please enter your username:\n> " + m.input
	}

	channelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF"))
	
	// Render messages
	msgs := ""
	for _, mm := range m.messages {
		ts := time.Unix(mm.Timestamp, 0).Format("15:04:05")
		msgs += lipgloss.NewStyle().Faint(true).Render(ts) + " " +
			lipgloss.NewStyle().Bold(true).Render(mm.Username) + ": " + mm.Content + "\n"
	}

	chatView := header +
		channelStyle.Render("Channel: #"+m.channel) + "\n\n" +
		msgs + "\n" +
		"─────────────────────────────────────\n" +
		"> " + m.input + "\n\n" +
		lipgloss.NewStyle().Faint(true).Render("Press Ctrl+C to quit")

	return chatView
}

// tea cmd to wait for one incoming message
func waitForIncoming(ch <-chan Message) tea.Cmd {
	return func() tea.Msg {
		m := <-ch
		return incomingMsg{Message: m}
	}
}
