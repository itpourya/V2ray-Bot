package app

import (
	"time"

	"gopkg.in/telebot.v3"
)

var (
	TOKEN  = "6466071910:AAEkVNnmsVrnfSmDfpmeEcSR-tT4phRvgNY"
	POLLER = &telebot.LongPoller{Timeout: 1 * time.Second}
)
