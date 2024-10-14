package main

import (
	"log"
	"os"

	"github.com/itpourya/Haze/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't load robot token envirement file")
	}

	app := app.GenarateAPP(os.Getenv("TOKEN"))

	go app.Start()

	ch := make(chan struct{})
	<-ch
}
