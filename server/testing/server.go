package testing

import "github.com/mindfork/mindfork/message"

// Server is a rudimentary echo server.Server that records Messages received.
type Server struct {
	Messages []message.Message
}

// Serve implements server.Server on Server.
func (s *Server) Serve(m message.Message) message.Message {
	s.Messages = append(s.Messages, m)
	return m
}
