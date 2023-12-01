
# Stribog
## config
Create Posgres DB & user, stribog by example.  
Start RabiitMQ, 


## Start project
First start
```
go mod init github.com/swenro11/stribog
#update and install
go get -u ./... && go mod tidy 
#just install
go mod tidy 
# apply migrations
go run -tags migrate ./cmd/app 
# start main task
```

## Project Info
Application Based on [Go-clean-template](https://github.com/evrone/go-clean-template)  
Sheduller Based on [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)  
Add to project - go get github.com/jasonlvhit/gocron  

Telegram API [Go Telegram Bot API](https://go-telegram-bot-api.dev/)  
Add to project -  go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5  

MongoDB, manual [Ubuntu](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-ubuntu/)

