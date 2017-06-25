package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

var symafore = make(chan struct{}, 1)

// Получить следующие цитаты
func nextQuote(g *gocui.Gui, v *gocui.View) error {
	iv, err := g.View("info")
	if err != nil {
		return err
	}

	iv.Clear()

	fmt.Fprintf(iv, "%s", "Загрузка цитат!!!")

	go g.Execute(func(g *gocui.Gui) error {
		defer func() { <-symafore }()
		symafore <- struct{}{}

		v, err := g.View("main")
		if err != nil {
			return err
		}
		iv, err := g.View("info")
		if err != nil {
			return err
		}

		v.Clear()
		iv.Clear()

		config.Cit += config.Cnt

		str, err := readBashorg(config.Cit, config.Cnt)
		v.Clear()
		if err != nil {
			fmt.Fprintf(v, "%s", err)
			return err
		} else {
			fmt.Fprintf(v, "%s", str)
		}

		iv.Clear()

		fmt.Fprintf(iv, "%s", "Готово!")

		return nil
	})

	return nil
}

// Предыдущие цитаты
func prevQuote(g *gocui.Gui, v *gocui.View) error {
	iv, err := g.View("info")
	if err != nil {
		return err
	}

	iv.Clear()

	fmt.Fprintf(iv, "%s", "Загрузка цитат!!!")

	go g.Execute(func(g *gocui.Gui) error {
		defer func() { <-symafore }()
		symafore <- struct{}{}

		v, err := g.View("main")
		if err != nil {
			return err
		}
		iv, err := g.View("info")
		if err != nil {
			return err
		}

		v.Clear()
		iv.Clear()

		if config.Cit-config.Cnt < 0 {
			config.Cit = 1
		} else {
			config.Cit -= config.Cnt
		}

		str, err := readBashorg(config.Cit, config.Cnt)
		v.Clear()
		if err != nil {
			fmt.Fprintf(v, "%s", err)
			return err
		} else {
			fmt.Fprintf(v, "%s", str)
		}

		iv.Clear()

		fmt.Fprintf(iv, "%s", "Готово!")

		return nil
	})

	return nil
}

// выйти и сохранить текущие цитаты
func quit(g *gocui.Gui, v *gocui.View) error {
	config.Cit += config.Cnt
	saveConfig()
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyF10, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, nextQuote); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, prevQuote); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("info", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		yopt := "yopt!!!|||"

		fmt.Fprintf(v, "%s", yopt)

		v.Frame = true
	}

	if v, err := g.SetView("main", 0, 3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		yopt := "Загрузка цитат!!!"

		fmt.Fprintf(v, "%s", yopt)

		str, err := readBashorg(config.Cit, config.Cnt)
		v.Clear()
		if err != nil {
			fmt.Fprintf(v, "%s", err)
		} else {
			fmt.Fprintf(v, "%s", str)
		}

		v.Frame = true
		v.Wrap = true
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}

	return nil
}
