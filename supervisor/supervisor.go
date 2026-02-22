package supervisor

// Supervisor will manage all jobs.

import (
	"fmt"
	"taskmaster/config"
	"taskmaster/process"
)

type Supervisor struct {
	Prs map[string]*process.Process
	cfg *config.Config
}



func New(cfg *config.Config) *Supervisor {
    return &Supervisor{Prs: make(map[string]*process.Process), cfg: cfg}
}

func (s *Supervisor) Start() error{
	for name, prog := range(s.cfg.Programs){
		s.Prs[name] = process.New(prog.Cmd)
		if prog.AutoStart{
			err := s.Prs[name].Start()
			if err != nil{
				return err
			}
			go s.Prs[name].Wait()
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
	go s.Prs[name].Wait()
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
