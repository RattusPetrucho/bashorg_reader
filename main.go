package main

import (
	"github.com/jroimartin/gocui"
	"log"
)

func init() {
	var err error
	config, err = readConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error

	g := gocui.NewGui()

	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.FgColor = gocui.ColorDefault
	g.BgColor = gocui.ColorDefault
	g.ShowCursor = false

	err = g.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
