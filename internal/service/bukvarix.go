package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
)

type BukvarixService struct {
	cfg *config.Config
	log *log.Logger
}

const (
	KeywordsUrl  = "http://api.bukvarix.com/v1/keywords/"
	MkeywordsUrl = "http://api.bukvarix.com/v1/mkeywords/"
	Separator    = "\r\n"
)

// docs - https://www.bukvarix.com/api_keywords.html
func NewBukvarixService(cfg *config.Config, l *log.Logger) *BukvarixService {
	return &BukvarixService{
		cfg: cfg,
		log: l,
	}
}

func (service *BukvarixService) Keywords(keyword string) ([]string, error) {
	url := KeywordsUrl + "?api_key=" + service.cfg.PARAM.BukvarixApiKey + "&q=" + url.QueryEscape(keyword)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("BukvarixService.Keywords - http.NewRequest: " + err.Error())
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BukvarixService.Keywords - client.Do: " + err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("BukvarixService.Keywords - io.ReadAll: " + err.Error())
	}

	return strings.Split(string(body), Separator), nil
}
