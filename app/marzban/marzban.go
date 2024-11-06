package marzban

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/itpourya/Haze/app/serializer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Marzban interface {
	CreateUserAccount(userID string, data_limit int) (serializer.Response, error)
	GetUser(userID string) (serializer.Response, string, error)
	ExpireUpdate(userID string) error
	DataLimitUpdate(userID string, charge string) error
}

type marzban struct{}

func NewMarzbanClient() Marzban {
	return &marzban{}
}

func (m *marzban) CreateUserAccount(username string, data_limit int) (serializer.Response, error) {
	expire := fmt.Sprint(CreateTime())
	limit := strconv.Itoa(GenerateData(data_limit))
	var resp *http.Response
	var response serializer.Response

	userFound, token, err := m.GetUser(username)
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
	json.Unmarshal(body, &response)

	defer resp.Body.Close()
	return response, nil
}

func (m *marzban) GetUser(username string) (serializer.Response, string, error) {
	var resp *http.Response
	var response serializer.Response
	token, err := auth()
	if err != nil {
		log.Println(err)
		return response, "", err
	}

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
		return response, token, errors.New("FAILD REQUEST | " + API_GET_USER)
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
	user, token, _ := m.GetUser(userID)
	expire := chargeMonth(user.Expire)
	client := &http.Client{}

  user.Expire = expire
  user.Status = "active"

  data := strings.NewReader(fmt.Sprint(user))

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
	defer resp.Body.Close()

	return nil
}

func (m *marzban) DataLimitUpdate(username string, charge string) error {
	var resp *http.Response
	user, token, _ := m.GetUser(username)
	client := &http.Client{}

  user.DataLimit = chargeDataLimit(user.DataLimit, charge)
  user.Status = "active"

  data := strings.NewReader(fmt.Sprint(user))

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

	if charge == "50GB" {
		return dataLimit + DATA_LIMIT_50GB
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
	defer resp.Body.Close()

	return jsonData.AccessToen, nil
}

func CreateTime() int64 {
	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)
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
