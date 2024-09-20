package validator

import (
	"encoding/json"
	"strings"

	"github.com/itpourya/Haze/app/cache"
)

func Validate(data string, userID string) cache.CachePayload {
	var payload cache.CachePayload
	data2 := strings.Replace(data, "getdel "+userID+": ", "", 1)
	data2 = strings.Replace(data2, "get "+userID+": ", "", 1)
	json.Unmarshal([]byte(data2), &payload)

	return payload
}
