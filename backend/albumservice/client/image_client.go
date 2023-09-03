package client

import (
	config "albumservice/configs"
	"fmt"
	"net/http"
	"strconv"
)

type ImageClient interface {
	DeleteImage(userId int, albumId string) error
}

type ImageClientImpl struct {
}

func NewImageClient() ImageClient {
	return &ImageClientImpl{}
}

func (client *ImageClientImpl) DeleteImage(userId int, albumId string) error {
	appConfig := config.GetConfig()
	request, _ := http.NewRequest("DELETE", appConfig.ImageSvcUrl+"/"+albumId+"/images", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("user-id", strconv.Itoa(userId))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	return nil
}
