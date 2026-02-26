package process

// Process handles launching and monitoring a single child process.

import (
	// "fmt"
	// "log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"taskmaster/logger"
	"time"
)

type State int 

const (
	Running = iota //automatic incrementation
	Stopped
)

var signals = map[string]syscall.Signal{
    "TERM": syscall.SIGTERM,
    "KILL": syscall.SIGKILL,
    "USR1": syscall.SIGUSR1,
    "USR2": syscall.SIGUSR2,
    "HUP":  syscall.SIGHUP,
    "INT":  syscall.SIGINT,
}

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
	stopSignal string
	stopTime int
	stdout string
	stderr string  
	env map[string]string
	workingDir string
	umask string
}

func New(cmd string, name string, signal string, timesig int, stdout string, stderr string, env map[string] string, workingDir string, umask string) *Process {
    return &Process{
		cmd: cmd, 
		state: Stopped,
		Name:name,
		Stopping: false,
		Done: make(chan error,1),
		stopSignal: signal,
		stopTime:timesig,
		stdout: stdout,
		stderr: stderr,
		env: env,
		workingDir: workingDir,
		umask: umask, //other file permissions 
}
}


func (p* Process) Start() error{
	split := strings.Fields(p.cmd)
	p.exec = exec.Command(split[0], split[1:]...)
	p.exec.Dir = p.workingDir
	for key, value := range p.env {
    	p.exec.Env = append(p.exec.Env, key+"="+value)
	}	
	if p.stdout != "" {
    	f, err := os.OpenFile(p.stdout, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    	if err != nil {
    	    return err
    	}
    	p.exec.Stdout = f
	}
	if p.stderr != "" {
    	f, err := os.OpenFile(p.stderr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    	if err != nil {
    	    return err
    	}
    	p.exec.Stderr = f
	}
	if p.umask != "" {
	    val, _ := strconv.ParseUint(p.umask, 8, 32)
	    old := syscall.Umask(int(val))
	    err := p.exec.Start()
	    syscall.Umask(old)
	    if err != nil {
	        return err
	    }
	} else {
	    err := p.exec.Start()
	    if err != nil {
	        return err
	    }
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
        logger.LogDied(p.Name, p.exec.ProcessState.ExitCode(), p.stopSignal)
		p.Done <- err
        return err
    }
	p.Done <- err
    return nil
}


func (p* Process) Stop() error{
	p.Stopping = true
	sig := signals[p.stopSignal]
	err := p.exec.Process.Signal(sig)
	// time.Sleep(time.Duration(p.stopTime) * time.Second)
	select {
		case <-p.Done:
		case <-time.After(time.Duration(p.stopTime) * time.Second):
		    p.exec.Process.Kill()  
	}
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

func (p *Process) SetStopSignal(sig string) {
    p.stopSignal = sig
}


func (p *Process) ExitCode() int {
    return p.exec.ProcessState.ExitCode()
}