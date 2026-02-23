package shell

// Shell provides the interactive control shell for the user.

import (
	"fmt"
	"strings"
	"taskmaster/config"
	"taskmaster/supervisor"

	"github.com/chzyer/readline"
)

type Shell struct {
	svr *supervisor.Supervisor
}

func New(svr *supervisor.Supervisor) *Shell{
	return &Shell{svr: svr}
}

func (s* Shell) Run() error{
	rl, _ := readline.New("taskmaster> ")
	for {
		line, err := rl.Readline()
		if err != nil {break}
		
		split := strings.Fields(line)
		if len(split) == 0{
			continue
		}
		switch  split[0]{
		case "status":
			s.svr.Status()
		case "start":
			if len(split) < 2 {
        		fmt.Println("Usage: start <name>")
        		continue
    		}
			s.svr.StartJob(split[1])
		case "stop":
			if len(split) < 2 {
        		fmt.Println("Usage: stop <name>")
        		continue
    		}
			s.svr.StopJob(split[1])
		case "restart":
			if len(split) < 2 {
        		fmt.Println("Usage: restart <name>")
        		continue
    		}
			s.svr.RestartJob(split[1])
		case "reload":
			newcfg, err := config.Load(s.svr.CfgPath)
			if err != nil{
				fmt.Println("reload error:", err)
    			continue
			}
			s.svr.Reload(newcfg)
			
		case "exit":
			return nil
		default:
			fmt.Println("Unknown command")
			
		}
	}
	return nil
	//handle Eof
}
