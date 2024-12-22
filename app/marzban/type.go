package marzban

type Token struct {
	AccessToen string `json:"access_token"`
	TokenType  string `json:"token_type"`
}

var (
	API_AUTH_URL    = "https://marz.redzedshop.ir:8000/api/admin/token"
	API_CREATE_USER = "https://marz.redzedshop.ir:8000/api/user"
	API_GET_USER    = "https://marz.redzedshop.ir:8000/api/user/"
)

const (
	DATA_LIMIT_10GB  = 10737418240
	DATA_LIMIT_15GB  = 16106127360
	DATA_LIMIT_20GB  = 21474836480
	DATA_LIMIT_50GB  = 53687091200
	DATA_LIMIT_100GB = 107374182400
	DATA_LIMIT_30GB  = 32212254720
	DATA_LIMIT_40GB  = 42949672960
	DATA_LIMIT_60GB  = 64424509440
	DATA_LIMIT_70GB  = 75161927680
	DATA_LIMIT_80GB  = 85899345920
	DATA_LIMIT_90GB  = 96636764160
)
