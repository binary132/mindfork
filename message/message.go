package message

// Type is a Message type constant.
type Type string

// Message encapsulates the different kinds of Mindfork messages.
type Message interface{}

// Encoder is a function for encoding a Message as a byte slice.
type Encoder func(Message) ([]byte, error)

// Decoder is a function for decoding a Message from a byte slice.
type Decoder func([]byte) (Message, error)

// Maker is a type having an Encoder and Decoder method.  These may be used to
// translate Messages between the wire protocol and Mindfork.  For a reference
// implementation, see core.Maker.
type Maker interface {
	Encoder() Encoder
	Decoder() Decoder
}
