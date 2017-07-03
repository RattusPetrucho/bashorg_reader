package main

import (
	"log"

	"github.com/jroimartin/gocui"
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

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	// g.SetLayout(layout)
	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.FgColor = gocui.ColorDefault
	g.BgColor = gocui.ColorDefault
	// g.ShowCursor = false
	g.Cursor = false

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
