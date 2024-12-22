package main

import (
	"github.com/itpourya/Haze/app"
)

func main() {
	tgApp := app.GenerateAPP("8063671366:AAG10SKs-3sWVSMocRimuBuepZjPQ94fZFU")

	go tgApp.Start()

	ch := make(chan struct{})
	<-ch
}
