package scanner

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/lewington/erosai-server/shared"
)

func TestScanner(t *testing.T) {
	in := make(chan shared.Link, 2)
	s := New(in)

	fmt.Println(s)

	in <- shared.Link{
		2,
		// "https://www.pornhub.com/",
		"https://www.ixxx.com/",
		// "https://www.instagram.com/caradelevingne/?hl=en",
		false,
		0,
	}

	time.Sleep(time.Second * 90)
}
