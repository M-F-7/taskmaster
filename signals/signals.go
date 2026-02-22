package signals

// Signals handles Unix signal interception (SIGHUP, SIGTERM, etc.)


import (
	"taskmaster/supervisor"
	"taskmaster/process"
	"os"
	"os/signal"
	"syscall"
	"taskmaster/logger"
	"taskmaster/config"
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
