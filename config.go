package main

import (
	"encoding/json"
	"os"
)

type configuration struct {
	Cit int //Последняя прочитанная цитата
	Cnt int //Кол-во запрашиваемых цитат
}

var config = new(configuration)

// Сохранение конфигурации
func saveConfig() {
	file, err := os.Create(os.Getenv("HOME") + "/.bashor_reader")
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

// инициализация
// читаем конфигурацию из файле
func readConfig() (conf *configuration, err error) {
	conf = new(configuration)

	if _, err = os.Stat(os.Getenv("HOME") + "/.bashor_reader"); err == nil {
		var file *os.File
		file, err = os.Open(os.Getenv("HOME") + "/.bashor_reader")
		if err != nil {
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)

		err = decoder.Decode(conf)
		if err != nil {
			return
		}
	} else {
		conf.Cit = 1
		conf.Cnt = 7
	}

	return
}
