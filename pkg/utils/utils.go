package utils

import (
	"github.com/itpourya/Haze/internal/service"
)

func IsManager(userID string, userService service.UserService) bool {
	manager := userService.GetManagerService(userID)

	if manager.UserID == userID {
		return true
	} else {
		return false
	}
}

func IsAdmin(userID string) bool {
	return userID == "6556338275"
}
