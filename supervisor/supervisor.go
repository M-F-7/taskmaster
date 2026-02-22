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