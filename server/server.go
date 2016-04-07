package server

import "github.com/mindfork/mindfork/message"

// Server must be implemented in order to wire up a Mindfork service.
type Server interface {
	// Serve specifies the routing and responses of Messages.
	Serve(m message.Message) message.Message
}
