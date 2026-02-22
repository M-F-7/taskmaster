package main

import (
	"fmt"
	"os"
	"taskmaster/config"
	// "taskmaster/logger"
	// "taskmaster/process"
	"taskmaster/supervisor"
	"taskmaster/shell"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config.yml>\n", os.Args[0])
		os.Exit(1)
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("Loaded %d program(s):\n", len(cfg.Programs))
	// for name, prog := range cfg.Programs {
	// 	fmt.Printf("  [%s] cmd=%q numprocs=%d autostart=%v exitcodes=%v\n",
	// 		name, prog.Cmd, prog.NumProcs, prog.AutoStart, prog.ExitCodes)
	// }

	spr := supervisor.New(cfg)
	spr.Start()

	shl := shell.New(spr)
	shl.Run()
	
}
