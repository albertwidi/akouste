package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type arrayFlags []string

func (af *arrayFlags) String() string {
	return ""
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

// Flags of application operator
type Flags struct {
	FilesToWatch arrayFlags
}

func main() {
	appFlags := Flags{}

	flag.Parse()
	flag.Var(&appFlags.FilesToWatch, "fw", "watch the file to change")

	log.Println("waiting for changes")

	signalReload := make(chan os.Signal, 1024)
	signal.Notify(signalReload, syscall.SIGUSR2)
	go func(sig chan os.Signal) {
		for {
			<-sig
			log.Println("for reloading")
		}
	}(signalReload)

	// exit when receive signal
	signalCh := make(chan os.Signal, 1024)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case sig := <-signalCh:
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("exiting operator")
			os.Exit(1)

		default:
			log.Println("UNKNOWN SIGNAL")
		}
	}
}
