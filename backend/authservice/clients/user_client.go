package clients

import (
	config "authservice/configs"
	"authservice/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserClient interface {
	Register(*models.UserRequest) error
	GetUser(user *models.UserLoginRequest) (map[string]interface{}, error)
}

type UserClientImpl struct {
}

func NewUserClient() UserClient {
	return &UserClientImpl{}
}

func (client *UserClientImpl) Register(user *models.UserRequest) error {
	appConfig := config.GetConfig()
	jsonData, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPost, appConfig.UserSvcUrl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	return nil
}

func (client *UserClientImpl) GetUser(user *models.UserLoginRequest) (map[string]interface{}, error) {

	appConfig := config.GetConfig()
	jsonData, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPost, appConfig.UserSvcUrl+"/find", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if resp.StatusCode != 200 {
		return nil, errors.New("Wrong username or password")
	}
	var response models.Response
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)
	m, _ := response.Data.(map[string]interface{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// defer response.Body.Close()
	return m, nil
}
