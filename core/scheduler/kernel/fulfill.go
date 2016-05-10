package kernel

import (
	"errors"

	"github.com/mindfork/mindfork/core/message"
	mfm "github.com/mindfork/mindfork/message"
)

// Fulfill implements Scheduler.Fulfill on Kernel.
func (k *Kernel) Fulfill(f message.Fulfillment) mfm.Message {
	return errors.New("implement me")
}
