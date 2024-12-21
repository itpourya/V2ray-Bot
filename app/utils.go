package app

func IsManager(userID string) bool {
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
