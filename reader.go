package main

import (
	"errors"
	// "fmt"
	"strconv"

	// парсиннг html
	// "golang.org/x/net/html"

	// get запрос
	"net/http"

	// Перевод из cp1251 в utf-8
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
	// "os"
)

type quote_res struct {
	Num   int
	Quote string
	Err   error
}

// var file *os.File

var quote_limit = make(chan struct{}, 10)

// Читаем цитату по номеру num
// и выбираем текст цитаты
func readQuote(num int, url string, quotes chan<- quote_res) {
	quote_limit <- struct{}{}
	defer func() { <-quote_limit }()

	qr := quote_res{Num: num}

	// Читаем цитату по номеру num
	// url := "http://bash.im/quote/" + fmt.Sprintf("%d", num)
	resp, err := http.Get(url)

	if err != nil {
		qr.Err = err
		quotes <- qr
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		qr.Err = err
		quotes <- qr
		return
	}

	for i := 0; i < len(buf); i++ {
		if buf[i] == 'd' && buf[i+1] == 'i' && buf[i+2] == 'v' {
			for buf[i] != '>' && buf[i] != 'c' {
				i++
				if i >= len(buf) {
					qr.Err = errors.New("no citations")
					quotes <- qr
					return
				}
			}
			if buf[i] == 'c' && buf[i+1] == 'l' && buf[i+2] == 'a' && buf[i+3] == 's' && buf[i+4] == 's' {
				i += 5
				if buf[i+2] == 't' && buf[i+3] == 'e' && buf[i+4] == 'x' && buf[i+5] == 't' {
					i += 8
					// start := i
					var buff []byte
					for string(buf[i:i+6]) != "</div>" {
						if buf[i] == '&' && buf[i+1] == 'l' && buf[i+2] == 't' && buf[i+3] == ';' {
							buff = append(buff, '<')
							i += 4
						} else if buf[i] == '&' && buf[i+1] == 'g' && buf[i+2] == 't' && buf[i+3] == ';' {
							buff = append(buff, '>')
							i += 4
						} else if buf[i] == '<' && buf[i+1] == 'b' && buf[i+2] == 'r' && buf[i+3] == '>' {
							buff = append(buff, '\n')
							i += 4
						} else if buf[i] == '<' && buf[i+1] == 'b' && buf[i+2] == 'r' && buf[i+3] == ' ' && buf[i+4] == '/' && buf[i+5] == '>' {
							buff = append(buff, '\n')
							i += 6
						} else if buf[i] == '&' && buf[i+1] == 'q' && buf[i+2] == 'u' && buf[i+3] == 'o' && buf[i+4] == 't' && buf[i+5] == ';' {
							buff = append(buff, '"')
							i += 6
						} else {
							buff = append(buff, buf[i])
							i++
						}
					}

					sr := strings.NewReader(string(buff))

					// перекодируем цитату из cp1251 в utf-8
					tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
					buff, err = ioutil.ReadAll(tr)
					if err != nil {
						qr.Err = err
						quotes <- qr
						return
					}

					qr.Quote += string(buff)
					quotes <- qr
					return
				}
			}
		}
	}
	qr.Err = errors.New("no citations")
	quotes <- qr
	return

	// body := html.NewTokenizer(resp.Body)

	// // В цикле пробегаем по всем токенам html
	// for {
	// 	// берём следующий токен и получаем его тип
	// 	tt := body.Next()
	// 	// Проверяем тип токена
	// 	switch tt {
	// 	case html.ErrorToken:
	// 		qr.Err = errors.New("no citations")
	// 		quotes <- qr
	// 		return
	// 	// Если это тэг
	// 	case html.StartTagToken:
	// 		tn, attr := body.TagName()
	// 		// тэг div с атрибутом
	// 		if string(tn) == "div" && attr {
	// 			// атрибут class со значением text
	// 			key, val, _ := body.TagAttr()
	// 			if string(key) == "class" && string(val) == "text" {
	// 				// то извлекаем цитату
	// 				for {
	// 					tt = body.Next()
	// 					sr := strings.NewReader(string(body.Text()))

	// 					// перекодируем цитату из cp1251 в utf-8
	// 					tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	// 					buf, err := ioutil.ReadAll(tr)

	// 					if err != nil {
	// 						qr.Err = err
	// 						quotes <- qr
	// 						return
	// 					}

	// 					qr.Quote += string(buf)

	// 					tt = body.Next()

	// 					tn, _ := body.TagName()
	// 					if string(tn) == "br" {
	// 						qr.Quote += "\n"
	// 					} else {
	// 						quotes <- qr
	// 						return
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func readBashorg(start, count int) (string, error) {
	var out_str string
	quotes := make(chan quote_res, count)

	for i := start; i < count+start; i++ {
		url := "http://bash.im/quote/" + strconv.Itoa(i)
		go readQuote(i, url, quotes)
	}

	quotes_arr := make([]string, count)
	for i := 0; i < count; i++ {
		qr := <-quotes
		if qr.Err != nil {
			return "", qr.Err
		}
		str := strconv.Itoa(qr.Num) + "\n" + qr.Quote + "\n\n"
		quotes_arr[qr.Num-start] = str
	}

	for _, val := range quotes_arr {
		out_str += val
	}

	return out_str, nil
}
