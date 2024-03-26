package service

import (
	"context"
	"fmt"
	"os"

	"github.com/henomis/lingoose/llm/huggingface"
	huggingfaceHupe1980 "github.com/hupe1980/go-huggingface"
	"github.com/swenro11/stribog/config"
	log "github.com/swenro11/stribog/pkg/logger"
)

const (
	HuggingfaceTokenParam = "HUGGING_FACE_HUB_TOKEN"
	//https://huggingface.co/{model_name}
	//HuggingfaceModelgpt2 = "gpt2"

	/* microsoft/phi-2
	The repository for microsoft/phi-2 contains custom code which must be executed to correctly load the model.
	You can inspect the repository content at https://hf.co/microsoft/phi-2.\nPlease pass the argument `trust_remote_code=True` to allow custom code to be run.
	*/

	/*
		The model TomGrc/FusionNet_7Bx2_MoE_14B is too large to be loaded automatically (25GB > 10GB).
		Please use Spaces (https://huggingface.co/spaces) or Inference Endpoints (https://huggingface.co/inference-endpoints)."}
	*/
	HuggingfaceModel7Bx2FusionNet = "TomGrc/FusionNet_7Bx2_MoE_14B"
	/*
		The model Weyaxi/Nous-Hermes-2-SUS-Chat-2x34B is too large to be loaded automatically (121GB > 10GB).
		Please use Spaces (https://huggingface.co/spaces) or Inference Endpoints (https://huggingface.co/inference-endpoints).

		For others model - same error about "is too large to be loaded automatically"
	*/
	HuggingfaceModel2x34bHermes2 = "Weyaxi/Nous-Hermes-2-SUS-Chat-2x34B"
	HuggingfaceModel67bDeepseek  = "deepseek-ai/deepseek-llm-67b-base"
	HuggingfaceModel70bCOKAL     = "DopeorNope/COKAL-v1-70B"
)

type HuggingfaceService struct {
	cfg *config.Config
	log *log.Logger
}

func NewHuggingfaceService(cfg *config.Config, l *log.Logger) *HuggingfaceService {
	return &HuggingfaceService{
		cfg: cfg,
		log: l,
	}
}

func (service *HuggingfaceService) CheckGetenv(enableLog bool) string {
	huggingfaceToken := os.Getenv(HuggingfaceTokenParam)
	if enableLog {
		service.log.Info("HuggingfaceService - checkGetenv -  os.Getenv = " + huggingfaceToken)
	}
	if len(huggingfaceToken) > 0 {
		return huggingfaceToken
	}

	os.Setenv(HuggingfaceTokenParam, service.cfg.AI.HuggingfaceToken)
	huggingfaceToken = os.Getenv(HuggingfaceTokenParam)
	if enableLog {
		service.log.Info("HuggingfaceService - checkGetenv -  second os.Getenv = " + huggingfaceToken)
	}

	return huggingfaceToken
}

func (service *HuggingfaceService) TextGenerationLingoose(model string, prompt string) (string, error) {

	llm := huggingface.New(model, 0.1, false).WithMode(huggingface.ModeTextGeneration)

	result, errCompletion := llm.Completion(context.Background(), prompt)
	if errCompletion != nil {
		return result, fmt.Errorf("HuggingfaceService - llm.Completion: " + errCompletion.Error())
	}

	return result, nil
}

func (service *HuggingfaceService) FusionNetTextGenLingoose(prompt string) (string, error) {
	huggingfaceToken := service.CheckGetenv(false)
	if len(huggingfaceToken) > 0 {
		return service.TextGenerationLingoose(HuggingfaceModel7Bx2FusionNet, prompt)
	} else {
		return "", fmt.Errorf("HuggingfaceService.FusionNetTextGenLingoose - Not found HuggingfaceToken")
	}
}

func (service *HuggingfaceService) TextGenerationHupe1980(model string, prompt string) (string, error) {

	ic := huggingfaceHupe1980.NewInferenceClient(service.cfg.AI.HuggingfaceToken)

	result, errTextGeneration := ic.TextGeneration(context.Background(), &huggingfaceHupe1980.TextGenerationRequest{
		Inputs: prompt,
		Model:  model,
	})
	if errTextGeneration != nil {
		return "", fmt.Errorf("HuggingfaceService.TextGenerationHupe1980 - " + errTextGeneration.Error())
	}

	return result[0].GeneratedText, nil
}

func (service *HuggingfaceService) FusionNetTextGenHupe1980(prompt string) (string, error) {
	return service.TextGenerationHupe1980(HuggingfaceModel7Bx2FusionNet, prompt)
}

func (service *HuggingfaceService) HermesTextGenHupe1980(prompt string) (string, error) {

	/*
		_HuggingfaceMode3blMicrosoftPhi2 - error
		The repository for microsoft/phi-2 contains custom code which must be executed to correctly load the model.
		You can inspect the repository content at https://hf.co/microsoft/phi-2.\nPlease pass the argument `trust_remote_code=True` to allow custom code to be run.
	*/

	return service.TextGenerationHupe1980(HuggingfaceModel2x34bHermes2, prompt)
}
