package main

import (
	"github.com/itpourya/Haze/config"
	"github.com/itpourya/Haze/internal/bot"
)

func main() {
	tgApp := bot.GenerateAPP(config.TOKEN)

	go tgApp.Start()

	ch := make(chan struct{})
	<-ch
}
