package supervisor

// Supervisor will manage all jobs.

import (
	"fmt"
	"reflect"
	"taskmaster/config"
	"taskmaster/logger"
	"taskmaster/process"
)

type Supervisor struct {
	Prs map[string]*process.Process
	cfg *config.Config
	CfgPath string
}



func New(cfg *config.Config, path string) *Supervisor {
    return &Supervisor{Prs: make(map[string]*process.Process), cfg: cfg, CfgPath: path}
}

func (s *Supervisor) Start() error{
	for name, prog := range(s.cfg.Programs){
		s.Prs[name] = process.New(prog.Cmd, name)
		if prog.AutoStart{
			err := s.Prs[name].Start()
			if err != nil{
				return err
			}
			go s.Watch(name)
		}
	}
	return nil
}


func (s *Supervisor) Status() {
    for name, p := range s.Prs {
        fmt.Printf("%-20s %-10s pid %d\n", name, p.GetState(), p.Pid())
    }
}

func (s *Supervisor) StartJob(name string) error{

	if _, ok := s.Prs[name]; !ok {
    	fmt.Printf("unknown program: %s\n", name)
    	return nil
	}
	if s.Prs[name].GetState() == process.Running{
		fmt.Printf("Process already running")
		return nil
	}
	err := s.Prs[name].Start()
	if err != nil{
		return err
	}
	go s.Watch(name)
	return nil

}

func (s *Supervisor) StopJob(name string) error{

	if _, ok := s.Prs[name]; !ok {
    	fmt.Printf("unknown program: %s\n", name)
    	return nil
	}
	if s.Prs[name].GetState() == process.Stopped{
		fmt.Printf("Process already stopped")
		return nil
	}
	err := s.Prs[name].Stop()
	if err != nil{
		return err
	}
	return nil

}


func (s *Supervisor) RestartJob(name string) error{
	if _, ok := s.Prs[name]; !ok {
		fmt.Printf("Unknown program: %s\n", name)
		return nil
	}
	err := s.Prs[name].Stop()
	if err != nil {
    	logger.Log(fmt.Sprintf("error stopping %s: %v", name, err))
	}
	err = s.Prs[name].Start()
	if err != nil{
		return err
	}
	s.Watch(name)
	return nil
}


func (s *Supervisor) Reload(newCfg *config.Config) {
	var flag int

	for name, prog := range newCfg.Programs {
		flag = 0
		for curr_name, curr_prog := range s.cfg.Programs{
			if curr_name == name{
				flag = 1
				if reflect.DeepEqual(curr_prog, prog){ // ==
					break
				} else {
					s.RestartJob(name)
				}
			}
		}
		if flag == 0{
			s.Prs[name] = process.New(prog.Cmd, name)

			if prog.AutoStart {
        		s.StartJob(name)
	    	}
		}
	}

	for name, p := range s.Prs {
    	if _, exists := newCfg.Programs[name]; !exists {
        	p.Stop()
        	delete(s.Prs, name)
    	}
	}
    s.cfg = newCfg
}


func (s *Supervisor) Watch(name string) {
	s.Prs[name].Wait()
    err := <-s.Prs[name].Done
    prog := s.cfg.Programs[name]
    
    switch prog.AutoRestart {
    case "always":
		logger.LogRestart(name)
        s.StartJob(name)
        s.Watch(name)
    case "unexpected":
        if err != nil {
			logger.LogRestart(name)
            s.StartJob(name)
        	s.Watch(name)
        }
    case "never":
    }
}