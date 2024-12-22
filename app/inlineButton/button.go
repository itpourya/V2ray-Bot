package inlinebutton

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/itpourya/Haze/app/entity"
	"github.com/itpourya/Haze/app/serializer"
	"gopkg.in/telebot.v3"
)

func StartUpPannel() (string, *telebot.ReplyMarkup) {
	text := `سلام به ربات RedZone خوش اومدی ❤️🫠

	اینجا ۲۴ ساعت سر کاریم 🫡

	اینجا بدون محدودیت به دنیا وصل میشیم 🛜

	ارتباط با ادمین @heredeveloper

	✅ /start`
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
			{
				{
					Text: "🛍️ | گرفتن نمایندگی فروش",
					Data: "sell",
				},
			},
		},
		OneTimeKeyboard: true,
	}

	return text, button
}

func DateLimitList() (string, *telebot.ReplyMarkup) {
	text := "لطفا مدت زمان اشتراکی که میخوای رو انتخاب کن 🤔"

	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "۱ ماهه",
					Data: "1 month",
				},
			},
			{
				{
					Text: "۲ ماهه",
					Data: "2 month",
				},
			},
			{
				{
					Text: "۳ ماهه",
					Data: "3 month",
				},
			},
			{
				{
					Text: "۴ ماهه",
					Data: "4 month",
				},
			},
			{
				{
					Text: "۵ ماهه",
					Data: "5 month",
				},
			},
			{
				{
					Text: "۶ ماهه",
					Data: "6 month",
				},
			},
		},
	}

	return text, button
}

func Checkout(price int64, datalimit string, datelimit string) (string, *telebot.ReplyMarkup) {
	text := `لیست سفارش شما 🛍️` + "\n" + "قیمت نهایی : " + fmt.Sprint(price) + "حجم : " + datalimit + "مدت زمان اشتراک : " + datelimit + "ماهه"
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "کارت به کارت 🏦",
					Data: "send",
				},
			},
			{
				{
					Text: "👝 پرداخت با کیف پول",
					Data: "paywallet",
				},
			},
			{
				{
					Text: "مدیر فروش هستم ✋🏻",
					Data: "CheckManager",
				},
			},
		},
	}

	return text, button
}

func ManagerPannel() (string, *telebot.ReplyMarkup) {
	text := `پنل مدیر فروش 🧬`
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "خرید کانفیگ جدید | 🛍️",
					Data: "ManagerBuy",
				},
			},
			{
				{
					Text: "مشاهده کاربران | 🧬",
					Data: "me",
				},
			},
			{
				{
					Text: "مشاهده صورتحساب | 📝",
					Data: "invoice",
				},
			},
		},
	}

	return text, button
}

func AdminPannel() (string, *telebot.ReplyMarkup) {
	text := `به پنل ادمین خوش اومدین ✅`
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "مشاهده مدیران فروش | 🛍️",
					Data: "ShowManagerList",
				},
			},
		},
	}

	return text, button
}

func InvoicePannel(dept int64) (string, *telebot.ReplyMarkup) {
	text := `صورت حساب شما بصورت زیر است 🗒️

	💢نکته : برای تصویه حساب میتونین با ادمین در ارتباط باشین
	@heredeveloper`

	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "بدهکاری به ربات | 🤖",
				},
			},
			{
				{
					Text: fmt.Sprint(dept) + " " + "تومان",
				},
			},
		},
	}

	return text, button
}

func DataLimitList() (string, *telebot.ReplyMarkup) {
	text := "لطفا حجمی که میخوای رو انتخاب کن 🛜"

	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "10 GB",
					Data: "10",
				},
			},
			{
				{
					Text: "15 GB",
					Data: "15",
				},
			},
			{
				{
					Text: "20 GB",
					Data: "20",
				},
			},
			{
				{
					Text: "50 GB",
					Data: "50",
				},
			},
			{
				{
					Text: "60 GB",
					Data: "60",
				},
			},
			{
				{
					Text: "70 GB",
					Data: "70",
				},
			},
			{
				{
					Text: "80 GB",
					Data: "80",
				},
			},
			{
				{
					Text: "90 GB",
					Data: "90",
				},
			},
			{
				{
					Text: "100 GB",
					Data: "100",
				},
			},
		},
	}

	return text, button
}

func ManagerAnswer() string {
	text := "درخواست شما به ادمین ارسال شد، لطفا تا جواب مدیر منتظر بمانید.❤️"

	return text
}

func Locations() (string, *telebot.ReplyMarkup) {
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

func WalletPannel(wallet entity.Wallet) (string, *telebot.ReplyMarkup) {
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

func Settlement() string {
	text := `
		🤳 عزیزم لطفا یه تصویر از فیش واریزی برام ارسال کن :

	🔰 6219861929816543 - پوریا صمیمی

	✅ بعد از اینکه پرداختت تایید شد ( لینک سرور ) به صورت خودکار از طریق همین ربات برات ارسال میشه!
		`

	return text
}

func ConfigPannel(userData serializer.Response) (string, *telebot.ReplyMarkup) {
	status := userData.Status
	configName := userData.Username
	link := "https://marz.redzedshop.ir:8000" + userData.SubscriptionURL
	buyTime := strings.Replace(userData.CreatedAt[0:10], "-", "/", 2)
	expire := validateTime(userData.Expire)
	dataLimit := validateDataLimit(userData.DataLimit - userData.UsedTraffic)

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

func ConfigList(users []entity.User) (string, *telebot.ReplyMarkup) {
	text := "یکی از اکانت هاتو انتخاب کن 🤔"

	userConfigsList := [][]telebot.InlineButton{}
	for _, user := range users {
		configsList := []telebot.InlineButton{
			{
				Text: "gt-" + user.UsernameSub,
				Data: "gt-" + user.UsernameSub,
			},
		}

		userConfigsList = append(userConfigsList, configsList)
	}

	button := &telebot.ReplyMarkup{
		InlineKeyboard:  userConfigsList,
		OneTimeKeyboard: true,
	}

	return text, button
}

func ManagerList(managers []entity.Manager) (string, *telebot.ReplyMarkup) {
	text := "یکی از اکانت هاتو انتخاب کن 🤔"

	managerList := [][]telebot.InlineButton{}
	for _, manager := range managers {
		configsList := []telebot.InlineButton{
			{
				Text: manager.UserID,
				Data: "mg-" + manager.UserID,
			},
		}

		managerList = append(managerList, configsList)
	}

	button := &telebot.ReplyMarkup{
		InlineKeyboard:  managerList,
		OneTimeKeyboard: true,
	}

	return text, button
}

func AdminManagerPannel(manager entity.Manager) (string, *telebot.ReplyMarkup) {
	text := "روش پرداخت رو انتخاب کن 🤔"
	button := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "Manager UserID",
				},
				{
					Text: manager.UserID,
				},
			},
			{
				{
					Text: "Dept :",
				},
			},
			{
				{
					Text: fmt.Sprint(manager.Dept) + " تومان",
				},
			},
			{
				{
					Text: "صاف کردن بدهی 💢",
					Data: "D-" + manager.UserID,
				},
			},
		},
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
			{
				{
					Text: "مدیر فروش هستم ✋🏻",
					Data: "CheckManager",
				},
			},
		},
		OneTimeKeyboard: true,
	}

	return text, button
}

func ChargeWalletPannel() string {
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

func validateDept(input int64) string {
	var result []string
	dept := strings.Split(strconv.Itoa(int(input)), "")
	counter := 0
	for i := 0; i <= len(dept)-1; i++ {
		if counter == 2 {
			result = append(result, ",")
			counter = 0
		}

		result = append(result, dept[i])
		counter++
	}

	return strings.Join(result, "")
}
