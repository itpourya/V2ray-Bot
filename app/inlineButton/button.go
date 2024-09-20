package inlinebutton

import (
	"fmt"
	"strings"
	"time"

	"github.com/itpourya/Haze/entity"
	"github.com/itpourya/Haze/serializer"
	"gopkg.in/telebot.v3"
)

func Start() (string, *telebot.ReplyMarkup) {
	text := `سلااام به ربات HAZE  خوش اومدی 🫡❤️

		اینجا ۲۴ ساعت در خدمت شما هستیم 🔥

		ما اینجاییم تا شما را بدون هیچ محدویتی به شبکه جهانی متصل کنیم ❤️

		✅ کیفیت در ساخت انواع کانکشن ها
		📡 برقرای امنیت در ارتباط شما
		☎️ پشتیبانی تا روز آخر


		ارتباط با ادمین @heredeveloper
		🚪 /start`
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "🛒 | خرید کانفیگ جدید",
					Data: "buy",
				},
			},
			{
				{
					Text: "📦 | کانفیگ های من",
					Data: "me",
				},
			},
			{
				{
					Text: "💸 | ارسال رسید شارژ",
					Data: "charge",
				},
			},
			{
				{
					Text: "💳 | مشاهده کیف پول",
					Data: "wallet",
				},
			},
		},
		OneTimeKeyboard: true,
	}

	return text, button
}

func Buy() (string, *telebot.ReplyMarkup) {
	text := `سلااام به ربالوکیشن مدنظر خودتون رو انتخاب کنین 🤔`
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "سرور آلمان | 🇩🇪",
					Data: "germany",
				},
			},
		},
	}

	return text, button
}

func Germany() (string, *telebot.ReplyMarkup) {
	text := `سلااام به ربالوکیشن مدنظر خودتون رو انتخاب کنین 🤔`
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "۱۰ گیگ |  28,000 تومان",
					Data: "10",
				},
			},
			{
				{
					Text: "۱۵ گیگ |  38,000 تومان",
					Data: "15",
				},
			},
			{
				{
					Text: "۲۰ گیگ |  50,000 تومان",
					Data: "20",
				},
			},
			{
				{
					Text: "۵۰ گیگ |  120,000 تومان",
					Data: "50",
				},
			},
			{
				{
					Text: "۱۰۰ گیگ |  180,000 تومان",
					Data: "100",
				},
			},
		},
	}

	return text, button
}

func ShowWallet(wallet entity.Wallet) (string, *telebot.ReplyMarkup) {
	text := "میتونین میزان شارژ کیف پول خودتون رو ببینین 😊"
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: fmt.Sprint(wallet.Balance) + " " + "تومان",
				},
			},
		},
	}

	return text, button
}

func Send() string {
	text := `
		🤳 عزیزم لطفا یه تصویر از فیش واریزی برام ارسال کن :

	🔰 6219861929816543 - پوریا صمیمی

	✅ بعد از اینکه پرداختت تایید شد ( لینک سرور ) به صورت خودکار از طریق همین ربات برات ارسال میشه!
		`

	return text
}

func Me(userData serializer.Response) (string, *telebot.ReplyMarkup) {
	status := userData.Status
	configName := userData.Username
	link := "https://marz.ikernel.ir:8000/" + userData.SubscriptionURL
	buyTime := strings.Replace(userData.CreatedAt[0:10], "-", "/", 2)
	expire := validateTime(userData.Expire)
	dataLimit := validateDataLimit(userData.DataLimit)

	text := `وضعیت کانفیگ: ` + status + `🟢

		🦠 نام کانفیگ :` + configName + `

		🔗 لینک اتصال:
		` + link
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				telebot.InlineButton{
					Text: buyTime,
				},
				{
					Text: "⏰ تاریخ خرید",
				},
			},
			{
				telebot.InlineButton{
					Text: expire,
				},
				{
					Text: "⏰ تاریخ انقضاء",
				},
			},
			{
				telebot.InlineButton{
					Text: dataLimit,
				},
				{
					Text: "⏳ حجم باقیمانده",
				},
			},
			{
				{
					Text: "♻️ تمدید یکماهه سرویس",
					Data: "remonth" + configName,
				},
			},
			{
				{
					Text: "🚀 افزایش حجم سرویس",
					Data: "retraffic" + configName,
				},
			},
		},
		OneTimeKeyboard: true,
	}

	return text, button
}

func ShowConfigsMe(users []entity.User) (string, *telebot.ReplyMarkup) {
	text := "یکی از اکانت هاتو انتخاب کن🤔"

	userdetail := [][]telebot.InlineButton{}
	for _, user := range users {
		ls := []telebot.InlineButton{
			{
				Text: "gt-" + user.UsernameSub,
				Data: "gt-" + user.UsernameSub,
			},
		}

		userdetail = append(userdetail, ls)
	}

	button := &telebot.ReplyMarkup{
		InlineKeyboard:  userdetail,
		OneTimeKeyboard: true,
	}

	return text, button
}

func Remonth() (string, *telebot.ReplyMarkup) {
	text := "روش پرداخت رو انتخاب کن 🤔"
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "تمدید یک ماهه | 32,000 تومان",
					Data: "remonthpay",
				},
			},
		},
		OneTimeKeyboard: true,
	}

	return text, button
}

func Wallet() string {
	text := `لطفا مقداری که میخوای کیف پولتو شارژ کنی رو به عدد وارد کن

		(توجه کن عددی که وارد میکنی بین ۵ هزار تومان تا  حداکثر ۱۰ میلیون تومان باشه 😇)`

	return text
}

func validateTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	formatted := t.Format("2006/01/02")

	return formatted
}

func validateDataLimit(dataLimit int64) string {
	gb := float64(dataLimit) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2f GB", gb)
}
