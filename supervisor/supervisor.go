package supervisor

// Supervisor will manage all jobs.

import (
	"fmt"
	"reflect"
	"taskmaster/config"
	"taskmaster/logger"
	"taskmaster/process"
	"time"
)

type Supervisor struct {
	Prs     map[string]*process.Process
	cfg     *config.Config
	CfgPath string
}

func New(cfg *config.Config, path string) *Supervisor {
	return &Supervisor{Prs: make(map[string]*process.Process), cfg: cfg, CfgPath: path}
}

func (s *Supervisor) Start() error {
	for name, prog := range s.cfg.Programs {
		// s.Prs[name] = process.New(prog.Cmd, name,prog.StopSignal, prog.StopTime, prog.Stdout, prog.Stderr, prog.Env, prog.WorkingDir, prog.Umask)
		var instanceName string
		for i := 0; i < prog.NumProcs; i++ {
		    instanceName = fmt.Sprintf("%s_%d", name, i)
		    s.Prs[instanceName] = process.New(prog.Cmd, instanceName,prog.StopSignal, prog.StopTime, prog.Stdout, prog.Stderr, prog.Env, prog.WorkingDir, prog.Umask)
			if prog.AutoStart {
				if instanceName != ""{
					err := s.Prs[instanceName].Start()
					if err != nil {
						return err
					}
					go s.Watch(instanceName)
				}
			}
		}
	}
	return nil
}

func (s *Supervisor) Status() {
	for name, p := range s.Prs {
		fmt.Printf("%-20s %-10s pid %d\n", name, p.GetState(), p.Pid())
	}
}

func (s *Supervisor) StartJob(name string) error {

	if _, ok := s.Prs[name]; !ok {
		fmt.Printf("unknown program: %s\n", name)
		return nil
	}
	if s.Prs[name].GetState() == process.Running {
		fmt.Printf("Process already running")
		return nil
	}
	err := s.Prs[name].Start()
	if err != nil {
		return err
	}
	go s.Watch(name)
	return nil

}

func (s *Supervisor) StopJob(name string) error {

	if _, ok := s.Prs[name]; !ok {
		fmt.Printf("unknown program: %s\n", name)
		return nil
	}
	if s.Prs[name].GetState() == process.Stopped {
		fmt.Printf("Process already stopped")
		return nil
	}
	err := s.Prs[name].Stop()
	if err != nil {
		return err
	}
	return nil

}

func (s *Supervisor) RestartJob(name string) error {
	if _, ok := s.Prs[name]; !ok {
		fmt.Printf("Unknown program: %s\n", name)
		return nil
	}
	select {
	case <-s.Prs[name].Done:
	default:
	}
	err := s.Prs[name].Stop()
	if err != nil {
		logger.Log(fmt.Sprintf("error stopping %s: %v", name, err))
	}
	// <-s.Prs[name].Done
	err = s.Prs[name].Start()
	if err != nil {
		return err
	}
	return nil
}

func (s *Supervisor) Reload(newCfg *config.Config) {
	var flag int

	for name, prog := range newCfg.Programs {
		flag = 0
		for curr_name, curr_prog := range s.cfg.Programs {
			if curr_name == name {
				flag = 1
				if reflect.DeepEqual(curr_prog, prog) { // ==
					break
				} else {
					s.RestartJob(name)
				}
			}
		}
		if flag == 0 {
			s.Prs[name] = process.New(prog.Cmd, name, prog.StopSignal, prog.StopTime, prog.Stdout, prog.Stderr, prog.Env, prog.WorkingDir, prog.Umask)

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

// func (s *Supervisor) Watch(name string) {
// 	retries := 0
// 	for {
// 		// _ = s.Prs[name].Wait()
// 		prog := s.cfg.Programs[name]
//         // lancer Wait() en arrière-plan
//         go s.Prs[name].Wait()
//         // phase 1 : starttime check
//         select {
//         case <-time.After(time.Duration(prog.StartTime) * time.Second):
//             retries = 0
//         case err := <-s.Prs[name].Done:
//             // mort trop tôt → gérer retries et autorestart
//             // ...
//             continue
//         }
//         // phase 2 : attendre la vraie mort
//         err := <-s.Prs[name].Done

// 		code := s.Prs[name].ExitCode()
// 		// prog := s.cfg.Programs[name]
		
// 		unexpected := true
// 		for _, expected := range prog.ExitCodes {
// 		    if code == expected {
// 		        unexpected = false
// 		        break
// 		    }
// 		}
// 		switch prog.AutoRestart {
// 		case "always":
// 			if retries >= prog.StartRetries {
//             	return
//         	}
//         	retries++
// 			logger.LogRestart(name)
// 			s.startOnly(name)
// 		case "unexpected":
// 			if unexpected {
// 				if retries >= prog.StartRetries {
//             		return
//         		}
//         		retries++
// 				logger.LogRestart(name)
// 				s.startOnly(name)
// 			} else {
// 				return
// 			}
// 		case "never":
// 			return
// 		}
// 	}
// }


func baseName(name string) string {
    for i := len(name) - 1; i >= 0; i-- {
        if name[i] == '_' {
            return name[:i]
        }
    }
    return name
}

func (s *Supervisor) Watch(name string) {
    retries := 0
    for {
		progName := baseName(name)
        prog := s.cfg.Programs[progName]
        // lancer Wait() en arrière-plan
        go s.Prs[name].Wait()
        // phase 1 : starttime check
        select {
        case <-time.After(time.Duration(prog.StartTime) * time.Second):
            retries = 0 // process stable, reset retries
        case err := <-s.Prs[name].Done:
			if s.shouldRestart(progName,name, err, &retries) {
        		logger.LogRestart(name)
        		s.startOnly(name)
    		} else {
        		return
		    }
            // continue
        }
	// default:
        // phase 2 : attendre la vraie mort
        err := <-s.Prs[name].Done
		if s.shouldRestart(progName,name, err, &retries) {
		    logger.LogRestart(name)
		    s.startOnly(name)
		} else {
		    return
		}
        _ = err
    }
}

func (s *Supervisor) shouldRestart(progName string, instanceName string, err error, retries *int) bool {
	_ = err 
	prog := s.cfg.Programs[progName]
    unexpected := true
    for _, expected := range prog.ExitCodes {
        if s.Prs[instanceName].ExitCode() == expected {
            unexpected = false
            break
        }
    }
    switch prog.AutoRestart {
    case "always":
        if *retries >= prog.StartRetries {
            return false
        }
        *retries++
        return true
    case "unexpected":
        if unexpected {
            if *retries >= prog.StartRetries {
                return false
            }
            *retries++
            return true
        }
        return false
    case "never":
        return false
    }
    return false
}

func (s *Supervisor) startOnly(name string) error {
	return s.Prs[name].Start()
}
