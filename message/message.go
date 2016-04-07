package message

import "io"

// Type is a Message type constant.
type Type string

// Message encapsulates the different kinds of Mindfork messages.
type Message interface{}

// Encoder is a type which implements Encode; for example,
// encoding/json.Encoder.
type Encoder interface {
	// Encode writes the encoding of the given Message to the Writer.
	Encode(Message) error
}

// Decoder is a type which implements Decode; for example,
// encoding/json.Decoder.
type Decoder interface {
	// Decode reads the next encoded Message from its input and stores it in
	// the given object.
	Decode(Message) error
}

// MessageMaker is a type having a NewEncoder and NewDecoder method.  These
// may be used to translate Messages between the wire protocol and Mindfork.
// For a reference implementation, see core.MessageMaker.
type MessageMaker interface {
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}
