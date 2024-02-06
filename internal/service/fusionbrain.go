package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
	"github.com/tidwall/gjson"
)

const (
	_BaseURL           = "https://api-key.fusionbrain.ai/key/api/v1/"
	_StylesURL         = "https://cdn.fusionbrain.ai/static/styles/api"
	_ModelsAddURL      = "models"
	_RunAddURL         = "text2image/run"
	_Kandinsky3ModelId = "4"
)

type FusionbrainService struct {
	cfg *config.Config
	log *log.Logger
}

type ResponseRun struct {
	uuid   string
	status string
}

type RequestRunModel struct {
	ModelID uint `json:"model_id"`
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

func NewFusionbrainService(cfg *config.Config, l *log.Logger) *FusionbrainService {
	return &FusionbrainService{
		cfg: cfg,
		log: l,
	}
}

func (service *FusionbrainService) AuthGetRequest(addURL string) (*http.Response, error) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}

	req, err := http.NewRequest(http.MethodGet, _BaseURL+addURL, nil)
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

// https://freshman.tech/snippets/go/multipart-upload-google-drive/
// data must be MultipartFormDataContent (sorry for my C# commentary in Go)
func (service *FusionbrainService) AuthPostRequest(addURL string, data []byte) (*http.Response, error) {
	client := http.Client{Timeout: time.Duration(3) * time.Second}

	req, err := http.NewRequest(http.MethodPost, _BaseURL+addURL, bytes.NewBuffer(data))
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
	response, errAuthNewRequest := service.AuthGetRequest(_ModelsAddURL)

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
	response, errAuthNewRequest := service.AuthGetRequest(_ModelsAddURL)

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

// {\"status\":\"INITIAL\",\"uuid\":\"0a5b8c21-4e59-4ab8-a592-093fc5b0cc77\"}
func (service *FusionbrainService) CreateTaskString(promt string, quantity uint, width uint, height uint, style string, negativePromptUnclip string) (string, error) {
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
		return "uuid.Nil", errors.WithStack(err)
	}

	paramsPayloadBytes, err := json.Marshal(&requestData)
	if err != nil {
		return "uuid.Nil", errors.WithStack(err)
	}

	_, err = paramsWriter.Write(paramsPayloadBytes)
	if err != nil {
		return "uuid.Nil", errors.WithStack(err)
	}

	err = writer.WriteField("model_id", "4")
	if err != nil {
		return "uuid.Nil", errors.WithStack(err)
	}

	err = writer.Close()
	if err != nil {
		return "uuid.Nil", errors.WithStack(err)
	}

	request, err := http.NewRequest(http.MethodPost, _BaseURL+_RunAddURL, payload)
	if err != nil {
		return "uuid.Nil", errors.WithStack(err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Add("X-Key", "Key "+service.cfg.AI.FusionbrainApi)
	request.Header.Add("X-Secret", "Secret "+service.cfg.AI.FusionbrainSecret)

	resp, err := client.Do(request)
	if err != nil {
		return "", errors.WithStack(err)
	}

	defer resp.Body.Close()

	respBytes, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		return "", fmt.Errorf("FusionbrainService.CreateTaskString - ioutil.ReadAll: " + errReadAll.Error())
	}

	return string(respBytes), nil
}

// always "00000000-0000-0000-0000-000000000000" in result, need fix this
func (service *FusionbrainService) CreateTask(promt string, quantity uint, width uint, height uint, style string, negativePromptUnclip string) (uuid.UUID, error) {
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
		return uuid.Nil, errors.WithStack(err)
	}

	paramsPayloadBytes, err := json.Marshal(&requestData)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	_, err = paramsWriter.Write(paramsPayloadBytes)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	err = writer.WriteField("model_id", _Kandinsky3ModelId)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	err = writer.Close()
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	request, err := http.NewRequest(http.MethodPost, _BaseURL+_RunAddURL, payload)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Add("X-Key", "Key "+service.cfg.AI.FusionbrainApi)
	request.Header.Add("X-Secret", "Secret "+service.cfg.AI.FusionbrainSecret)

	resp, err := client.Do(request)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return uuid.Nil, errors.WithStack(err)
	}

	imageId, err := uuid.Parse(gjson.GetBytes(respBytes, "uuid").String())

	return imageId, nil
}

/*
func (service *FusionbrainService) CreateTask(promt string, quantity uint, width uint, height uint) (*ResponseRun, error) {
	var request = RequestRun{
		ModelID: _Kandinsky3ModelId,
		Params: RequestRunParams{
			Type:      "GENERATE",
			NumImages: quantity,
			Height:    height,
			Width:     width,
			GenerateParams: struct {
				Query string "json:\"query\""
			}{promt},
		}}

	requestBody, errMarshal := json.Marshal(request)

	if errMarshal != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - json.Marshal: ", errMarshal.Error())
	}

	response, errAuthPostRequest := service.AuthPostRequest(_RunAddURL, requestBody)

	if errAuthPostRequest != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - AuthPostRequest: ", errAuthPostRequest.Error())
	}

	defer response.Body.Close()

	var target *ResponseRun
	errDecode := json.NewDecoder(response.Body).Decode(target)
	if errDecode != nil {
		return nil, fmt.Errorf("FusionbrainService.CreateTask - Decode: ", errDecode.Error())
	}

	return target, nil
}
*/

/*
func (service *FusionbrainService) GetStylesAsync()
{
	var uri = new Uri("https://cdn.fusionbrain.ai/static/styles/api");
	var response = await _httpClient.GetAsync(uri, token);

	response.EnsureSuccessStatusCode();
	return await response.Content.ReadFromJsonAsync<IEnumerable<Style>>(cancellationToken: token) ?? ImmutableList.Create<Style>();
}
*/
