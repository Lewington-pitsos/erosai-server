package lg

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const linesPerFile = 10000000

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

type Lggr struct {
	input      chan statement
	fileNumber int
	lines      int
	file       *os.File
}

func (l *Lggr) Debug(str string, args ...interface{}) {
	l.input <- statement{
		kind:    "DEBUG",
		message: fmt.Sprintf(str, args...),
	}
}

func (l *Lggr) Info(str string, args ...interface{}) {
	l.input <- statement{
		kind:    "INFO",
		message: fmt.Sprintf(str, args...),
	}
}

func (l *Lggr) Warn(str string, args ...interface{}) {
	l.input <- statement{
		kind:    "WARN",
		message: fmt.Sprintf(str, args...),
	}
}

func (l *Lggr) Error(str string, args ...interface{}) {
	l.input <- statement{
		kind:    "ERROR",
		message: fmt.Sprintf(str, args...),
	}
}

func (l *Lggr) Fatal(str string, args ...interface{}) {
	l.input <- statement{
		kind:    "FATAL",
		message: fmt.Sprintf(str, args...),
	}
}

func (l *Lggr) keepLogging() {
	for s := range l.input {
		l.lines++
		if l.lines > linesPerFile {
			l.startNewFile()
		}
		switch s.kind {
		case "DEBUG":
			log.Debug(s.message)
		case "INFO":
			log.Info(s.message)
		case "WARN":
			log.Warn(s.message)
		case "ERROR":
			log.Error(s.message)
		case "FATAL":
			log.Fatal(s.message)
		}
	}
}

func (l *Lggr) dirName() string {
	return build.Default.GOPATH + "/src/bitbucket.org/lewington/erosai-server/log/"
}

func (l *Lggr) fileName() string {
	return fmt.Sprintf("%vmain-%v.log", l.dirName(), l.fileNumber)
}

func (l *Lggr) SetLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
}

func (l *Lggr) LogGeneralStatsEvery(wait int) {
	go l.keepLoggingStats(wait)
}

func (l *Lggr) keepLoggingStats(wait int) {
	for {
		time.Sleep(time.Second * time.Duration(wait))
		l.generalStats()
	}
}

func (l *Lggr) generalStats() {
	l.MemUsage()
	l.OpenFiles()
}

func (l *Lggr) startNewFile() {
	if l.file != nil {
		l.file.Close()
	}
	l.lines = 0
	l.fileNumber++
	var file, err = os.OpenFile(l.fileName(), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	l.file = file
	log.SetOutput(file)
}

func (l *Lggr) MemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	l.Info("****** MEMORY USAGE *******")
	l.Info("Alloc = %v MiB", bToMb(m.Alloc))
	l.Info("HeapIdle = %v MiB", bToMb(m.HeapIdle))
	l.Info("HeapInuse = %v MiB", bToMb(m.HeapInuse))
	l.Info("HeapReleased = %v MiB", bToMb(m.HeapReleased))
	l.Info("Malloc = %v\n", bToMb(m.Mallocs))
	l.Info("Frees = %v\n", bToMb(m.Frees))
	l.Info("TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	l.Info("Sys = %v MiB", bToMb(m.Sys))
	l.Info("StackInUse = %v MiB", bToMb(m.StackInuse))
	l.Info("StackSys = %v MiB", bToMb(m.StackSys))
	// l.Info("LibxmlSize = %v B", mem.AllocSize())
}

func (l *Lggr) OpenFiles() {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		l.Error(err.Error())
	}
	lines := strings.Split(string(out), "\n")
	l.Info("%v Open Files", int64(len(lines)-1))
}

func NewLggr() *Lggr {
	l := &Lggr{
		make(chan statement, 200),
		0,
		0,
		nil,
	}
	if flag.Lookup("test.v") == nil || specialTest == true {
		os.Rename(l.dirName()+"main-1.log", oldLogDir()+fmt.Sprintf("%v.log", time.Now().Format("02-01-2006_15:04:05.000")))
		oldLogFiles, err := ioutil.ReadDir(oldLogDir())

		if err != nil {
			panic(err)
		}

		deleteOldestFiles(oldLogFiles)
	}

	refreshDir(l.dirName())

	log.SetLevel(log.DebugLevel)
	l.startNewFile()

	go l.keepLogging()

	return l
}

var L *Lggr = NewLggr()

func oldLogDir() string {
	return build.Default.GOPATH + oldLogDirPath
}
