package main

import (
	"github.com/itpourya/Haze/app"
)

func main() {
	tgApp := app.GenerateAPP("7394922553:AAHow5sFxgLnzIaJHXNXPHBpiVwYq_Cr8ao")

	go tgApp.Start()

	ch := make(chan struct{})
	<-ch
}
