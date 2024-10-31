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
		log.Fatal("Can't load robot token environment file")
	}

	tgApp := app.GenerateAPP(os.Getenv("TOKEN"))

	go tgApp.Start()

	ch := make(chan struct{})
	<-ch
}
