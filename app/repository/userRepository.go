package repository

import (
	"errors"
	"log"

	"github.com/itpourya/Haze/app/entity"
	"gorm.io/gorm"
)

type APIrepository interface {
	CreateUser(userID string, userSub string) error
	GetUser(userID string) []entity.User
	IncreseBalance(userID string, charge int) error
	CreateWallet(userID string) error
	GetWallet(userID string) entity.Wallet
	DecreaseBalance(userID string, amount int) error
}

type apiRepository struct {
	db *gorm.DB
}

func NewRepository(conn *gorm.DB) APIrepository {
	return &apiRepository{
		db: conn,
	}
}

func (marz *apiRepository) CreateUser(userID string, userSub string) error {
	var user entity.User
	user.UserID = userID
	user.UsernameSub = userSub

	marz.db.Save(&user)

	return nil
}

func (marz *apiRepository) GetUser(userID string) []entity.User {
	var user []entity.User

	_ = marz.db.Model(&entity.User{}).Where("user_id == ?", userID).Find(&user)

	return user
}

func (marz *apiRepository) IncreseBalance(userID string, charge int) error {
	walletExist := marz.GetWallet(userID)

	if walletExist.UserID != userID {
		marz.CreateWallet(userID)
		walletExist = marz.GetWallet(userID)
	}

	walletExist.Balance += int64(charge)
	err := marz.db.Model(&entity.Wallet{}).Where("user_id = ?", userID).Update("balance", walletExist.Balance)
	if err != nil {
		return err.Error
	}

	return nil
}

func (marz *apiRepository) DecreaseBalance(userID string, amount int) error {
	walletExist := marz.GetWallet(userID)

	if walletExist.UserID != userID {
		return errors.New("wallet not found")
	}

	walletExist.Balance -= int64(amount)
	if walletExist.Balance < 0 {
		return errors.New("user can not pay with wallet")
	}
	err := marz.db.Model(&entity.Wallet{}).Where("user_id = ?", userID).Update("balance", walletExist.Balance)
	if err != nil {
		return err.Error
	}

	return nil
}

func (marz *apiRepository) CreateWallet(userID string) error {
	var wallet entity.Wallet

	walletExist := marz.GetWallet(userID)
	if walletExist.UserID == userID {
		return errors.New("wallet found")
	}

	wallet.UserID = userID
	wallet.Balance = 0

	marz.db.Save(&wallet)

	return nil
}

func (marz *apiRepository) GetWallet(userID string) entity.Wallet {
	var wallet entity.Wallet

	err := marz.db.Model(&entity.Wallet{}).Where("user_id = ?", userID).Take(&wallet)
	if err != nil {
		log.Println(err)
	}

	return wallet
}
