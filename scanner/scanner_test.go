package scanner

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/lewington/erosai/shared"
)

func TestScanner(t *testing.T) {
	in := make(chan shared.Link, 2)
	s := New(in)

	fmt.Println(s)

	in <- shared.Link{
		2,
		// "https://www.pornhub.com/",
		"https://www.ixxx.com/",
		false,
		0,
	}

	time.Sleep(time.Second * 90)
}
