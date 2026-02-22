package logger

// Logger handles writing events to a log file.

import (
	"fmt"
	"os"
	"log"
)

var logger *log.Logger

func Init(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil{
		return err
	}
	logger = log.New(f, "", log.LstdFlags)
	return nil
}


func Log(msg string){
	logger.Println(msg)
}


func LogStart(name string, pid int){
	Log(fmt.Sprintf("[%s] started (pid %d)", name, pid))
}

func LogStop(name string) {
	Log(fmt.Sprintf("[%s] stopped ", name))
}
func LogDied(name string, exitCode int) {
	Log(fmt.Sprintf("[%s] died unexpectedly (exit code %d) ", name, exitCode))
}
func LogRestart(name string) {
	Log(fmt.Sprintf("[%s] restarting", name))
}
func LogReload() {
	Log("config reloaded")
}