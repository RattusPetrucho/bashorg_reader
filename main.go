package main

import (
	"errors"
	"fmt"
	"os"
	// get запрос
	"net/http"
	// парсиннг html
	"golang.org/x/net/html"

    // for config
    "encoding/json"

	// Перевод из cp1251 в utf-8
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"

    "github.com/jroimartin/gocui"
    "log"
)

type configuration struct {
    Cit int64
    Cnt int64
}

var config = new(configuration)

// Читаем цитату по номеру num
// и выбираем текст цитаты
func read_qoute(num int64) (string, error) {
	// Читаем цитату по номеру num
	url := "http://bash.im/quote/" + fmt.Sprintf("%d", num)
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body := html.NewTokenizer(resp.Body)

	// В цикле пробегаем по всем токенам html
	for {
		// берём следующий токен и получаем его тип
		tt := body.Next()
		// Проверяем тип токена
		switch tt {
		case html.ErrorToken:
			err := errors.New("no citations")
			return "", err
		// Если это тэг
		case html.StartTagToken:
			tn, attr := body.TagName()
			// тэг div с атрибутом
			if string(tn) == "div" && attr {
				// атрибут class со значением text
				key, val, _ := body.TagAttr()
				if string(key) == "class" && string(val) == "text" {
					// то извлекаем цитату
                    var s string;
                    for {
                        tt = body.Next()
                        sr := strings.NewReader(string(body.Text()))

                        // перекодируем цитату из cp1251 в utf-8
                        tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
                        buf, err := ioutil.ReadAll(tr)

                        if err != nil {
                            return "", err
                        }

                        s += string(buf)

                        tt = body.Next()

                        tn, _ := body.TagName()
                        if string(tn) == "br" {
                            s += "\n"
                        } else {
                            return s, nil
                        }
                    }
				}
			}
		}
	}
}

func read_bashorg(start, count int64) (string, error) {

    var out_str string;

    for i := start; i < count+start; i++ {
        str, err := read_qoute(i)
        if err != nil {
            return "", err
        }

        out_str += fmt.Sprintf("#%d\n%s", i, str)
        out_str += "\n\n"
    }

    return out_str, nil
}


// инициализация
// читаем конфигурацию из файле
func init_reader() (*configuration) {
    conf := configuration{}

    if _, err := os.Stat(os.Getenv("HOME") + "/bashor_reader"); err == nil {
        file, err := os.Open(os.Getenv("HOME") + "/bashor_reader")
        if err != nil {
            fmt.Println("error:", err)
            os.Exit(1)
        }
        defer file.Close()

        decoder := json.NewDecoder(file)
        
        err = decoder.Decode(&conf)

        if err != nil {
            fmt.Println("error:", err)
            os.Exit(1)
        }
    } else {
        conf.Cit = 1
        conf.Cnt = 7
    }

    return &conf
}

func save_config() {
    file, err := os.Create(os.Getenv("HOME") + "/bashor_reader")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)

    err = encoder.Encode(config)
    if err != nil {
        panic(err)
    }
}

// Получить следующие цитаты
func next_quote(g *gocui.Gui, v *gocui.View) error {
    v.Clear()

    config.Cit += config.Cnt

    str, err := read_bashorg(config.Cit, config.Cnt)
    v.Clear()
    if err != nil {
        fmt.Fprintf(v, "%s", err)
        return err
    } else {
        fmt.Fprintf(v, "%s", str)
    }

    iv, _ := g.View("info")
    iv.Clear()

    fmt.Fprintf(iv, "%s", "Готово!")

    return nil
}

func next_quote_info(g *gocui.Gui, v *gocui.View) error {
    iv, _ := g.View("info")

    iv.Clear()

    yopt := "Загрузка цитат!!!"

    fmt.Fprintf(iv, "%s", yopt)

    g.Flush()

    return nil
}

// Предыдущие цитаты
func prev_quote(g *gocui.Gui, v *gocui.View) error {
    v.Clear()

    if config.Cit - config.Cnt < 0 {
        config.Cit = 1
    } else {
        config.Cit -= config.Cnt
    }

    str, err := read_bashorg(config.Cit, config.Cnt)
    v.Clear()
    if err != nil {
        fmt.Fprintf(v, "%s", err)
        return err
    } else {
        fmt.Fprintf(v, "%s", str)
    }

    iv, _ := g.View("info")
    iv.Clear()

    fmt.Fprintf(iv, "%s", "Готово!")

    return nil
}

func prev_quote_info(g *gocui.Gui, v *gocui.View) error {
    iv, _ := g.View("info")

    iv.Clear()

    yopt := "Загрузка цитат!!!"

    fmt.Fprintf(iv, "%s", yopt)

    g.Flush()

    return nil
}

// выйти и сохранить текущие цитаты
func quit(g *gocui.Gui, v *gocui.View) error {
    config.Cit += config.Cnt
    save_config()
    return gocui.Quit
}

func keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyF10, gocui.ModNone, quit); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, next_quote_info); err != nil {
        return err
    }

    if err := g.SetKeybinding("main", gocui.KeyArrowRight, gocui.ModNone, next_quote); err != nil {
        return err
    }

    if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, prev_quote_info); err != nil {
        return err
    }

    if err := g.SetKeybinding("main", gocui.KeyArrowLeft, gocui.ModNone, prev_quote); err != nil {
        return err
    }

    return nil
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()

    if v, err := g.SetView("info", 0, 0, maxX-1, 2); err != nil {
        if err != gocui.ErrorUnkView {
            return err
        }

        yopt := "yopt!!!|||"

        fmt.Fprintf(v, "%s", yopt)

        v.Frame = true
    }

    if v, err := g.SetView("main", 0, 2, maxX-1, maxY-1); err != nil {
        if err != gocui.ErrorUnkView {
            return err
        }

        yopt := "Загрузка цитат!!!"

        fmt.Fprintf(v, "%s", yopt)

        str, err := read_bashorg(config.Cit, config.Cnt)
        v.Clear()
        if err != nil {
            fmt.Fprintf(v, "%s", err)
        } else {
            fmt.Fprintf(v, "%s", str)
        }

        v.Frame = true
        v.Wrap = true
        if err := g.SetCurrentView("main"); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    var err error
    config = init_reader()

    // defer save_config()

    // str := read_bashorg(config.Cit, config.Cnt)

    // fmt.Println(str)

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
