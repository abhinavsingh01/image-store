package tests

import (
	"albumservice/client"
	config "albumservice/configs"
	"albumservice/controllers"
	"albumservice/models"
	"albumservice/server"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"albumservice/middleware"
	"albumservice/repository"

	"albumservice/services"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.uber.org/dig"
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

type MockImageClient struct {
	mock.Mock
}

func NewMockImageClient() client.ImageClient {
	return &MockImageClient{}
}

func (m *MockImageClient) DeleteImage(userId int, albumId string) error {
	m.Called(userId, albumId)
	return nil
}

func addNewAlbum(router *gin.Engine) string {
	album := models.AlbumRequest{
		AlbumName: "test",
	}
	jsonData, _ := json.Marshal(album)
	request, _ := http.NewRequest(http.MethodPost, "/v1/album", bytes.NewBuffer(jsonData))
	request.Header.Add("user-id", "1")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, request)
	var response models.Response
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)
	m, _ := response.Data.(map[string]interface{})
	albumId := fmt.Sprintf("%v", m["albumId"])
	return albumId
}

func testDBINit() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Album{})
	return db, nil
}

func setupTestContainer() *dig.Container {
	mockImageClient := new(MockImageClient)
	mockImageClient.On("DeleteImage", 1, 3).Return(nil)
	container := dig.New()
	container.Provide(config.LoadConfig)
	container.Provide(testDBINit)
	container.Provide(func() client.ImageClient {
		return mockImageClient
	})

	container.Provide(repository.NewAlbumRepoImpl)
	container.Provide(services.NewAlbumServiceImpl)
	container.Provide(controllers.NewAlbumController)
	container.Provide(middleware.NewAuth)
	container.Provide(server.NewRouter)

	return container
}

var _ = Describe("AlbumController", func() {
	var router *gin.Engine
	var albumId string
	BeforeSuite(func() {
		cont := setupTestContainer()
		err := cont.Invoke(func(r *server.Router) {
			router = r.InitRouter()
		})
		fmt.Println(err)

	})

	Describe("Get all albums of user", func() {
		Context("If there is no album for user", func() {
			BeforeEach(func() {

			})
			It("should return No album found message", func() {
				request, err := http.NewRequest(http.MethodGet, "/v1/album/all", nil)
				request.Header.Add("user-id", "1")
				Expect(err).ToNot(HaveOccurred())
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				var response models.Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
				Expect(response.Message).To(Equal("No album found"))
			})
		})
		Context("If there is some album for user", func() {
			BeforeEach(func() {
				albumId = addNewAlbum(router)
			})
			It("should return No album found message", func() {
				request, err := http.NewRequest(http.MethodGet, "/v1/album/all", nil)
				request.Header.Add("user-id", "1")
				Expect(err).ToNot(HaveOccurred())
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				var response models.Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
				Expect(response.Data).ShouldNot(BeEmpty())
			})
		})
	})
	Describe("Get album for user", func() {
		Context("If we pass correct album id", func() {
			BeforeEach(func() {
				albumId = addNewAlbum(router)
			})
			It("should return album", func() {
				request, err := http.NewRequest(http.MethodGet, "/v1/album/"+albumId, nil)
				request.Header.Add("user-id", "1")
				Expect(err).ToNot(HaveOccurred())
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				var response models.Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
			})

		})
		Context("If we pass wrong album id", func() {
			It("should return error", func() {
				request, err := http.NewRequest(http.MethodGet, "/v1/album/wrong", nil)
				request.Header.Add("user-id", "1")
				Expect(err).ToNot(HaveOccurred())
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(400))
			})
		})
	})

	Describe("Delete album for user", func() {
		Context("If we pass correct album id", func() {
			BeforeEach(func() {
				albumId = addNewAlbum(router)
			})
			It("should delete album", func() {

				request, err := http.NewRequest(http.MethodDelete, "/v1/album/"+albumId, nil)
				request.Header.Add("user-id", "1")
				Expect(err).ToNot(HaveOccurred())
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				var response models.Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
			})

		})
		Context("If we pass wrong album id", func() {
			It("should return error", func() {
				request, err := http.NewRequest(http.MethodDelete, "/v1/album/wrong", nil)
				request.Header.Add("user-id", "1")
				Expect(err).ToNot(HaveOccurred())
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(400))
			})
		})
	})

	Describe("Create new album for user", func() {
		Context("If we pass correct payload", func() {
			It("should return album id", func() {
				album := models.AlbumRequest{
					AlbumName: "test",
				}
				jsonData, _ := json.Marshal(album)
				request, _ := http.NewRequest(http.MethodPost, "/v1/album", bytes.NewBuffer(jsonData))
				request.Header.Add("user-id", "1")
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(200))
			})
		})

		Context("If we pass wrong payload", func() {
			It("should return error", func() {
				album := models.AlbumRequest{}
				jsonData, _ := json.Marshal(album)
				request, _ := http.NewRequest(http.MethodPost, "/v1/album", bytes.NewBuffer(jsonData))
				request.Header.Add("user-id", "1")
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(400))
			})
		})
	})

})
