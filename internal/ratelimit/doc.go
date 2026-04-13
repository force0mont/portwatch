// Package ratelimit provides a sliding-window rate limiter keyed by an
// arbitrary string (e.g. port+protocol pair).  Each key is tracked
// independently; once the burst limit is exceeded within the window the
// key is blocked until the window resets.
package ratelimit
