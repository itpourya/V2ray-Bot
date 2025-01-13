package bot

import (
	"context"
	"fmt"
	cache2 "github.com/itpourya/Haze/internal/cache"
	"github.com/itpourya/Haze/internal/database"
	"github.com/itpourya/Haze/internal/inlineButton"
	"github.com/itpourya/Haze/internal/marzban"
	"github.com/itpourya/Haze/internal/repository"
	"github.com/itpourya/Haze/internal/service"
	"github.com/itpourya/Haze/internal/validator"
	"github.com/itpourya/Haze/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"gopkg.in/telebot.v3"
)

var (
	mrz            = marzban.NewMarzbanClient()
	rdb            = cache2.NewCache()
	ctxb           = context.Background()
	db             = database.New()
	userRepository = repository.NewRepository(db)
	userService    = service.NewUserService(userRepository)
)

func start(ctx telebot.Context) error {
	sender := strconv.Itoa(int(ctx.Sender().ID))

	if utils.IsManager(sender, userService) {
		text, button := inlinebutton.ManagerPannel()

		return ctx.Send(text, button)
	}

	if utils.IsAdmin(sender) {
		text, button := inlinebutton.AdminPannel()

		return ctx.Send(text, button)
	}

	text, button := inlinebutton.StartUpPannel()
	err := userService.CreateUserWalletService(fmt.Sprint(ctx.Sender().ID))
	if err != nil {
		log.Error("Handler Error", err.Error())
	}

	return ctx.Send(text, button)
}

func text(ctx telebot.Context) error {
	data := rdb.GetDel(ctxb, fmt.Sprint(ctx.Sender().ID))
	validate := validator.Validate(data.String(), fmt.Sprint(ctx.Sender().ID))
	pattern := `^\d+\|[a-zA-Z]+$`

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Error("Error compiling regex", err)
	}

	if validate.ReChargeWallet {
		validate.ChargeWallet = ctx.Text()
		err := rdb.Set(ctxb, fmt.Sprint(ctx.Sender().ID), validate, time.Duration(1*time.Hour))
		if err != nil {
			log.Info(err)
		}
		return ctx.Send(inlinebutton.Settlement())
	}

	if re.MatchString(ctx.Text()) {
		data := strings.Split(ctx.Text(), "|")
		TargetUserID := data[0]
		TargetUserOwnerName := data[1]

		userService.EnterConfigOwnerNameService(TargetUserOwnerName, TargetUserID)
		return ctx.Send("Changed.")
	}

	return ctx.Send("خیلی ببخشید ولی نفهمیدم چی میخواین 🤔")
}

func inline(ctx telebot.Context) error {
	command := ctx.Data()
	userData := ctx.Sender()
	userIDstr := strconv.Itoa(int(userData.ID))
	AdminUserID := int64(6556338275)

	if command == "EnterOwnerName" {
		ctx.Send("لطفا اسم صاحب کانفیگ را بهمراه ایدی کاربر بصورت زیر وارد کنید\nId+username")
		return nil
	}

	if command == "ShowManagerList" {
		managersList := userService.GetManagerListService()
		text, button := inlinebutton.ManagerList(managersList)

		return ctx.Edit(text, button)
	}

	if strings.HasPrefix(command, "mg-") {
		managerUserID := strings.Replace(command, "mg-", "", 1)
		manager := userService.GetManagerService(managerUserID)
		text, button := inlinebutton.AdminManagerPannel(manager)

		return ctx.Edit(text, button)
	}

	if strings.HasPrefix(command, "D-") {
		managerUserID := strings.Replace(command, "D-", "", 1)
		userService.ClearManagerDeptService(managerUserID)

		return ctx.Edit("بدهی منیجر کاملا تسویه شد ✅")
	}

	if strings.HasSuffix(command, " month") {
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		validate.DateLimit = strings.Replace(command, " month", "", 1)
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
		if err != nil {
			log.Info(err)
		}

		text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)

		return ctx.Edit(text, button)
	}

	if command == "invoice" {
		dept := userService.GetInvoiceService(userIDstr)
		text, button := inlinebutton.InvoicePannel(dept)

		return ctx.Edit(text, button)
	}

	if command == "CheckManager" {
		manager := utils.IsManager(userIDstr, userService)
		if manager {
			data := rdb.GetDel(ctxb, userIDstr).String()
			validate := validator.Validate(data, fmt.Sprint(userIDstr))
			data_limit, _ := strconv.Atoi(strings.Replace(validate.DataLimitCharge, "GB", "", 1))

			if validate.Buy {
				username := userService.GenerateUsernameService(userIDstr)
				response, _ := mrz.CreateMarzbanUser(username, int(data_limit), validate.DateLimit)
				userService.IncreaseManagerDeptService(userIDstr, validate.Price)
				return ctx.Send(`
												😍 سفارش جدید شما
									📡 پروتکل: vless
									🔋حجم سرویس: ` + validate.DataLimitCharge + `
									⏰ مدت سرویس: 30 روز

									subscription : ` + `https://marz.ikernel.ir:8000` + response.SubscriptionURL + `
												`)
			}

			if validate.Recharge {
				err := mrz.DataLimitUpdate(validate.ConfigName, validate.DataLimitCharge)
				if err != nil {
					log.Info(err)
				}

				userService.IncreaseManagerDeptService(userIDstr, validate.Price)

				ctx.Send("پرداخت شما با موفقیت انجام شد. حجم مورد نظر به سرویس شما اضافه شد")
				return ctx.Delete()
			}

			if validate.Remonth {
				err := mrz.ExpireUpdate(validate.ConfigName)
				if err != nil {
					log.Info(err)
				}

				userService.IncreaseManagerDeptService(userIDstr, validate.Price)

				ctx.Send("پرداخت شما با موفقیت انجام شد. یک ماه به سرویس شما اضافه شد")
				return ctx.Delete()
			}
		}
	}

	if command == "ManagerBuy" {
		var payload cache2.CachePayload
		payload.Buy = true
		payload.Manager = true
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Info(err)
		}

		text, button := inlinebutton.DataLimitList()
		return ctx.Edit(text, button)
	}

	if command == "sell" {
		text := inlinebutton.ManagerAnswer()

		ctx.Bot().Send(&telebot.Chat{ID: AdminUserID}, "Manager : "+userIDstr+" | "+userData.Username, &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{
					{
						Text: "Add",
						Data: "add" + userIDstr,
					},
				},
			},
			OneTimeKeyboard: true,
		})

		return ctx.Send(text)
	}

	if strings.HasPrefix(command, "add") {
		managerID := strings.Replace(command, "add", "", 1)
		managerIDinteger, _ := strconv.Atoi(managerID)
		userService.CreateManagerService(managerID)
		ctx.Bot().Send(&telebot.Chat{ID: int64(managerIDinteger)}, "تبریک میگم شما به عنوان مدیر فروش در ربات ما ثبت شدین، میتونین با اجرای دوباره ی ربات به پنل مدیریت مدیر های فروش دسترسی داشته باشین ❤️")
		return ctx.Delete()
	}

	if strings.HasPrefix(command, "user") {
		userID := strings.Replace(command, "user", "", 1)
		recipientID, _ := strconv.Atoi(userID)
		data := rdb.GetDel(ctxb, userID).String()
		validate := validator.Validate(data, fmt.Sprint(userID))
		data_limit, _ := strconv.Atoi(strings.Replace(validate.DataLimitCharge, "GB", "", 1))

		if validate.Buy {
			username := userService.GenerateUsernameService(userID)
			response, _ := mrz.CreateMarzbanUser(username, int(data_limit), validate.DateLimit)
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
				log.Info(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(recipientID)}, "پرداخت شما با موفقیت انجام شد. حجم مورد نظر به سرویس شما اضافه شد")
			return ctx.Delete()
		}

		if validate.Remonth {
			err := mrz.ExpireUpdate(validate.ConfigName)
			if err != nil {
				log.Info(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(recipientID)}, "پرداخت شما با موفقیت انجام شد. یک ماه به سرویس شما اضافه شد")
			return ctx.Delete()
		}

		if validate.ReChargeWallet && validate.ChargeWallet != "" {
			charge, _ := strconv.Atoi(validate.ChargeWallet)
			err := userService.IncreaseUserBalanceService(userID, charge)
			if err != nil {
				log.Info(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(recipientID)}, "پرداخت شما با موفقیت انجام شد. مبلغ مورد نظر به کیف شما اضافه شد🚀")
			return ctx.Delete()
		}

		return ctx.Delete()
	}

	if command == "wallet" {
		userWallet := userService.GetUserWalletService(userIDstr)
		text, button := inlinebutton.WalletPannel(userWallet)

		return ctx.Edit(text, button)
	}

	if command == "charge" {
		var payload cache2.CachePayload
		payload.Buy = false
		payload.ReChargeWallet = true
		payload.Recharge = false
		payload.Remonth = false
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Info(err)
		}

		return ctx.Edit(inlinebutton.ChargeWalletPannel())
	}

	if command == "paywallet" {
		data := rdb.GetDel(ctxb, userIDstr)
		validate := validator.Validate(data.String(), userIDstr)
		err := userService.DicreaseUserBalanceService(userIDstr, int(validate.Price))
		if err != nil {
			return ctx.Edit("شرمندت کیف پول شما موجودی کافی رو برای پرداخت نداره 🫠")
		}

		data_limit, _ := strconv.Atoi(strings.Replace(validate.DataLimitCharge, "GB", "", 1))

		if validate.Buy {
			username := userService.GenerateUsernameService(userIDstr)
			response, _ := mrz.CreateMarzbanUser(username, int(data_limit), validate.DateLimit)
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
				log.Info(err)
			}

			_, _ = ctx.Bot().Send(&telebot.Chat{ID: int64(userData.ID)}, "پرداخت شما با موفقیت انجام شد. حجم مورد نظر به سرویس شما اضافه شد")
			return ctx.Delete()
		}

		if validate.Remonth {
			err := mrz.ExpireUpdate(validate.ConfigName)
			if err != nil {
				log.Info(err)
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
		var payload cache2.CachePayload
		manager := utils.IsManager(userIDstr, userService)
		if manager {
			payload.Manager = true
		}
		payload.ConfigName = strings.Replace(command, "retraffic", "", 1)
		payload.Remonth = false
		payload.Recharge = true
		payload.Buy = false
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Info(err)
		}

		text, button := inlinebutton.DataLimitList()

		return ctx.Edit(text, button)
	}

	if strings.HasPrefix(command, "gt-") {
		username := strings.Replace(command, "gt-", "", 1)
		user, _, err := mrz.GetMarzbanUser(username)
		if err != nil {
			log.Info(err)
		}

		text, button := inlinebutton.ConfigPannel(user)

		return ctx.Edit(text, button)
	}

	if strings.HasPrefix(command, "remonthpay") {
		return ctx.Send(inlinebutton.Settlement())
	}

	if strings.HasPrefix(command, "remonth") {
		var payload cache2.CachePayload
		manager := utils.IsManager(userIDstr, userService)
		if manager {
			payload.Manager = true
			payload.Price = 15000
		} else {
			payload.Price = 32000
		}
		payload.Buy = false
		payload.Recharge = false
		payload.Remonth = true
		payload.ConfigName = strings.Replace(command, "remonth", "", 1)
		err := rdb.Set(ctxb, userIDstr, payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Info(err)
		}
		text, button := inlinebutton.Remonth()

		ctx.Edit(text, button)
	}

	if command == "me" {
		user := userService.GetUserByUserIDService(userIDstr)

		text, button := inlinebutton.ConfigList(user)
		return ctx.Edit(text, button)
	}

	if command == "buy" {
		var payload cache2.CachePayload
		payload.Buy = true
		err := rdb.Set(ctxb, fmt.Sprint(userData.ID), payload, time.Duration(1*time.Hour))
		if err != nil {
			log.Info(err)
		}
		text, button := inlinebutton.Locations()
		return ctx.Edit(text, button)
	}

	if command == "germany" {
		text, button := inlinebutton.DataLimitList()
		return ctx.Edit(text, button)
	}

	if command == "send" {
		return ctx.Send(inlinebutton.Settlement())
	}

	switch command {
	case "10":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "10GB"
			if validate.Manager {
				validate.Price = 10000
			} else {
				validate.Price = 28000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "10GB"
			if validate.Manager {
				validate.Price = 10
			} else {
				validate.Price = 28000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}

		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "15":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "15GB"
			if validate.Manager {
				validate.Price = 15000
			} else {
				validate.Price = 38000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "15GB"
			if validate.Manager {
				validate.Price = 15000
			} else {
				validate.Price = 38000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "20":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "20GB"
			if validate.Manager {
				validate.Price = 20000
			} else {
				validate.Price = 50000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "20GB"
			if validate.Manager {
				validate.Price = 20000
			} else {
				validate.Price = 50000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "50":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "50GB"
			if validate.Manager {
				validate.Price = 50000
			} else {
				validate.Price = 120000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "50GB"
			if validate.Manager {
				validate.Price = 50000
			} else {
				validate.Price = 120000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "100":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "100GB"
			if validate.Manager {
				validate.Price = 100000
			} else {
				validate.Price = 180000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "100GB"
			if validate.Manager {
				validate.Price = 100000
			} else {
				validate.Price = 180000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "60":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "60GB"
			if validate.Manager {
				validate.Price = 60000
			} else {
				validate.Price = 150000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "60GB"
			if validate.Manager {
				validate.Price = 60000
			} else {
				validate.Price = 150000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "70":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "70GB"
			if validate.Manager {
				validate.Price = 70000
			} else {
				validate.Price = 175000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "70GB"
			if validate.Manager {
				validate.Price = 70000
			} else {
				validate.Price = 175000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "80":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "80GB"
			if validate.Manager {
				validate.Price = 80000
			} else {
				validate.Price = 200000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "80GB"
			if validate.Manager {
				validate.Price = 60000
			} else {
				validate.Price = 200000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)

	case "90":
		dataCache := rdb.GetDel(ctxb, fmt.Sprint(userData.ID))
		validate := validator.Validate(dataCache.String(), userIDstr)
		if validate.Buy {
			validate.DataLimitCharge = "90GB"
			if validate.Manager {
				validate.Price = 90000
			} else {
				validate.Price = 225000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}
		}

		if validate.Recharge {
			validate.DataLimitCharge = "90GB"
			if validate.Manager {
				validate.Price = 90000
			} else {
				validate.Price = 225000
			}
			err := rdb.Set(ctxb, fmt.Sprint(userData.ID), validate, time.Duration(1*time.Hour))
			if err != nil {
				log.Info(err)
			}

			text, button := inlinebutton.Checkout(validate.Price, validate.DataLimitCharge, validate.DateLimit)
			return ctx.Edit(text, button)
		}
		text, button := inlinebutton.DateLimitList()
		return ctx.Edit(text, button)
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
