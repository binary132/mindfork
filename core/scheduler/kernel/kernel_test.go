package kernel_test

import (
	"github.com/mindfork/mindfork/core/scheduler"
	"github.com/mindfork/mindfork/core/scheduler/kernel"
)

var _ = scheduler.Scheduler(&kernel.Kernel{})
