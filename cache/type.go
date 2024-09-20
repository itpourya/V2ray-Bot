package cache

import "encoding/json"

type CachePayload struct {
	ConfigName      string
	Remonth         bool
	Recharge        bool
	Buy             bool
	DataLimitCharge string
	ReChargeWallet  bool
	ChargeWallet    string
	Price           int64
}

func (i CachePayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
