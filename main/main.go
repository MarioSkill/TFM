package main

import (
	"gitlab.com/marioskill/configuration"
	"gitlab.com/marioskill/global"
	"gitlab.com/marioskill/gui"
)

func main() {
	global.Init()
	gui.Start(configuration.GUI_port)
}
