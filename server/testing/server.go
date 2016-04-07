package testing

import "github.com/mindfork/mindfork/message"

// Server is a rudimentary echo server.Server.
type Server struct{}

// Serve implements server.Server on Server.
func (s *Server) Serve(m message.Message) message.Message { return m }
