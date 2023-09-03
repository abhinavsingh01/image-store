package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image-store-api-gateway/middleware"
	"image-store-api-gateway/server"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type MockAuth struct {
	mock.Mock
}

func NewMockAuth() middleware.Auth {
	return &MockAuth{}
}

func (m *MockAuth) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Add("user-id", "1")
	}
}

type ResponseRecorderWrapper struct {
	*httptest.ResponseRecorder
	closeNotifyChan chan bool
}

func (rw *ResponseRecorderWrapper) CloseNotify() <-chan bool {
	if rw.closeNotifyChan == nil {
		rw.closeNotifyChan = make(chan bool, 1)
	}
	return rw.closeNotifyChan
}

func createNewUser(router *gin.Engine) string {
	user := UserRequest{
		Username: "mocktest1",
		Name:     "mocktest1",
		Password: "mocktest",
		Email:    "mocktest@test.com",
	}
	jsonData, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPost, "/gw/api/v1/auth/register", bytes.NewBuffer(jsonData))
	request.RequestURI = "/gw/api/v1/auth/register"
	resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
	router.ServeHTTP(resp, request)
	return loginUser(router)
}

func loginUser(router *gin.Engine) string {
	user := UserLogin{
		Username: "mocktest1",
		Password: "mocktest",
	}
	jsonData, _ := json.Marshal(user)
	request, _ := http.NewRequest(http.MethodPost, "/gw/api/v1/auth/login", bytes.NewBuffer(jsonData))
	request.RequestURI = "/gw/api/v1/auth/login"
	resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
	router.ServeHTTP(resp, request)
	var response Response
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)
	m, _ := response.Data.(map[string]interface{})
	token := fmt.Sprintf("%v", m["token"])
	return token
}

func addNewAlbum(router *gin.Engine, token string) string {
	album := AlbumRequest{
		AlbumName: "test",
	}
	jsonData, _ := json.Marshal(album)
	request, _ := http.NewRequest(http.MethodPost, "/gw/api/v1/album/new", bytes.NewBuffer(jsonData))
	request.RequestURI = "/gw/api/v1/album/new"
	request.Header.Set("Authorization", token)
	resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
	router.ServeHTTP(resp, request)
	var response Response
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &response)
	m, _ := response.Data.(map[string]interface{})
	albumId := fmt.Sprintf("%v", m["albumId"])
	return albumId
}

func uploadImage(router *gin.Engine, token string, albumId string) string {
	fileDir, _ := os.Getwd()
	filePath := filepath.Join(fileDir, "test.jpg")
	file, _ := os.Open(filePath)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("files", filepath.Base(filePath))
	io.Copy(part, file)
	writer.Close()
	url := "/gw/api/v1/image/album/" + albumId + "/image/upload"
	request, _ := http.NewRequest("POST", url, body)
	request.RequestURI = url
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)
	resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
	router.ServeHTTP(resp, request)
	g := resp.Result()
	var response Response
	bodyBytes, _ := ioutil.ReadAll(g.Body)
	json.Unmarshal(bodyBytes, &response)
	data, _ := response.Data.(([]interface{}))
	imageId := data[0].(string)
	return imageId
}

func deleteAlbum(router *gin.Engine, token string, albumId string) {
	request, _ := http.NewRequest(http.MethodDelete, "/gw/api/v1/album/"+albumId, nil)
	request.Header.Set("Authorization", token)
	request.RequestURI = "/gw/api/v1/album/" + albumId
	resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
	router.ServeHTTP(resp, request)
}

func deleteUser(router *gin.Engine) {
	// impelement delete user after
}

var _ = Describe("Integration testing", func() {
	var router *gin.Engine
	var token string
	var albumId string
	// auth := new(MockAuth)
	auth := middleware.NewAuth()
	BeforeSuite(func() {
		gin.SetMode(gin.TestMode)
		router = server.InitRoutes(auth)
		token = "Bearer " + loginUser(router)
	})

	Describe("Get all albums of user", func() {
		Context("If there is no album for user", func() {
			BeforeEach(func() {

			})
			It("should return No album found message", func() {
				request, err := http.NewRequest(http.MethodGet, "/gw/api/v1/album/all", nil)
				request.RequestURI = "/gw/api/v1/album/all"
				request.Header.Set("Authorization", token)
				Expect(err).ToNot(HaveOccurred())
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				var response Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
				Expect(response.Message).To(Equal("No album found"))
			})
		})
		Context("If there is some album for user", func() {
			BeforeEach(func() {
				albumId = addNewAlbum(router, token)
			})
			It("should return album list", func() {
				request, err := http.NewRequest(http.MethodGet, "/gw/api/v1/album/all", nil)
				request.RequestURI = "/gw/api/v1/album/all"
				request.Header.Set("Authorization", token)
				Expect(err).ToNot(HaveOccurred())
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				var response Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
				Expect(response.Data).ShouldNot(BeEmpty())
			})
			AfterEach(func() {
				deleteAlbum(router, token, albumId)
			})
		})
	})

	Describe("Create new album for user", func() {
		Context("If we pass correct payload", func() {
			It("should return album id", func() {
				album := AlbumRequest{
					AlbumName: "test",
				}
				jsonData, _ := json.Marshal(album)
				request, _ := http.NewRequest(http.MethodPost, "/gw/api/v1/album/new", bytes.NewBuffer(jsonData))
				request.Header.Set("Authorization", token)
				request.RequestURI = "/gw/api/v1/album/new"
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(200))
			})
			AfterEach(func() {
				deleteAlbum(router, token, albumId)
			})
		})

		Context("If we pass wrong payload", func() {
			It("should return error", func() {
				album := AlbumRequest{}
				jsonData, _ := json.Marshal(album)
				request, _ := http.NewRequest(http.MethodPost, "/gw/api/v1/album/new", bytes.NewBuffer(jsonData))
				request.Header.Set("Authorization", token)
				request.RequestURI = "/gw/api/v1/album/new"
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(400))
			})
		})
	})

	Describe("Delete album", func() {
		Context("If we pass correct id", func() {
			BeforeEach(func() {
				albumId = addNewAlbum(router, token)
			})
			It("should return success", func() {
				request, _ := http.NewRequest(http.MethodDelete, "/gw/api/v1/album/"+albumId, nil)
				request.Header.Set("Authorization", token)
				request.RequestURI = "/gw/api/v1/album/" + albumId
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(200))
			})
			AfterEach(func() {
				deleteAlbum(router, token, albumId)
			})
		})
		Context("If we pass wrong id", func() {
			It("should return error", func() {
				url := "/gw/api/v1/album/wrong"
				request, _ := http.NewRequest(http.MethodDelete, url, nil)
				request.Header.Set("Authorization", token)
				request.RequestURI = url
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(400))
			})
		})
	})

	Describe("Get all images of an album", func() {
		BeforeEach(func() {
			albumId = addNewAlbum(router, token)
		})
		Context("if there are no images", func() {
			It("should return empty images", func() {
				url := "/gw/api/v1/image/album/" + albumId + "/images"
				request, err := http.NewRequest(http.MethodGet, url, nil)
				request.RequestURI = url
				request.Header.Set("Authorization", token)
				Expect(err).ToNot(HaveOccurred())
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				var response Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
				Expect(response.Message).To(Equal("No images found"))
			})
		})
		Context("If there are images", func() {
			BeforeEach(func() {
				uploadImage(router, token, albumId)
			})
			It("should return images", func() {
				url := "/gw/api/v1/image/album/" + albumId + "/images"
				request, err := http.NewRequest(http.MethodGet, url, nil)
				request.RequestURI = url
				request.Header.Set("Authorization", token)
				Expect(err).ToNot(HaveOccurred())
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				var response Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(200))
				Expect(response.Message).To(Equal("Images found"))
			})
		})
		Context("If we pass wrong album id", func() {
			BeforeEach(func() {
				uploadImage(router, token, albumId)
			})
			It("should return error", func() {
				url := "/gw/api/v1/image/album/wrong/images"
				request, err := http.NewRequest(http.MethodGet, url, nil)
				request.RequestURI = url
				request.Header.Set("Authorization", token)
				Expect(err).ToNot(HaveOccurred())
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				var response Response
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &response)
				Expect(resp.Code).To(Equal(400))
			})
		})
		AfterEach(func() {
			deleteAlbum(router, token, albumId)
		})
	})

	Describe("Upload image", func() {
		Context("if we pass correct album id", func() {
			BeforeEach(func() {
				albumId = addNewAlbum(router, token)
			})
			It("should return success", func() {
				fileDir, _ := os.Getwd()
				filePath := filepath.Join(fileDir, "test.jpg")
				file, _ := os.Open(filePath)
				defer file.Close()
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("files", filepath.Base(filePath))
				io.Copy(part, file)
				writer.Close()
				url := "/gw/api/v1/image/album/" + albumId + "/image/upload"
				request, _ := http.NewRequest("POST", url, body)
				request.RequestURI = url
				request.Header.Set("Content-Type", writer.FormDataContentType())
				request.Header.Set("Authorization", token)
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(200))
			})
			AfterEach(func() {
				deleteAlbum(router, token, albumId)
			})
		})
		Context("if we pass wrong album id", func() {
			It("should return error", func() {
				fileDir, _ := os.Getwd()
				filePath := filepath.Join(fileDir, "test.jpg")
				file, _ := os.Open(filePath)
				defer file.Close()
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("files", filepath.Base(filePath))
				io.Copy(part, file)
				writer.Close()
				url := "/gw/api/v1/image/album/wrong/image/upload"
				request, _ := http.NewRequest("POST", url, body)
				request.RequestURI = url
				request.Header.Set("Content-Type", writer.FormDataContentType())
				request.Header.Set("Authorization", token)
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(400))
			})
		})
	})

	Describe("Get image", func() {
		var imageId string
		BeforeEach(func() {
			albumId = addNewAlbum(router, token)
			imageId = uploadImage(router, token, albumId)
		})

		Context("if we pass correct album id", func() {
			It("should return success", func() {
				url := "/gw/api/v1/image/" + imageId + "/download"
				request, err := http.NewRequest(http.MethodGet, url, nil)
				request.RequestURI = url
				request.Header.Set("Authorization", token)
				Expect(err).ToNot(HaveOccurred())
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(200))
			})
		})
		Context("if we pass wrong image id", func() {
			It("should return error", func() {
				url := "/gw/api/v1/image/wrong/download"
				request, err := http.NewRequest(http.MethodGet, url, nil)
				request.RequestURI = url
				request.Header.Set("Authorization", token)
				Expect(err).ToNot(HaveOccurred())
				resp := &ResponseRecorderWrapper{ResponseRecorder: httptest.NewRecorder()}
				router.ServeHTTP(resp, request)
				Expect(resp.Code).To(Equal(400))
			})
		})
		AfterEach(func() {
			deleteAlbum(router, token, albumId)
		})
	})
})
