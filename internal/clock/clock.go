// Package clock provides a simple wall-clock abstraction used throughout
// portwatch to allow deterministic time injection in tests.
package clock

import "time"

// Clock returns the current time.
type Clock func() time.Time

// Real is a Clock backed by time.Now.
var Real Clock = time.Now

// Fixed returns a Clock that always returns t.
func Fixed(t time.Time) Clock {
	return func() time.Time { return t }
}

// Advance returns a new Clock whose base time is t advanced by d.
func Advance(base Clock, d time.Duration) Clock {
	t := base().Add(d)
	return Fixed(t)
}
