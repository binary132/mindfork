package message

import mfm "github.com/mindfork/mindfork/message"

const (
	TIntention mfm.Type = "intention"
	TEcho      mfm.Type = "echo"
	TSource    mfm.Type = "source"
)

type Echo struct{}
type Source struct{}
