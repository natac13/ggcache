package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/natac13/ggcache/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	// Embed the ServerOpts struct to inherit its methods.
	ServerOpts
	followers map[net.Conn]struct{}
	cache     cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		followers:  make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]", s.ListenAddr)

	if !s.IsLeader {
		go func() {
			conn, err := net.Dial("tcp", s.LeaderAddr)
			fmt.Println("connected with leader: ", s.LeaderAddr)
			if err != nil {
				log.Fatal(err)
			}

			s.handleConn(conn)
		}()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			// continue to accept new connections
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	if s.IsLeader {
		s.followers[conn] = struct{}{}
	}

	fmt.Println("connected with client: ", conn.RemoteAddr())
	for {
		// read from the connection
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("connection closed by %s\n", conn.RemoteAddr())
			} else {
				log.Printf("read error: %s\n", err)
			}
			break
		}

		go s.handleCommand(conn, buf[:n])
	}

}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)

	if err != nil {
		log.Printf("failed to parse command error: %s\n", err)
		conn.Write([]byte(err.Error()))
		return
	}

	switch msg.Command {
	case SET:
		err = s.handleSetCmd(conn, msg)
	case GET:
		err = s.handleGetCmd(conn, msg)
	case DEL:
		err = s.handleDelCmd(conn, msg)
	}

	if err != nil {
		log.Printf("failed to handle command error: %s\n", err)
		conn.Write([]byte(err.Error()))
	}

}

// handleSetCmd handles the SET command.
func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return fmt.Errorf("failed to set key: %s", err)
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

// sendToFollowers sends the message to all followers.
func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		fmt.Println("sending to follower: ", conn.RemoteAddr())
		_, err := conn.Write(msg.ToBytes())
		if err != nil {
			log.Printf("failed to write to follower: %s", err)
			continue
		}
	}
	return nil
}

// handleGetCmd handles the GET command.
func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.Key)
	if err != nil {
		return fmt.Errorf("failed to get key: %s", err)
	}

	_, err = conn.Write(val)
	if err != nil {
		return fmt.Errorf("failed to write response: %s", err)
	}

	return nil
}

func (s *Server) handleDelCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Delete(msg.Key); err != nil {
		return fmt.Errorf("failed to delete key: %s", err)
	}

	return nil
}
