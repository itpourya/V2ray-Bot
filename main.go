package main

import (
	"os"

	"github.com/itpourya/Haze/app"
)

func main() {
	tgApp := app.GenerateAPP(os.Getenv("TOKEN"))

	go tgApp.Start()

	ch := make(chan struct{})
	<-ch
}
