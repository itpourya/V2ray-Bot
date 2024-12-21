package repository

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/itpourya/Haze/app/entity"
	"gorm.io/gorm"
)

type APIrepository interface {
	RegisterUser(userID string, userSub string) error
	GetUserConfigsAccount(userID string) []entity.User
	IncreseUserBalance(userID string, charge int) error
	CreateUserWallet(userID string) error
	GetUserWallet(userID string) entity.Wallet
	DecreaseUserBalance(userID string, amount int) error
	CreateManager(userID string) error
	GetManager(userID string) entity.Manager
	IncreaseManagerDept(userID string, price int64) bool
	ClearManagerDept(userID string) bool
	GetInvoice(userID string) int64
	GetManagerList() []entity.Manager
}

type apiRepository struct {
	db *gorm.DB
}

func NewRepository(conn *gorm.DB) APIrepository {
	return &apiRepository{
		db: conn,
	}
}

func (session *apiRepository) GetManagerList() []entity.Manager {
	var manager []entity.Manager

	err := session.db.Find(&manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
	}

	return manager
}

func (session *apiRepository) GetInvoice(userID string) int64 {
	var manager entity.Manager

	err := session.db.Where("user_id = ?", userID).Find(&manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
	}

	return manager.Dept
}

func (session *apiRepository) ClearManagerDept(userID string) bool {
	var manager entity.Manager

	err := session.db.Where("user_id = ?", userID).Find(&manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
		return false
	}

	manager.Dept = 0

	err = session.db.Save(manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
		return false
	}

	return true
}

func (session *apiRepository) IncreaseManagerDept(userID string, price int64) bool {
	var manager entity.Manager

	err := session.db.Where("user_id = ?", userID).Find(&manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
		return false
	}

	manager.Dept += price

	err = session.db.Save(manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
		return false
	}

	return true
}

func (session *apiRepository) GetManager(userID string) entity.Manager {
	var manager entity.Manager

	err := session.db.Where("user_id = ?", userID).Find(&manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
	}

	return manager
}

func (session *apiRepository) CreateManager(userID string) error {
	var manager entity.Manager

	manager.UserID = userID

	err := session.db.Create(&manager)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
		return err.Error
	}

	return nil
}

func (session *apiRepository) RegisterUser(userID string, userSub string) error {
	var user entity.User
	user.UserID = userID
	user.UsernameSub = userSub

	err := session.db.Save(&user)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
		return err.Error
	}

	return nil
}

func (session *apiRepository) GetUserConfigsAccount(userID string) []entity.User {
	var user []entity.User

	err := session.db.Model(&entity.User{}).Where("user_id = ?", userID).Find(&user)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
	}

	return user
}

func (session *apiRepository) IncreseUserBalance(userID string, charge int) error {
	walletExist := session.GetUserWallet(userID)

	if walletExist.UserID != userID {
		session.CreateUserWallet(userID)
		walletExist = session.GetUserWallet(userID)
	}

	walletExist.Balance += int64(charge)
	err := session.db.Model(&entity.Wallet{}).Where("user_id = ?", userID).Update("balance", walletExist.Balance)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
		return err.Error
	}

	return nil
}

func (session *apiRepository) DecreaseUserBalance(userID string, amount int) error {
	walletExist := session.GetUserWallet(userID)

	if walletExist.UserID != userID {
		log.Error("Repository Error", "User wallet not found.")
		return errors.New("wallet not found")
	}

	walletExist.Balance -= int64(amount)
	if walletExist.Balance < 0 {
		log.Error("User wallet is not enough.")
		return errors.New("user can not pay with wallet")
	}
	err := session.db.Model(&entity.Wallet{}).Where("user_id = ?", userID).Update("balance", walletExist.Balance)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
	}

	return nil
}

func (session *apiRepository) CreateUserWallet(userID string) error {
	var wallet entity.Wallet

	walletExist := session.GetUserWallet(userID)
	if walletExist.UserID == userID {
		log.Error("Repository Error", "User wallet not found.")
		return errors.New("wallet found")
	}

	wallet.UserID = userID
	wallet.Balance = 0

	err := session.db.Save(&wallet)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
	}

	return nil
}

func (session *apiRepository) GetUserWallet(userID string) entity.Wallet {
	var wallet entity.Wallet

	err := session.db.Model(&entity.Wallet{}).Where("user_id = ?", userID).Take(&wallet)
	if err.Error != nil {
		log.Error("Repository Error", err.Error)
	}

	return wallet
}
