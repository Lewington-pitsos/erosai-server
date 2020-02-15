package assist

import (
	"testing"
	"time"
)

func TestGeneralTime(t *testing.T) {
	unixTime := 1554340278
	currentTime := time.Unix(int64(unixTime), 0)

	betfairTime := ToBetfair(currentTime)

	if unixTime != int(betfairTime.Unix()) {
		t.Fatalf("expected betfair timestamp to be %v, got %v", unixTime, betfairTime.Unix())
	}

	if betfairTime.Hour() != 1 {
		t.Fatalf("expected betfair hout to be %v, got %v", 1, betfairTime.Hour())
	}

	if betfairTime.Minute() != 11 {
		t.Fatalf("expected betfair minute to be %v, got %v", 1, betfairTime.Minute())
	}

	dayEnd := EndOfDay(currentTime)

	if unixTime == int(dayEnd.Unix()) {
		t.Fatalf("expected betfair timestamp to be %v, got %v", unixTime, dayEnd.Unix())
	}

	if dayEnd.Hour() != 23 {
		t.Fatalf("expected betfair hour to be %v, got %v", 23, dayEnd.Hour())
	}

	if dayEnd.Minute() != 59 {
		t.Fatalf("expected betfair minute to be %v, got %v", 59, dayEnd.Minute())
	}

}

func TestGeneralTimeSameDay(t *testing.T) {
	day, err := time.Parse("2006-01-02 15:04:05 MST", "2019-04-06 08:45:00 AEST")
	Check(err)
	gmtDay, err := time.Parse("2006-01-02 15:04:05 MST", "2019-04-05 22:45:00 GMT")
	Check(err)

	if !SameDay(ToBetfair(day), gmtDay) {
		t.Fatalf("expected AEST and GMT dates %v and %v to be on the same day since they actually refer to the same time", day, gmtDay)
	}
}

func TestGeneralMath(t *testing.T) {
	actual := Rounded(67, 50)
	rounded := 50
	if actual != rounded {
		t.Fatalf("expected rounded number to be %v, got %v", rounded, actual)
	}
	actual = Rounded(50, 50)
	rounded = 50
	if actual != rounded {
		t.Fatalf("expected rounded number to be %v, got %v", rounded, actual)
	}
	actual = Rounded(151, 50)
	rounded = 150
	if actual != rounded {
		t.Fatalf("expected rounded number to be %v, got %v", rounded, actual)
	}
	actual = Rounded(149, 50)
	rounded = 100
	if actual != rounded {
		t.Fatalf("expected rounded number to be %v, got %v", rounded, actual)
	}
	actual = Rounded(100, 50)
	rounded = 100
	if actual != rounded {
		t.Fatalf("expected rounded number to be %v, got %v", rounded, actual)
	}

	actual = FloatRepDivide(1000, 100)
	divided := 1000
	if actual != divided {
		t.Fatalf("expected divided number to be %v, got %v", divided, actual)
	}
	actual = FloatRepDivide(1000, 1000)
	divided = 100
	if actual != divided {
		t.Fatalf("expected divided number to be %v, got %v", divided, actual)
	}
	actual = FloatRepDivide(10000, 400)
	divided = 2500
	if actual != divided {
		t.Fatalf("expected divided number to be %v, got %v", divided, actual)
	}
	actual = FloatRepDivide(1000, 300)
	divided = 333
	if actual != divided {
		t.Fatalf("expected divided number to be %v, got %v", divided, actual)
	}
	actual = FloatRepDivide(1500, 400)
	divided = 375
	if actual != divided {
		t.Fatalf("expected divided number to be %v, got %v", divided, actual)
	}
	actual = FloatRepDivide(1000, 1500)
	divided = 66
	if actual != divided {
		t.Fatalf("expected divided number to be %v, got %v", divided, actual)
	}
	actual = FloatRepDivide(1000, 2000)
	divided = 50
	if actual != divided {
		t.Fatalf("expected divided number to be %v, got %v", divided, actual)
	}
}

func TestGeneralTypecasting(t *testing.T) {
	actual := IntToNoDecimalNumber(100)
	rounded := "1"
	if actual != rounded {
		t.Fatalf("expected no decimal version to be %v, got %v", rounded, actual)
	}
	actual = IntToNoDecimalNumber(157)
	rounded = "1"
	if actual != rounded {
		t.Fatalf("expected no decimal version to be %v, got %v", rounded, actual)
	}
	actual = IntToNoDecimalNumber(15732)
	rounded = "157"
	if actual != rounded {
		t.Fatalf("expected no decimal version to be %v, got %v", rounded, actual)
	}
	actual = IntToNoDecimalNumber(1699)
	rounded = "16"
	if actual != rounded {
		t.Fatalf("expected no decimal version to be %v, got %v", rounded, actual)
	}
	actual = IntToNoDecimalNumber(122)
	rounded = "1"
	if actual != rounded {
		t.Fatalf("expected no decimal version to be %v, got %v", rounded, actual)
	}
	actual = IntToNoDecimalNumber(50)
	rounded = "0"
	if actual != rounded {
		t.Fatalf("expected no decimal version to be %v, got %v", rounded, actual)
	}
	actual = IntToNoDecimalNumber(0)
	rounded = "0"
	if actual != rounded {
		t.Fatalf("expected no decimal version to be %v, got %v", rounded, actual)
	}
}
