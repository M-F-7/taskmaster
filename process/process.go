package process

// Process handles launching and monitoring a single child process.

import (
	// "fmt"
	// "log"
	// "os"
	"os/exec"
	"strings"
	"syscall"
	"taskmaster/logger"
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
	Name string
	Stopping bool
	Done chan error

}

func New(cmd string, name string) *Process {
    return &Process{cmd: cmd, state: Stopped, Name:name, Stopping: false, Done: make(chan error, 1)}
}


func (p* Process) Start() error{
	split := strings.Fields(p.cmd)
	p.exec = exec.Command(split[0], split[1:]...)
	err := p.exec.Start()
	if err != nil{
		return err
	}
	p.state = Running
	logger.LogStart(p.Name, p.Pid())
	return nil
}


func (p *Process) Wait() error {
    err := p.exec.Wait()
    wasStopping := p.Stopping
    p.state = Stopped
    p.Stopping = false
    if err != nil && !wasStopping {
        logger.LogDied(p.Name, p.exec.ProcessState.ExitCode())
		p.Done <- err
        return err
    }
	p.Done <- err
    return nil
}


func (p* Process) Stop() error{
	p.Stopping = true
	err := p.exec.Process.Signal(syscall.SIGTERM)
	if err != nil{
		return err
	}
	p.state = Stopped
	logger.LogStop(p.Name)
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