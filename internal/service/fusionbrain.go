package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/swenro11/stribog/config"
	"github.com/swenro11/stribog/internal/entity"
	log "github.com/swenro11/stribog/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	BaseURL              = "https://api-key.fusionbrain.ai/key/api/v1/"
	StylesURL            = "https://cdn.fusionbrain.ai/static/styles/api"
	ModelsAddURL         = "models"
	RunAddURL            = "text2image/run"
	GetAddURL            = "text2image/status/"
	Kandinsky3ModelId    = "4"
	emptyUuid            = "00000000-0000-0000-0000-000000000000"
	TaskStatusInitial    = "INITIAL"    // the request has been received, is in the queue for processing
	TaskStatusProcessing = "PROCESSING" // the request is being processed
	TaskStatusDone       = "DONE"       // task completed
	TaskStatusFail       = "FAIL"       // the task could not be completed.
	ErrorTaskNotFound    = "404 Not Found"
	LetterBytes          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type FusionbrainService struct {
	cfg *config.Config
	log *log.Logger
}

type ResponseRun struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
}

type RequestRunParams struct {
	Type                 string `json:"type"`
	NumImages            uint   `json:"numImages"`
	Width                uint   `json:"width"`
	Height               uint   `json:"height"`
	Style                string `json:"style,omitempty"`
	NegativePromptUnclip string `json:"negativePromptUnclip"`
	GenerateParams       struct {
		Query string `json:"query"`
	} `json:"generateParams"`
}

type ResponseModels struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Version float64 `json:"version"`
	Type    string  `json:"type"`
}

/*
	{
	  "uuid": "string",
	  "status": "string",
	  "images": ["string"],
	  "errorDescription": "string",
	  "censored": "false"
	}
*/
type ResponseStatus struct {
	Uuid             string   `json:"uuid"`
	Status           string   `json:"status"`
	Images           []string `json:"images"`
	ErrorDescription string   `json:"errorDescription"`
	Censored         bool     `json:"censored"`
}

func NewFusionbrainService(cfg *config.Config, l *log.Logger) *FusionbrainService {
	return &FusionbrainService{
		cfg: cfg,
		log: l,
	}
}

func (service *FusionbrainService) AuthGetRequest(addURL string) (*http.Response, error) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}

	req, err := http.NewRequest(http.MethodGet, BaseURL+addURL, nil)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.AuthGetRequest - http.NewRequest: " + err.Error())
	}

	req.Header.Add("X-Key", "Key "+service.cfg.AI.FusionbrainApi)
	req.Header.Add("X-Secret", "Secret "+service.cfg.AI.FusionbrainSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.AuthGetRequest - client.Do: " + err.Error())
	}

	return resp, nil
}

/*
Obsolete
https://freshman.tech/snippets/go/multipart-upload-google-drive/
data must be MultipartFormDataContent (sorry for my C# commentary in Go)
*/
func (service *FusionbrainService) AuthPostRequest(addURL string, data []byte) (*http.Response, error) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}

	req, err := http.NewRequest(http.MethodPost, BaseURL+addURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.AuthPostRequest - http.NewRequest: " + err.Error())
	}

	req.Header.Set("Content-Type", "multipart/related") //"application/json; charset=UTF-8"
	req.Header.Add("X-Key", "Key "+service.cfg.AI.FusionbrainApi)
	req.Header.Add("X-Secret", "Secret "+service.cfg.AI.FusionbrainSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.AuthPostRequest - client.Do: " + err.Error())
	}

	return resp, nil
}

// {\"id\":4,\"name\":\"Kandinsky\",\"version\":3.0,\"type\":\"TEXT2IMAGE\"}
func (service *FusionbrainService) GetStringModels() (string, error) {
	response, errAuthNewRequest := service.AuthGetRequest(ModelsAddURL)

	if errAuthNewRequest != nil {
		return "", fmt.Errorf("FusionbrainService.GetStringModels - AuthGetRequest: " + errAuthNewRequest.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("FusionbrainService.GetStringModels - ioutil.ReadAll: " + err.Error())
	}

	return string(body), nil
}

func (service *FusionbrainService) GetModels() (*ResponseModels, error) {
	response, errAuthNewRequest := service.AuthGetRequest(ModelsAddURL)

	if errAuthNewRequest != nil {
		return nil, fmt.Errorf("FusionbrainService.GetModels - AuthGetRequest: " + errAuthNewRequest.Error())
	}

	defer response.Body.Close()

	var models []*ResponseModels
	errDecode := json.NewDecoder(response.Body).Decode(&models)
	if errDecode != nil {
		return nil, fmt.Errorf("FusionbrainService.GetModels - Decode: " + errDecode.Error())
	}

	return models[0], nil
}

// TODO: quantity don't work. I send 5 -> get 1 image.
func (service *FusionbrainService) CreateTask(promt string, quantity uint, width uint, height uint, style string, negativePromptUnclip string, enableLog bool) (*ResponseRun, error) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}

	var requestData = RequestRunParams{
		Type:                 "GENERATE",
		NumImages:            quantity,
		Height:               height,
		Width:                width,
		Style:                style,
		NegativePromptUnclip: negativePromptUnclip,
		GenerateParams: struct {
			Query string "json:\"query\""
		}{promt},
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	paramsPart := make(map[string][]string)
	paramsPart["Content-Disposition"] = append(paramsPart["Content-Disposition"], "form-data; name=\"params\"")
	paramsPart["Content-Type"] = append(paramsPart["Content-Type"], "application/json")

	paramsWriter, err := writer.CreatePart(paramsPart)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - writer.CreatePart: ", err.Error())
	}

	paramsPayloadBytes, err := json.Marshal(&requestData)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - json.Marshal: ", err.Error())
	}

	_, err = paramsWriter.Write(paramsPayloadBytes)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - paramsWriter.Write: ", err.Error())
	}

	err = writer.WriteField("model_id", Kandinsky3ModelId)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - writer.WriteField: ", err.Error())
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - writer.Close: ", err.Error())
	}

	request, err := http.NewRequest(http.MethodPost, BaseURL+RunAddURL, payload)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - http.NewRequest: ", err.Error())
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Add("X-Key", "Key "+service.cfg.AI.FusionbrainApi)
	request.Header.Add("X-Secret", "Secret "+service.cfg.AI.FusionbrainSecret)

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - client.Do: ", err.Error())
	}

	defer response.Body.Close()

	responseBytes, errReadAll := io.ReadAll(response.Body)
	if errReadAll != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - ioutil.ReadAll: " + errReadAll.Error())
	}

	if enableLog {
		service.log.Info("FusionbrainService.CreateTask - string(response) = " + string(responseBytes))
	}

	var target *ResponseRun
	errUnmarshal := json.Unmarshal(responseBytes, &target)
	if errUnmarshal != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - json.Unmarshal: ", errUnmarshal.Error())
	}

	if target.Uuid != emptyUuid {
		db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
		if err != nil {
			service.log.Fatal("gorm.Open error: %s", err)
		}

		db.Create(&entity.Task{Uuid: target.Uuid, Status: target.Status, Promt: &promt})
	}

	return target, nil
}

func (service *FusionbrainService) CreateTaskForImage(image entity.Image, width uint, height uint, style string, negativePromptUnclip string, enableLog bool) (*ResponseRun, error) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}

	var requestData = RequestRunParams{
		Type:                 "GENERATE",
		NumImages:            1, //!
		Height:               height,
		Width:                width,
		Style:                style,
		NegativePromptUnclip: negativePromptUnclip,
		GenerateParams: struct {
			Query string "json:\"query\""
		}{*image.Promt}, //!
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	paramsPart := make(map[string][]string)
	paramsPart["Content-Disposition"] = append(paramsPart["Content-Disposition"], "form-data; name=\"params\"")
	paramsPart["Content-Type"] = append(paramsPart["Content-Type"], "application/json")

	paramsWriter, err := writer.CreatePart(paramsPart)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - writer.CreatePart: ", err.Error())
	}

	paramsPayloadBytes, err := json.Marshal(&requestData)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - json.Marshal: ", err.Error())
	}

	_, err = paramsWriter.Write(paramsPayloadBytes)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - paramsWriter.Write: ", err.Error())
	}

	err = writer.WriteField("model_id", Kandinsky3ModelId)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - writer.WriteField: ", err.Error())
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - writer.Close: ", err.Error())
	}

	request, err := http.NewRequest(http.MethodPost, BaseURL+RunAddURL, payload)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - http.NewRequest: ", err.Error())
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Add("X-Key", "Key "+service.cfg.AI.FusionbrainApi)
	request.Header.Add("X-Secret", "Secret "+service.cfg.AI.FusionbrainSecret)

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - client.Do: ", err.Error())
	}

	defer response.Body.Close()

	responseBytes, errReadAll := io.ReadAll(response.Body)
	if errReadAll != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - ioutil.ReadAll: " + errReadAll.Error())
	}

	if enableLog {
		service.log.Info("FusionbrainService.CreateTask - string(response) = " + string(responseBytes))
	}

	var target *ResponseRun
	errUnmarshal := json.Unmarshal(responseBytes, &target)
	if errUnmarshal != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - json.Unmarshal: ", errUnmarshal.Error())
	}

	if target.Uuid != emptyUuid {
		db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
		if err != nil {
			service.log.Fatal("gorm.Open error: %s", err)
		}

		db.Create(&entity.Task{Uuid: target.Uuid, Status: target.Status, Promt: image.Promt})
	}

	return target, nil
}

func (service *FusionbrainService) GetImages(task *entity.Task, enableLog bool) (*ResponseStatus, error) {
	response, errAuthNewRequest := service.AuthGetRequest(GetAddURL + task.Uuid)

	if errAuthNewRequest != nil {
		return nil, fmt.Errorf("FusionbrainService.GetImages - AuthGetRequest: " + errAuthNewRequest.Error())
	}

	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("FusionbrainService.GetImages - io.ReadAll: " + err.Error())
	}

	stringResult := string(responseBytes)
	if strings.Contains(stringResult, ErrorTaskNotFound) {
		db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
		if err != nil {
			service.log.Fatal("gorm.Open error: %s", err)
		}

		db.Delete(task)

		return nil, fmt.Errorf(gorm.ErrRecordNotFound.Error())
	}

	if enableLog {
		service.log.Info("FusionbrainService.GetImages - string(response) = " + string(responseBytes))
	}

	var target *ResponseStatus
	errUnmarshal := json.Unmarshal(responseBytes, &target)
	if errUnmarshal != nil {
		return nil, fmt.Errorf("FusionbrainService.GetImages - json.Unmarshal: ", errUnmarshal.Error())
	}

	if target.Uuid != emptyUuid {
		db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
		if err != nil {
			service.log.Fatal("gorm.Open error: %s", err)
		}
		service.log.Info("FusionbrainService.GetImages - Updates, Task.Uuid = " + target.Uuid)
		db.Model(&task).Updates(entity.Task{Status: target.Status, ErrorDescription: &target.ErrorDescription})
	}

	if target.Status == TaskStatusDone {
		db, err := gorm.Open(postgres.Open(service.cfg.PG.URL), &gorm.Config{})
		if err != nil {
			service.log.Fatal("gorm.Open error: %s", err)
		}

		for _, image := range target.Images {
			db.Create(&entity.Image{Base64: &image, ArticleID: 1, Promt: task.Promt})
		}
	}

	return target, nil
}

func (service *FusionbrainService) SaveImageToFileSystem(img entity.Image, path string) error {
	strPointerValue := *img.Base64
	unbased, err := base64.StdEncoding.DecodeString(strPointerValue)
	if err != nil {
		return fmt.Errorf("FusionbrainService.SaveImageToFileSystem - DecodeString: ", err.Error())
	}
	r := bytes.NewReader(unbased)
	imgDecode, err := jpeg.Decode(r)
	if err != nil {
		return fmt.Errorf("FusionbrainService.SaveImageToFileSystem - jpeg.Decode: ", err.Error())
	}

	jpgFilename := path + service.RandStringBytes(7) + ".jpg" //image.Slug
	file, err := os.OpenFile(jpgFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("FusionbrainService.SaveImageToFileSystem - os.OpenFile: ", err.Error())
	}

	err = png.Encode(file, imgDecode)
	if err != nil {
		return fmt.Errorf("FusionbrainService.SaveImageToFileSystem - png.Encode: ", err.Error())
	}
	fmt.Println("JPEG file", jpgFilename, "created")

	return nil
}

func (service *FusionbrainService) RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LetterBytes[rand.Intn(len(LetterBytes))]
	}
	return string(b)
}

/*
func (service *FusionbrainService) GetStylesAsync()
{
	var uri = new Uri("https://cdn.fusionbrain.ai/static/styles/api");
	var response = await _httpClient.GetAsync(uri, token);

	response.EnsureSuccessStatusCode();
	return await response.Content.ReadFromJsonAsync<IEnumerable<Style>>(cancellationToken: token) ?? ImmutableList.Create<Style>();
}
*/
