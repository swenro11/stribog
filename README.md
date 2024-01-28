
# Stribog
## Start project in docker
```
docker-compose up -d
### if u need to rebuild all docker containers
docker-compose up -d --force-recreate
```
Output u can see in Docker Dashboard

## Start project local
But with enviroment in docker containers. 
First start
```
docker-compose up -d
go mod init github.com/swenro11/stribog
#update and install
go get -u ./... && go mod tidy 
#just install
go mod tidy 
# apply migrations
go run -tags migrate ./cmd/app 
# start main task
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

### Online AI
[Go framework for building awesome LLM apps](https://lingoose.io/), [Github](https://github.com/henomis/lingoose)
Use models from huggingface, [Open LLM Leaderboard](https://huggingface.co/spaces/HuggingFaceH4/open_llm_leaderboard).
For free account use models < 10GB. 

### LocalAI
Start docker container with default model.  
```
git clone https://github.com/mudler/LocalAI.git 
cd ~/projects/LocalAI$ 
# configure default LLM (chat-gpt-3.5-turbo)
docker compose up -d
```

## Roadmap
v0.3.2 - ArticleService. Hugo
v0.3.1 - ArticleService. Images
v0.3.0 - ImageService
v0.2.3 - ArticleService. CheckUnique & CheckNotAI
v0.2.2 - ArticleService. SeoOptimization
v0.2.1 - ArticleService. Mock Statuses

## Changelog
v0.2.0 - ArticleService - Status New. GORM with models - Article & Image. Add generate.go instead of migrate.go
v0.1.2 - LocalAIService - LinGoose.
v0.1.1 - Cohere Service - LinGoose. Huggingface Service - hupe1980/go-huggingface
v0.1.0 - Huggingface Service - LinGoose. 
v0.0.1 - Start in docker & local