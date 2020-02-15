package auth

import (
	"testing"
	"time"
)

func TestTimeThrottle(t *testing.T) {
	timeout := time.Millisecond * 100
	th := NewTimeThrottle(3, timeout)
	if !th.Allow() {
		t.Fatal("expected timethrottle to allow the first request")
	}

	th.Fail()
	th.Fail()
	th.Fail()

	if th.Allow() {
		t.Fatal("expected timethrottle not to allow requests just after 3 fails")
	}

	time.Sleep(timeout * time.Duration(2))

	if !th.Allow() {
		t.Fatal("expected timethrottle not to allow requests after timeout finished")
	}

	th.Fail()
	th.Fail()
	th.Fail()

	if th.Allow() {
		t.Fatal("expected timethrottle not to allow requests just after 3 fails")
	}

	th.Fail()
	th.Fail()

	time.Sleep(timeout * time.Duration(2))

	if !th.Allow() {
		t.Fatal("expected timethrottle not to allow requests after timeout finished, even if fails occured during the timeout.")
	}

	th.Succeed()
	th.Fail()

	if !th.Allow() {
		t.Fatal("expected timethrottle to reset fails after each success")
	}

	th.Fail()
	th.Fail()
	th.Fail()

	if th.Allow() {
		t.Fatal("expected timethrottle not to allow requests just after 3 fails")
	}
}
