package service

import (
	"errors"
	"fmt"
	entity2 "github.com/itpourya/Haze/internal/entity"
	"github.com/itpourya/Haze/internal/repository"
	"hash/fnv"
	"strconv"

	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	GenerateUsernameService(userID string) string
	GetUserByUserIDService(userID string) []entity2.User
	IncreaseUserBalanceService(userID string, charge int) error
	CreateUserWalletService(userID string) error
	GetUserWalletService(userID string) entity2.Wallet
	DicreaseUserBalanceService(userID string, amount int) error
	CreateManagerService(userID string) error
	GetManagerService(userID string) entity2.Manager
	IncreaseManagerDeptService(userID string, price int64) bool
	GetInvoiceService(userID string) int64
	GetManagerListService() []entity2.Manager
	ClearManagerDeptService(userID string) bool
	EnterConfigOwnerNameService(name string, userID string) bool
}

type userService struct {
	repo repository.APIrepository
}

func NewUserService(repo repository.APIrepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (service *userService) EnterConfigOwnerNameService(name string, userID string) bool {
	if name == "" && userID == "" {
		return false
	}

	service.repo.EnterConfigOwnerName(name, userID)
	return true
}

func (service *userService) ClearManagerDeptService(userID string) bool {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return false
	}

	ok := service.repo.ClearManagerDept(userID)

	return ok
}

func (service *userService) GetManagerListService() []entity2.Manager {
	manager := service.repo.GetManagerList()

	return manager
}

func (service *userService) GetInvoiceService(userID string) int64 {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return 0
	}
	dept := service.repo.GetInvoice(userID)

	return dept
}

func (service *userService) IncreaseManagerDeptService(userID string, price int64) bool {
	if userID == "" || price <= 0 {
		log.Error("Service Error", "userID or price is not acceptable")
		return false
	}
	ok := service.repo.IncreaseManagerDept(userID, price)

	if !ok {
		log.Error("Service Error", "Failed increasing manager dept")
		return false
	}

	return ok
}

func (service *userService) GetManagerService(userID string) entity2.Manager {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return entity2.Manager{}
	}

	manager := service.repo.GetManager(userID)

	return manager
}

func (service *userService) CreateManagerService(userID string) error {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return errors.New("UserIS is empty")
	}
	service.repo.CreateManager(userID)

	return nil
}

func (service *userService) GenerateUsernameService(userID string) string {
	secret := strconv.Itoa(int(timestamppb.Now().Seconds))
	hashpwd := fmt.Sprint(hash(userID + secret))
	service.repo.RegisterUser(userID, string(hashpwd))

	return hashpwd
}

func (service *userService) GetUserByUserIDService(userID string) []entity2.User {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return nil
	}

	user := service.repo.GetUserConfigsAccount(userID)

	return user
}

func (service *userService) GetUserWalletService(userID string) entity2.Wallet {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return entity2.Wallet{}
	}
	userWallet := service.repo.GetUserWallet(userID)

	return userWallet
}

func (service *userService) IncreaseUserBalanceService(userID string, charge int) error {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return errors.New("UserID is empty")
	}
	err := service.repo.IncreaseUserBalance(userID, charge)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) DicreaseUserBalanceService(userID string, amount int) error {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return errors.New("UserID is empty")
	}

	err := service.repo.DecreaseUserBalance(userID, amount)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) CreateUserWalletService(userID string) error {
	if userID == "" {
		log.Error("Service Error", "UserID is empty")
		return errors.New("UserID is empty")
	}
	err := service.repo.CreateUserWallet(userID)
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
