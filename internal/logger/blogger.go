package logger

import (
	"log"
	"sync"
)

//BLogger Object
type BLogger struct {
}

var blogger *BLogger
var blSync sync.Once

//GetLog return native logger
func GetLog() *BLogger {
	blSync.Do(func() {
		blogger = &BLogger{}
	})
	return blogger
}

//Console print to console
func (BLogger *BLogger) Console(str string) {
	log.SetFlags(0)
	log.Print(str)
}

//Table print to console
func (BLogger *BLogger) Table(v interface{}) {
	log.SetFlags(0)
	log.Print(v)
}
