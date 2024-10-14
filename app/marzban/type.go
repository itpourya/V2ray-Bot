package marzban

import (
	"os"

	"github.com/joho/godotenv"
)

type Token struct {
	AccessToen string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

var (
	_               = godotenv.Load(".env")
	API_AUTH_URL    = os.Getenv("API_AUTH_URL")
	API_CREATE_USER = os.Getenv("API_CREATE_USER")
	API_GET_USER    = os.Getenv("API_GET_USER")
)

const (
	DATA_LIMIT_10GB  = 10737418240
	DATA_LIMIT_15GB  = 16106127360
	DATA_LIMIT_20GB  = 21474836480
	DATA_LIMIT_50GB  = 53687091200
	DATA_LIMIT_100GB = 107374182400
)
