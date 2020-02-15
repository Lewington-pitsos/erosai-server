package lg

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestLoggerFiles(t *testing.T) {
	specialTest = true
	tmp := oldLogDirPath
	oldLogDirPath = "/src/bitbucket.org/lewington/autoroller/old-log-tst/"
	refreshDir(oldLogDir())

	oldLogFiles, err := ioutil.ReadDir(oldLogDir())
	if err != nil {
		panic(err)
	}
	oldLogFileCount := len(oldLogFiles)
	if oldLogFileCount != 1 {
		t.Fatalf("expected there to be 0 old logs initially, got %v", oldLogFileCount)
	}

	_ = NewLggr()

	oldLogFiles, err = ioutil.ReadDir(oldLogDir())
	if err != nil {
		panic(err)
	}
	oldLogFileCount = len(oldLogFiles)
	if 2 != oldLogFileCount {
		t.Fatalf("expected 2 files after creating a new logger, got %v", oldLogFileCount)
	}

	for i := 0; i < 7; i++ {
		time.Sleep(time.Millisecond * 300)
		_ = NewLggr()
	}

	oldLogFiles, err = ioutil.ReadDir(oldLogDir())
	if err != nil {
		panic(err)
	}

	oldLogFileCount = len(oldLogFiles)
	if oldLogFileCount > 6 {
		t.Fatalf("expected log files to start getting removed so there are never more than 5 log files (6 files in total), got %v", oldLogFileCount)
	}

	refreshDir(oldLogDir())
	oldLogDirPath = tmp
	specialTest = false

}
