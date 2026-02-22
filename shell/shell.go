package shell

// Shell provides the interactive control shell for the user.

import (
	"fmt"
	"taskmaster/supervisor"
	"github.com/chzyer/readline"
	"strings"
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
			;
		case "stop":
			;
		case "restart":
			;
		case "exit":
			return nil
		default:
			fmt.Println("Unknown command")
			
		}
	}
	return nil
	//handle Eof
}
