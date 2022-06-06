package main

import "github.com/zserge/lorca"

func main() {
	var ui lorca.UI
	ui, _ = lorca.New("https://www.baidu.com", "", 1000, 800)
	select {}
	ui.Close()
}
