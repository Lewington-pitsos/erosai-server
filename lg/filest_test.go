package lg

import (
	"io/ioutil"
	"testing"
)

func TestOldLogRefresh(t *testing.T) {
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

	refreshDir(oldLogDir())
	oldLogDirPath = tmp
	specialTest = false
}
