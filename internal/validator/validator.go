package validator

import (
	"encoding/json"
	"github.com/itpourya/Haze/internal/cache"
	"strings"
)

func Validate(data string, userID string) cache.CachePayload {
	var payload cache.CachePayload
	data2 := strings.Replace(data, "getdel "+userID+": ", "", 1)
	data2 = strings.Replace(data2, "get "+userID+": ", "", 1)
	err := json.Unmarshal([]byte(data2), &payload)
	if err != nil {
		return cache.CachePayload{}
	}

	return payload
}

func ValidateAuth(data string) cache.CacheAuthToken {
	var payload cache.CacheAuthToken
	data2 := strings.Replace(data, "get TOKEN: ", "", 1)
	err := json.Unmarshal([]byte(data2), &payload)
	if err != nil {
		return cache.CacheAuthToken{}
	}

	return payload
}
