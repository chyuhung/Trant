package main

import (
	"Trant/config"
	"Trant/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zserge/lorca"
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
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	return exitChan
}
