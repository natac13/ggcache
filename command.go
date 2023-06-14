package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	SET Command = "SET"
	GET Command = "GET"
	DEL Command = "DEL"
)

type Message struct {
	Command Command
	Key     []byte
	Value   []byte
	TTL     time.Duration
}

func (m *Message) ToBytes() []byte {
	switch m.Command {
	case SET:
		return []byte(fmt.Sprintf("%s %s %s %d", m.Command, m.Key, m.Value, m.TTL))
	case GET, DEL:
		return []byte(fmt.Sprintf("%s %s", m.Command, m.Key))
	default:
		panic("unknown command")
	}
}

func parseMessage(raw []byte) (*Message, error) {
	rawStr := string(raw)
	parts := strings.Split(rawStr, " ")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid command: %s", rawStr)
	}

	cmd := Command(parts[0])

	msg := &Message{
		Command: cmd,
		Key:     []byte(parts[1]),
	}

	if cmd == SET {
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid SET command: %s", rawStr)
		}

		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("invalid SET TTL: %s", parts[3])
		}

		msg.Value = []byte(parts[2])
		msg.TTL = time.Duration(ttl)
	}

	return msg, nil
}
