package main

import "github.com/itpourya/Haze/app"

func main() {
	app := app.GenarateAPP()

	go app.Start()

	ch := make(chan struct{})
	<-ch
}
