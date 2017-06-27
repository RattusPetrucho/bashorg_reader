package main

import (
	"errors"
	"strconv"

	// get запрос
	"net/http"

	// Перевод из cp1251 в utf-8
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type quote_res struct {
	Num   int
	Quote string
	Err   error
}

// Семафор дл  ограничивающий кол-во запросов к башоргу 10ю
var semaphore = make(chan struct{}, 10)

// Читаем цитату по номеру num
// и выбираем текст цитаты
func readQuote(num int, url string, quotes chan<- quote_res) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

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
