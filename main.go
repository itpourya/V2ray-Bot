package main

import (
	"github.com/itpourya/Haze/app"
)

func main() {
	tgApp := app.GenerateAPP("6466071910:AAEkVNnmsVrnfSmDfpmeEcSR-tT4phRvgNY")

	go tgApp.Start()

	ch := make(chan struct{})
	<-ch
}
