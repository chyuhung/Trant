package main

import (
	"Trant/config"
	"Trant/server"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	chBrowerDead := make(chan struct{})
	chServerDead := make(chan struct{})
	go server.Run()
	go startBrowser(chBrowerDead, chServerDead)
	chSignal := listenInterrupt()
	for {
		select {
		case <-chSignal:
			chServerDead <- struct{}{}
		case <-chBrowerDead:
			os.Exit(0)
		}
	}
}

func startBrowser(chChromeDead chan struct{}, chServerDead chan struct{}) {
	port := config.GetPort()
	url := "http://127.0.0.1:" + port + "/static/index.html"
	chromePath := "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	cmd := exec.Command(chromePath, "--app="+url)
	err := cmd.Start()
	if err != nil {
		fmt.Println("浏览器启动失败，请手动访问：", url)
		fmt.Println(err)
	}
	go func() {
		<-chServerDead
		cmd.Process.Kill()
	}()
	go func() {
		cmd.Wait()
		chChromeDead <- struct{}{}
	}()

}
func listenInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
