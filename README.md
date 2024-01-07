
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

## Roadmap
v0.2.0 Keyword Research Tool

## Changelog
v0.1.0 - LinGoose. Huggingface Service
v0.0.1 - Start in docker & local