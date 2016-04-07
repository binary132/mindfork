package message_test

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Test hooks up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MessageSuite struct{}

var _ = Suite(&MessageSuite{})

func sampleMessages(s string) string {
	return map[string]string{
		"emptyMessage": ``,
		"emptyObject":  `{}`,
		"emptyString":  `""`,
		"validIntention": `
{"Who": "Bodie", "What": "To seek the Holy Grail"}`[1:],
		"timedIntention": `
{"Who": "Bodie", "What": "Something neat", "When": "2009-11-10T23:00:00Z"}`[1:],
	}[s]
}
