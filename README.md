
# Stribog

## Start docker
```bash
docker-compose up -d
# if u need to rebuild all docker containers
docker-compose up -d --force-recreate
```
Output u can see in Docker Dashboard

## Init, update and install
```bash
docker-compose up -d
go mod init github.com/swenro11/stribog
go get -u ./... && go mod tidy 
# just install
go mod tidy
```

## Start project 
Localy, with enviroment in docker containers.  
```bash
go run ./cmd/app
```

## Project Info
Application Based on [Go-clean-template](https://github.com/evrone/go-clean-template)  
Sheduller Based on [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)  
Add to project - go get github.com/jasonlvhit/gocron  

Telegram API [Go Telegram Bot API](https://go-telegram-bot-api.dev/)  
Add to project - go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5 

## ORM
[GORM](https://gorm.io/), top [from list](https://github.com/d-tsuji/awesome-go-orms) 

## AI/LLM
[Go framework for building awesome LLM apps](https://lingoose.io/), [Github](https://github.com/henomis/lingoose)  

### Online AI
- Huggingface, [Open LLM Leaderboard](https://huggingface.co/spaces/HuggingFaceH4/open_llm_leaderboard). **Free account allows you to use models < 10GB.**  
- [Cohere](https://cohere.com/) with LinGoose. 

### LocalAI
Start docker container with default model.  
```bash
git clone https://github.com/mudler/LocalAI.git 
cd ~/projects/LocalAI$ 
# configure default LLM (chat-gpt-3.5-turbo)
docker compose up -d
```

## Roadmap
**v1.0.0-BETA** - DeployToProdHugo.  
v0.9.0 - DeployToTestHugo. WriterService.ConvertToMd. 
v0.8.0 - CheckNotAI.  
v0.7.0 - CheckUnique.  
v0.6.0 - Generating & ReadyWithImages.  
v0.5.0 - KeywordsService. Get tags/keywords list by short description. Save to DB.  
V0.4.0 - HugoService.New  
V0.3.1 - Rename ArticleService to WriterService. Add WriterService.CreateArticleWithImages  

### Maybe
SaveCDN, DeleteCDN based on https://github.com/cloudflare/cloudflare-go  
DeleteImageFromFileSystem  

## Changelog
v0.3.0 - FusionbrainService.SaveImageToFileSystem
v0.2.5 - mock TasksService.Flow, refactoring, update Readme 
v0.2.4 - FusionbrainService.GetImages & Image.Base64 (save to DB)   
v0.2.3 - FusionbrainService.CreateTask & Task Entity  
v0.2.2 - FusionbrainService.CreateTaskString, Update dependencies  
v0.2.1 - KeywordService  
v0.2.0 - ArticleService - Status New. GORM with models - Article & Image. Add generate.go instead of migrate.go  
v0.1.2 - LocalAIService - LinGoose.  
v0.1.1 - Cohere Service - LinGoose. Huggingface Service - hupe1980/go-huggingface  
v0.1.0 - Huggingface Service - LinGoose.  
v0.0.1 - Start in docker & local  