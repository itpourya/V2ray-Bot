package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/proxy"

	"gopkg.in/telebot.v3"
)

func GenarateAPP() *telebot.Bot {
	transport := &http.Transport{}
	proxyurl := "localhost:12334"
	dialer, err := proxy.SOCKS5("tcp", proxyurl, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	cl := &http.Client{Transport: transport}

	setings := telebot.Settings{
		Client: cl,
		Token:  TOKEN,
		Poller: POLLER,
	}

	app, err := telebot.NewBot(setings)
	if err != nil {
		log.Fatal(err)
	}

	app.Handle("/start", start)
	app.Handle(telebot.OnText, text)
	app.Handle(telebot.OnCallback, inline)
	app.Handle(telebot.OnPhoto, recivePhoto)

	return app
}
