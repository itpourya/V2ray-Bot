package cache

import (
	"encoding/json"
)

type CachePayload struct {
	ConfigName      string
	Remonth         bool
	Recharge        bool
	Buy             bool
	DataLimitCharge string
	ReChargeWallet  bool
	ChargeWallet    string
	Price           int64
	Manager         bool
	DateLimit       string
}

type CacheAuthToken struct {
	AuthToken string
}

func (i CachePayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i CacheAuthToken) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
