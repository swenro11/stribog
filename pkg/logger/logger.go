package logger

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// Logger -.
type Logger struct {
	logger      *zerolog.Logger
	tgtokenapi  string
	tgchatid    string
	mongoClient *mongo.Client
}

var _ Interface = (*Logger)(nil)

// New -.
func New(ctx context.Context, level string, tgtoken string, chatid string, mongouri string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	mongoClientOptions := options.Client().ApplyURI(mongouri)
	client, err := mongo.Connect(ctx, mongoClientOptions)
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger:      &logger,
		tgtokenapi:  tgtoken,
		tgchatid:    chatid,
		mongoClient: client,
	}
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}

	l.msg("error", message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info().Msg(message)
	} else {
		l.logger.Info().Msgf(message, args...)
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}

func (l *Logger) Telegram(message string) {
	bot, err := tgbotapi.NewBotAPI(l.tgtokenapi)
	if err != nil {
		panic(err)
	}
	intChatId, err := strconv.ParseInt(l.tgchatid, 10, 64)
	if err != nil {
		panic(err)
	}
	msg := tgbotapi.NewMessage(intChatId, message)

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func (l *Logger) Mongo(ctx context.Context, TsUid int64, message string) {
	errPing := l.mongoClient.Ping(ctx, nil)
	if errPing != nil {
		panic(errPing)
	}

	mongoDatabase := l.mongoClient.Database("gotasks")
	logsCollection := mongoDatabase.Collection("logs")
	mongoLogsResult, errInsert := logsCollection.InsertOne(ctx, bson.D{
		{Key: "TS", Value: time.Now().UnixMicro()},
		{Key: "TsUid", Value: TsUid},
		{Key: "message", Value: message},
	})
	if errInsert != nil {
		panic(errInsert)
	}
	l.Debug("Mongo InsertedID %v \n", mongoLogsResult.InsertedID)
}
