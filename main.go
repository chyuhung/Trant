package main

import (
	"Trant/server"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

func main() {

	go server.Run()
	cmd := startBrowser()
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt)
	// 监听中断信号
	select {
	case <-chanSignal:
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("浏览器关闭失败，请手动关闭！")
			fmt.Println(err)
		}
	}
}

func startBrowser() *exec.Cmd {
	port := "27149"
	url := "http://127.0.0.1:" + port + "/static/index.html"
	chromePath := "C:\\Program Files (x86)\\Microsoft\\EdgeCore\\102.0.1245.30\\msedge.exe"
	cmd := exec.Command(chromePath, "--app="+url)
	err := cmd.Start()
	if err != nil {
		fmt.Println("浏览器启动失败，请手动访问：", url)
		fmt.Println(err)
	}
	return cmd
}
