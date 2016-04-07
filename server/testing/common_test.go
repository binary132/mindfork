package testing_test

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestingSuite struct{}

var _ = Suite(&TestingSuite{})
