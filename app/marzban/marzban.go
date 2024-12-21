package marzban

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/itpourya/Haze/app/cache"
	"github.com/itpourya/Haze/app/serializer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Marzban interface {
	CreateMarzbanUser(userID string, dataLimit int, dateLimit string) (serializer.Response, error)
	GetMarzbanUser(userID string) (serializer.Response, string, error)
	ExpireUpdate(userID string) error
	DataLimitUpdate(userID string, charge string) error
}

var (
	ctxb = context.Background()
	rdp  = cache.NewCache()
)

type marzban struct{}

func NewMarzbanClient() Marzban {
	return &marzban{}
}

func (m *marzban) CreateMarzbanUser(username string, dataLimit int, dateLimit string) (serializer.Response, error) {
	expire := fmt.Sprint(CreateTime(dateLimit))
	limit := strconv.Itoa(GenerateData(dataLimit))
	var resp *http.Response
	var response serializer.Response

	userFound, token, err := m.GetMarzbanUser(username)
	if err != nil {
		return response, err
	}
	if userFound.Status == "active" {
		return userFound, nil
	}

	data := strings.NewReader(`{
	  "username": ` + username + `,
	  "proxies": {
	    "vless": ""
	  },
	  "expire": ` + expire + `,
	  "data_limit": ` + limit + `,
	  "data_limit_reset_strategy": "no_reset",
	  "status": "active",
	  "note": "",
	  "on_hold_timeout": "2023-11-03T20:30:00",
	  "on_hold_expire_duration": 0
	}`)
	req, err := http.NewRequest("POST", API_CREATE_USER, data)
	if err != nil {
		return response, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if resp == nil {
		return response, errors.New("FAILED REQUEST | " + API_CREATE_USER)
	}
	if err != nil {
		return response, err
	}

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	return response, nil
}

func (m *marzban) GetMarzbanUser(username string) (serializer.Response, string, error) {
	var resp *http.Response
	var response serializer.Response
	var token string

	token, _ = auth()
	req, err := http.NewRequest("GET", API_GET_USER+username, nil)
	if err != nil {
		log.Println(err)
		return response, "", err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	resp, err = client.Do(req)
	if resp == nil {
		return response, token, errors.New("FAILED REQUEST | " + API_GET_USER)
	}
	if err != nil {
		log.Println(err)
		return response, "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return response, "", err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return response, "", err
	}

	return response, token, nil
}

func (m *marzban) ExpireUpdate(userID string) error {
	var resp *http.Response
	user, token, _ := m.GetMarzbanUser(userID)
	expire := fmt.Sprint(chargeMonth(user.Expire))
	client := &http.Client{}

	data := strings.NewReader(`{
			 "proxies": {
			   "vless": {}
			 },
			 "inbounds": {
			   "vless": []
			 },
			 "expire": ` + expire + `,
			 "data_limit": ` + fmt.Sprint(user.DataLimit) + `,
			 "data_limit_reset_strategy": "no_reset",
			 "status": "active",
			 "note": "",
			 "on_hold_timeout": "2023-11-03T20:30:00",
			 "on_hold_expire_duration": 0
			}`)

	req, err := http.NewRequest("PUT", API_GET_USER+userID, data)
	if err != nil {
		return errors.New("FAILED CREATE PUT REQUEST | " + API_GET_USER + " " + err.Error())
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if resp == nil {
		return errors.New("FAILED PUT REQUEST | " + API_GET_USER)
	}
	if err != nil {
		return errors.New("FAILED PUT REQUEST | " + API_GET_USER + " " + err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	return nil
}

func (m *marzban) DataLimitUpdate(username string, charge string) error {
	var resp *http.Response
	user, token, _ := m.GetMarzbanUser(username)
	client := &http.Client{}

	data := strings.NewReader(`{
			 "proxies": {
			   "vless": {}
			 },
			 "inbounds": {
			   "vless": []
			 },
			 "expire": ` + fmt.Sprint(user.Expire) + `,
			 "data_limit": ` + fmt.Sprint(chargeDataLimit(user.DataLimit, charge)) + `,
			 "data_limit_reset_strategy": "no_reset",
			 "status": "active",
			 "note": "",
			 "on_hold_timeout": "2023-11-03T20:30:00",
			 "on_hold_expire_duration": 0
			}`)

	req, err := http.NewRequest("PUT", API_GET_USER+username, data)
	if err != nil {
		return errors.New("FAILED CREATE PUT REQUEST | " + API_GET_USER + " " + err.Error())
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if resp == nil {
		return errors.New("FAILED PUT REQUEST | " + API_GET_USER)
	}
	if err != nil {
		return errors.New("FAILED PUT REQUEST | " + API_GET_USER + " " + err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("FAILED CLOSE REQUEST | " + API_GET_USER)
		}
	}(resp.Body)

	return nil
}

func chargeDataLimit(dataLimit int64, charge string) int64 {
	if charge == "10GB" {
		return dataLimit + DATA_LIMIT_10GB
	}

	if charge == "15GB" {
		return dataLimit + DATA_LIMIT_15GB
	}

	if charge == "20GB" {
		return dataLimit + DATA_LIMIT_20GB
	}

	if charge == "30GB" {
		return dataLimit + DATA_LIMIT_30GB
	}

	if charge == "40GB" {
		return dataLimit + DATA_LIMIT_40GB
	}

	if charge == "50GB" {
		return dataLimit + DATA_LIMIT_50GB
	}

	if charge == "70GB" {
		return dataLimit + DATA_LIMIT_70GB
	}

	if charge == "80GB" {
		return dataLimit + DATA_LIMIT_80GB
	}

	if charge == "90GB" {
		return dataLimit + DATA_LIMIT_90GB
	}

	if charge == "100GB" {
		return dataLimit + DATA_LIMIT_100GB
	}

	return dataLimit + DATA_LIMIT_100GB
}

func auth() (string, error) {
	var resp *http.Response
	payload := strings.NewReader(`grant_type=&username=admin&password=Poorya1382&scope=&client_id=&client_secret=`)
	req, err := http.NewRequest("POST", API_AUTH_URL, payload)
	if err != nil {
		return "", err
	}

	client := http.Client{}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json")

	resp, _ = client.Do(req)

	if resp == nil {
		return "", errors.New("nil response")
	}
	var jsonData Token

	err = json.NewDecoder(resp.Body).Decode(&jsonData)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var cacheAuthToken cache.CacheAuthToken
	cacheAuthToken.AuthToken = jsonData.AccessToen
	cacheMsg := rdp.Set(ctxb, "TOKEN", cacheAuthToken, time.Duration(1*time.Hour))
	if cacheMsg != nil {
		log.Println(err)
	}

	return jsonData.AccessToen, nil
}

func CreateTime(month string) int64 {
	now := time.Now()
	var futureDate time.Time

	if month == "1" {
		futureDate = now.AddDate(0, 1, 0)
	}
	if month == "2" {
		futureDate = now.AddDate(0, 2, 0)
	}
	if month == "3" {
		futureDate = now.AddDate(0, 3, 0)
	}
	if month == "4" {
		futureDate = now.AddDate(0, 4, 0)
	}
	if month == "5" {
		futureDate = now.AddDate(0, 5, 0)
	}
	if month == "6" {
		futureDate = now.AddDate(0, 6, 0)
	}

	timestamp := time.Date(futureDate.Year(), futureDate.Month(), futureDate.Day(), 0, 0, 0, 0, time.UTC)

	t := timestamppb.New(timestamp).Seconds
	return t
}

func GenerateData(dataLimit int) int {
	if dataLimit == 10 {
		return DATA_LIMIT_10GB
	}

	if dataLimit == 15 {
		return DATA_LIMIT_15GB
	}

	if dataLimit == 20 {
		return DATA_LIMIT_20GB
	}

	if dataLimit == 30 {
		return DATA_LIMIT_30GB
	}

	if dataLimit == 40 {
		return DATA_LIMIT_40GB
	}

	if dataLimit == 60 {
		return DATA_LIMIT_60GB
	}

	if dataLimit == 70 {
		return DATA_LIMIT_70GB
	}

	if dataLimit == 80 {
		return DATA_LIMIT_80GB
	}

	if dataLimit == 90 {
		return DATA_LIMIT_90GB
	}

	if dataLimit == 50 {
		return DATA_LIMIT_50GB
	}

	if dataLimit == 100 {
		return DATA_LIMIT_100GB
	}

	return 0
}

func chargeMonth(expire int64) int64 {
	t := time.Unix(expire, 0)
	t = t.AddDate(0, 1, 0)
	timestamp := timestamppb.New(t).Seconds

	return timestamp
}
