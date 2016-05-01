package message

import (
	"time"

	mfm "github.com/mindfork/mindfork/message"
)

// Core Message types.
const (
	TIntention mfm.Type = "intention"
	TEcho      mfm.Type = "echo"
	TSource    mfm.Type = "source"
)

// Echo contains a time which should be the time the message was created.
type Echo struct {
	When time.Time
}

// Source contains the core's source reference and license.
type Source struct{}
