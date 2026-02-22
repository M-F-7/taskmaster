package signals

// Signals handles Unix signal interception (SIGHUP, SIGTERM, etc.)

import (
	"os"
	"os/signal"
	// "strings"
	"syscall"
	"taskmaster/config"
	"taskmaster/logger"
	"taskmaster/process"
	"taskmaster/supervisor"
)

func Listen(svr *supervisor.Supervisor) error{
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM)

	for {
		sig := <-ch
		switch sig {
			case syscall.SIGHUP:
				newCfg, err := config.Load(svr.CfgPath)
				if err != nil{
					logger.Log(err.Error())
					break
				}
				svr.Reload(newCfg)
				// split := strings.Fields(newCfg.Programs.Cmd)
				// logger.LogRestart(split[0])
			case syscall.SIGTERM:
				for _, p := range svr.Prs {
    				if p.GetState() == process.Running {
        				p.Stop()
					}
				}
				os.Exit(0)
		}
	}
}
