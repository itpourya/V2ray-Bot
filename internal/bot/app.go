package bot

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"

	"gopkg.in/telebot.v3"
)

func GenerateAPP(token string) *telebot.Bot {
	transport := &http.Transport{}
	proxying := "localhost:12334"
	dialer, err := proxy.SOCKS5("tcp", proxying, nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	cl := &http.Client{Transport: transport}

	settings := telebot.Settings{
		Client: cl,
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 1 * time.Second},
	}

	app, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	app.Handle("/start", start)
	app.Handle(telebot.OnText, text)
	app.Handle(telebot.OnCallback, inline)
	app.Handle(telebot.OnPhoto, recivePhoto)

	return app
}
