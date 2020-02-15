package assist

import "testing"

func TestRandom(t *testing.T) {
	for i := 0; i < 100; i++ {
		x := RandomBetween(1, 10)
		if x > 10 || x < 1 {
			t.Fatalf("result %v was not between the required range", x)
		}
	}

	for i := 0; i < 100; i++ {
		x := RandomBetween(7, 10)
		if x > 10 || x < 7 {
			t.Fatalf("result %v was not between the required range", x)
		}
	}

	for i := 0; i < 100; i++ {
		x := RandomBetween(-10, -4)
		if x > -4 || x < -10 {
			t.Fatalf("result %v was not between the required range", x)
		}
	}
}
