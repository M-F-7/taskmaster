package process

// Process handles launching and monitoring a single child process.

import (
	// "fmt"s
	// "log"
	// "os"
	"os/exec"
	"strings"
	"syscall"
)

type State int 

const (
	Running = iota //automatic incrementation
	Stopped
)

func (s State) String() string {
    switch s {
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


func (p* Process) Stop() error{
	err := p.exec.Process.Signal(syscall.SIGTERM)
	if err != nil{
		return err
	}
	p.state = Stopped
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

func (p* Process) SetState(state State) {
	p.state = state
}