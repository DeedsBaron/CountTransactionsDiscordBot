package main

import (
	"fmt"
	"go-transaction-bot/bot"
	"go-transaction-bot/config"
)

var (
	configLoad = make(chan struct{})
)

func main() {
	err := config.ReadConfig()
	close(configLoad)
	fmt.Println(config.Config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bot.Start()
	<-make(chan struct{})
	return
}
