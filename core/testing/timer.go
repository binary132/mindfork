package testing

import "time"

// TestTimer is a time.Time implementing core.Timer.  Its Now always returns its
// value.
type TestTimer time.Time

// Now implements core.Timer for TestTimer.
func (t TestTimer) Now() time.Time { return time.Time(t) }
