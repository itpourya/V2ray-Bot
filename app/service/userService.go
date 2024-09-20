package service

import (
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/itpourya/Haze/app/entity"
	"github.com/itpourya/Haze/app/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	CreateUsername(userID string) string
	GetUserByUserID(userID string) []entity.User
	IncreaseUserBalance(userID string, charge int) error
	CreateUserWallet(userID string) error
	GetUserWallet(userID string) entity.Wallet
	DicreaseUserBalance(userID string, amount int) error
}

type userService struct {
	repo repository.APIrepository
}

func NewUserService(repo repository.APIrepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (service *userService) CreateUsername(userID string) string {
	secret := strconv.Itoa(int(timestamppb.Now().Seconds))
	hashpwd := fmt.Sprint(hash(userID + secret))
	service.repo.CreateUser(userID, string(hashpwd))

	return hashpwd
}

func (service *userService) GetUserByUserID(userID string) []entity.User {
	user := service.repo.GetUser(userID)

	return user
}

func (service *userService) GetUserWallet(userID string) entity.Wallet {
	userWallet := service.repo.GetWallet(userID)

	return userWallet
}

func (service *userService) IncreaseUserBalance(userID string, charge int) error {
	err := service.repo.IncreseBalance(userID, charge)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) DicreaseUserBalance(userID string, amount int) error {
	err := service.repo.DecreaseBalance(userID, amount)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) CreateUserWallet(userID string) error {
	err := service.repo.CreateWallet(userID)
	if err != nil {
		return err
	}

	return nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
