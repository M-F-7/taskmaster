package process

// Process handles launching and monitoring a single child process.

import (
	// "fmt"s
	// "log"
	// "os"
	"strings"
	"os/exec"
)

type State int 

const (
	Started State = iota //automatic incrementation
	Running
	Stopped
)

func (s State) String() string {
    switch s {
    case Started:
        return "STARTING"
    case Running:
        return "RUNNING"
    case Stopped:
        return "STOPPED"
    default:
        return "UNKNOWN"
    }
}

type Process struct {
	cmd string
	state State
	exec *exec.Cmd

}

func New(cmd string) *Process {
    return &Process{cmd: cmd, state: Stopped}
}


func (p* Process) Start() error{
	split := strings.Fields(p.cmd)
	p.exec = exec.Command(split[0], split[1:]...)

	err := p.exec.Start()
	if err != nil{
		return err
	}
	p.state = Running
	return nil
}


func (p* Process) Wait() error{
	err := p.exec.Wait()
	p.state = Stopped
	if err != nil{
		return err
	}
	return nil
}


func (p* Process) Pid () int{
	if p.exec == nil || p.exec.Process == nil {
        return 0
    }
	return p.exec.Process.Pid
}

func (p* Process) GetState() State{
	return p.state
}