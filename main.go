package main

import (
	"Trant/config"
	"Trant/server"
	"github.com/zserge/lorca"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	chBrowserExit := make(chan struct{})
	chServerExit := make(chan struct{})
	go server.Run()
	go startBrowser(chBrowserExit, chServerExit)
	chExit := listenInterrupt()
	for {
		select {
		case <-chExit:
			chServerExit <- struct{}{}
		case <-chBrowserExit:
			os.Exit(0)
		}
	}
}

func startBrowser(chBrowserExit chan struct{}, chServerExit chan struct{}) {
	port := config.GetPort()
	url := "http://127.0.0.1:" + port + "/static/index.html"
	ui, err := lorca.New(url, "", 600, 600)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-chServerExit
		ui.Close()
	}()
	go func() {
		<-ui.Done()
		chBrowserExit <- struct{}{}
	}()

}
func listenInterrupt() chan os.Signal {
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	return exitChan
}
