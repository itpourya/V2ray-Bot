package marzban

type Token struct {
	AccessToen string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

const (
	API_AUTH_URL     = "https://marz.ikernel.ir:8000/api/admin/token"
	API_CREATE_USER  = "https://marz.ikernel.ir:8000/api/user"
	API_GET_USER     = "https://marz.ikernel.ir:8000/api/user/"
	DATA_LIMIT_10GB  = 10737418240
	DATA_LIMIT_15GB  = 16106127360
	DATA_LIMIT_20GB  = 21474836480
	DATA_LIMIT_50GB  = 53687091200
	DATA_LIMIT_100GB = 107374182400
)
