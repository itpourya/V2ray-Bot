package main

import (
	"time"

	"gopkg.in/telebot.v3"
)

var (
	TOKEN  = "7517309635:AAEQBzrbNsC3T-scPjqHix-jwUq8P04DXwk"
	POLLER = &telebot.LongPoller{Timeout: 1 * time.Second}
)
