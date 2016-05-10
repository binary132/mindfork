package message

// Fulfillment is a Message which fulfills an Intention.
type Fulfillment struct {
	ID int64 `json:",omitempty"`

	// Which is the Intention which is fulfilled by this Fulfillment.
	Which int64 `json:",omitempty"`

	// Who is the user which created the Fulfillment.
	Who string

	// What is the value which the Intention wants.
	What interface{} `json:",omitempty"`
}
