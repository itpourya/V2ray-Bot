package app

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/itpourya/Haze/app/cache"
	"github.com/itpourya/Haze/app/database"
	inlinebutton "github.com/itpourya/Haze/app/inlineButton"
	"github.com/itpourya/Haze/app/marzban"
	"github.com/itpourya/Haze/app/repository"
	"github.com/itpourya/Haze/app/service"
	"github.com/itpourya/Haze/app/validator"
	"gopkg.in/telebot.v3"
)

var (
	mrz            = marzban.NewMarzbanClient()
	rdb            = cache.NewCache()
	ctxb           = context.Background()
	db             = database.New()
	userRepository = repository.NewRepository(db)
	userService    = service.NewUserService(userRepository)
)

func start(ctx telebot.Context) error {
	text, button := inlinebutton.Start()
	err := userService.CreateUserWallet(fmt.Sprint(ctx.Sender().ID))
	if err != nil {
		log.Println(err)
	}

	return ctx.Send(text, button)
}

func text(ctx telebot.Context) error {
	data := rdb.GetDel(ctxb, fmt.Sprint(ctx.Sender().ID))
	validate := validator.Validate(data.String(), fmt.Sprint(ctx.Sender().ID))

	if validate.ReChargeWallet {
		validate.ChargeWallet = ctx.Text()
		err := rdb.Set(ctxb, fmt.Sprint(ctx.Sender().ID), validate, time.Duration(1*time.Hour))
		if err != nil {
			log.Println(err)
		}
		return ctx.Send(inlinebutton.Send())
	}
	return ctx.Send("invalid")
}

func inline(ctx telebot.Context) error {
	command := ctx.Data()
	userData := ctx.Sender()
	userIDstr := strconv.Itoa(int(userData.ID))

	if strings.HasPrefix(command, "user") {
		userID := strings.Replace(command, "user", "", 1)
		recipientID, _ := strconv.Atoi(userID)
		data := rdb.GetDel(ctxb, userID).String()
		validate := validator.Validate(data, fmt.Sprint(userID))
		data_limit, _ := strconv.Atoi(strings.Replace(validate.DataLimitCharge, "GB", "", 1))

		if validate.Buy {
			username := userService.CreateUsername(userID)
			response, _ := mrz.CreateUserAccount(username, int(data_limit))
			ctx.Bot().Send(&telebot.Chat{ID: int64(recipientID)}, `
											😍 سفارش جدید شما
								📡 پروتکل: vless
								🔋حجم سرویس: `+validate.DataLimitCharge+`
								⏰ مدت سرویس: 30 روز

								subscription : `+`https://marz.ikernel.ir:8000`+response.SubscriptionURL+`
											`)
		}

		if validate.Recharge {
			err := mrz.DataLimitUpdate(validate.ConfigName, validate.DataLimitCharge)
			if err != nil {
				log.Println(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(recipientID)}, "پرداخت شما با موفقیت انجام شد. حجم مورد نظر به سرویس شما اضافه شد")
			return ctx.Delete()
		}

		if validate.Remonth {
			err := mrz.ExpireUpdate(validate.ConfigName)
			if err != nil {
				log.Println(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(recipientID)}, "پرداخت شما با موفقیت انجام شد. یک ماه به سرویس شما اضافه شد")
			return ctx.Delete()
		}

		if validate.ReChargeWallet && validate.ChargeWallet != "" {
			charge, _ := strconv.Atoi(validate.ChargeWallet)
			fmt.Println(charge)
			err := userService.IncreaseUserBalance(userID, charge)
			if err != nil {
				log.Println(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(recipientID)}, "پرداخت شما با موفقیت انجام شد. مبلغ مورد نظر به کیف شما اضافه شد🚀")
			return ctx.Delete()
		}

		return ctx.Delete()
	}

	if command == "wallet" {
		userWallet := userService.GetUserWallet(userIDstr)
		text, button := inlinebutton.ShowWallet(userWallet)

		return ctx.Edit(text, button)
	}

	if command == "charge" {
		var payload cache.CachePayload
		payload.Buy = false
		payload.ReChargeWallet = true
		payload.Recharge = false
		payload.Remonth = false
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Println(err)
		}

		return ctx.Edit(inlinebutton.Wallet())
	}

	if command == "paywallet" {
		data := rdb.GetDel(ctxb, userIDstr)
		validate := validator.Validate(data.String(), userIDstr)
		err := userService.DicreaseUserBalance(userIDstr, int(validate.Price))
		if err != nil {
			return ctx.Edit("شرمندت کیف پول شما موجودی کافی رو برای پرداخت نداره 🫠")
		}

		data_limit, _ := strconv.Atoi(strings.Replace(validate.DataLimitCharge, "GB", "", 1))

		if validate.Buy {
			username := userService.CreateUsername(userIDstr)
			response, _ := mrz.CreateUserAccount(username, int(data_limit))
			ctx.Bot().Send(&telebot.Chat{ID: int64(userData.ID)}, `
											😍 سفارش جدید شما
								📡 پروتکل: vless
								🔋حجم سرویس: `+validate.DataLimitCharge+`
								⏰ مدت سرویس: 30 روز

								subscription : `+`https://marz.ikernel.ir:8000/`+response.SubscriptionURL+`
											`)
		}

		if validate.Recharge {
			err := mrz.DataLimitUpdate(validate.ConfigName, validate.DataLimitCharge)
			if err != nil {
				log.Println(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(userData.ID)}, "پرداخت شما با موفقیت انجام شد. حجم مورد نظر به سرویس شما اضافه شد")
			return ctx.Delete()
		}

		if validate.Remonth {
			err := mrz.ExpireUpdate(validate.ConfigName)
			if err != nil {
				log.Println(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(userData.ID)}, "پرداخت شما با موفقیت انجام شد. یک ماه به سرویس شما اضافه شد")
			return ctx.Delete()
		}

	}

	if strings.HasPrefix(command, "dis") {
		userID := strings.Replace(command, "dis", "", 1)
		recive, _ := strconv.Atoi(userID)

		ctx.Bot().Send(&telebot.Chat{ID: int64(recive)}, "رسید پرداخت ارسالی از طرف شما قابل قبول نیست و تایید نمیشود ❤️")
		return ctx.Delete()
	}

	if strings.HasPrefix(command, "retraffic") {
		var payload cache.CachePayload
		payload.ConfigName = strings.Replace(command, "retraffic", "", 1)
		payload.Remonth = false
		payload.Recharge = true
		payload.Buy = false
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Println(err)
		}

		text, button := inlinebutton.Germany()

		return ctx.Edit(text, button)
	}

	if strings.HasPrefix(command, "gt-") {
		username := strings.Replace(command, "gt-", "", 1)
		user, _, err := mrz.GetUser(username)
		if err != nil {
			log.Println(err)
		}

		text, button := inlinebutton.Me(user)

		return ctx.Edit(text, button)
	}

	if strings.HasPrefix(command, "remonthpay") {
		return ctx.Send(inlinebutton.Send())
	}

	if strings.HasPrefix(command, "remonth") {
		var payload cache.CachePayload
		payload.Buy = false
		payload.Recharge = false
		payload.Remonth = true
		payload.ConfigName = strings.Replace(command, "remonth", "", 1)
		err := rdb.Set(ctxb, userIDstr, payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Println(err)
		}
		text, button := inlinebutton.Remonth()

		ctx.Edit(text, button)
	}

	if command == "me" {
		user := userService.GetUserByUserID(userIDstr)

		text, button := inlinebutton.ShowConfigsMe(user)
		return ctx.Edit(text, button)
	}

	if command == "buy" {
		var payload cache.CachePayload
		payload.Buy = true
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Println(err)
		}
		text, button := inlinebutton.Buy()
		return ctx.Edit(text, button)
	}

	if command == "germany" {
		text, button := inlinebutton.Germany()
		return ctx.Edit(text, button)
	}

	if command == "send" {
		return ctx.Send(inlinebutton.Send())
	}

	switch command {
	case "10":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "10GB"
			validate.Price = 28000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "10GB"
			validate.Price = 28000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}
		return ctx.Edit(`
		〽️ نام پلن: 10 گیگ
        ➖➖➖➖➖➖➖
        💎 قیمت پنل : 28,000 تومان
        ➖➖➖➖➖➖➖
		`, &telebot.ReplyMarkup{
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
			},
		})

	case "15":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "15GB"
			validate.Price = 38000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "15GB"
			validate.Price = 38000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}
		return ctx.Edit(`〽️ نام پلن: 15 گیگ
        ➖➖➖➖➖➖➖
        💎 قیمت پنل : 38,000 تومان
        ➖➖➖➖➖➖➖`, &telebot.ReplyMarkup{
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
			},
		})

	case "20":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "20GB"
			validate.Price = 50000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "20GB"
			validate.Price = 50000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}
		return ctx.Edit(`〽️ نام پلن: 20 گیگ
        ➖➖➖➖➖➖➖
        💎 قیمت پنل : 50,000 تومان
        ➖➖➖➖➖➖➖`, &telebot.ReplyMarkup{
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
			},
		})

	case "50":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "50GB"
			validate.Price = 120000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "50GB"
			validate.Price = 120000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}
		return ctx.Edit(`〽️ نام پلن: 50 گیگ
        ➖➖➖➖➖➖➖
        💎 قیمت پنل : 120,000 تومان
        ➖➖➖➖➖➖➖`, &telebot.ReplyMarkup{
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
			},
		})

	case "100":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "100GB"
			validate.Price = 180000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "100GB"
			validate.Price = 180000
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Println(err)
			}
		}
		return ctx.Edit(`〽️ نام پلن: 100 گیگ
        ➖➖➖➖➖➖➖
        💎 قیمت پنل : 180,000 تومان
        ➖➖➖➖➖➖➖`, &telebot.ReplyMarkup{
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
			},
		})
	}
	return nil
}

func recivePhoto(ctx telebot.Context) error {
	photo := ctx.Message().Photo
	sender := ctx.Sender().ID
	senderIDstr := strconv.Itoa(int(sender))
	recipientID := int64(6556338275)
	data := rdb.Get(ctxb, senderIDstr).String()
	message := ""

	validate := validator.Validate(data, senderIDstr)

	if strings.HasPrefix(data, "redis") {
		message = "FAKE"
	}

	if validate.Remonth {
		message = "یک ماهه"
		_, _ = ctx.Bot().Send(&telebot.Chat{ID: recipientID}, &telebot.Photo{File: photo.File}, &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{
					{
						Text: "تایید ✅",
						Data: "user" + fmt.Sprint(sender),
					},
					telebot.InlineButton{
						Text: message,
					},
				},
				{
					{
						Text: "غیرقابل قبول ❌",
						Data: "dis" + fmt.Sprint(sender),
					},
				},
			},
			OneTimeKeyboard: true,
		})

		return nil
	}

	if validate.ReChargeWallet {
		message = validate.ChargeWallet
		_, _ = ctx.Bot().Send(&telebot.Chat{ID: recipientID}, &telebot.Photo{File: photo.File}, &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{
					{
						Text: "تایید ✅",
						Data: "user" + fmt.Sprint(sender),
					},
					telebot.InlineButton{
						Text: message,
					},
				},
				{
					{
						Text: "غیرقابل قبول ❌",
						Data: "dis" + fmt.Sprint(sender),
					},
				},
			},
			OneTimeKeyboard: true,
		})

		return nil
	}

	_, _ = ctx.Bot().Send(&telebot.Chat{ID: recipientID}, &telebot.Photo{File: photo.File}, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Text: "تایید ✅",
					Data: "user" + fmt.Sprint(sender),
				},
				telebot.InlineButton{
					Text: validate.DataLimitCharge,
				},
			},
			{
				{
					Text: "غیرقابل قبول ❌",
					Data: "dis" + fmt.Sprint(sender),
				},
			},
		},
		OneTimeKeyboard: true,
	})

	return nil
}
