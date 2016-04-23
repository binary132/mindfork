package scheduler_test

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type SchedulerSuite struct{}

var _ = Suite(&SchedulerSuite{})
