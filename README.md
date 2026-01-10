# Textale

Textale is a real-time disappearing message chat accessible through SSH

## How to use

Connect instantly from any terminal:

```bash
ssh textale.mafupa.dev
```

## Features

### For Everyone

- **Minimalist TUI**: Clean, distraction-free interface built for the terminal
- **Multi-Channel Navigation**: Seamlessly switch between different conversation rooms
- **Presence Awareness**: See who's currently online in real-time
- **Custom Usernames**: Set your display name on the fly
- **Ephemeral Messages**: All messages disappear after a configurable time period
- **SSH-Native**: Access from anywhere with SSH—no browser or app required

### For Administrators

- **Channel Management**: Create, rename, or archive conversation channels
- **Moderation Tools**: Ban or allow specific users to maintain community standards
- **Server Statistics**: Monitor active users, message throughput, and system health
- **Configurable Retention**: Set message lifetime per channel

## Technology Stack

### Core

- **Go**: High-performance backend with excellent concurrency support
- **Bubble Tea**: Terminal UI framework for beautiful, interactive TUIs
- **Lip Gloss**: Styling and layout for terminal interfaces
- **SSH Server**: Custom SSH server implementation for connection handling

### Data & Storage

- **Redis**: In-memory message queuing and session management
- **SQLite**: Lightweight persistence for user settings and channel metadata

### Infrastructure

- **Docker**: Containerized deployment for easy scaling
- **systemd**: Process management and auto-restart on the server

## Contributing

This is a passion project exploring the intersection of retro computing and modern real-time systems. Contributions are welcomed.

## License

MIT
