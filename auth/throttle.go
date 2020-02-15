package auth

// Throttle keeps track of the number of recent
// auth requests and decides whether we should allow
// another request right now.
type Throttle interface {
	Allow() bool
	Fail()
	Succeed()
}
