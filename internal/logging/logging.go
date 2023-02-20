package logging

import (
	"github.com/upper-institute/flipbook/internal/helpers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevel_flag = "log.level"
	logEnv_flag   = "log.env"
)

var (
	Logger    *zap.Logger
	zapConfig zap.Config
)

func BindOptions(binder helpers.FlagBinder) {

	binder.BindString(logLevel_flag, "debug", "Logging level of stdout (debug, info or error)")
	binder.BindString(logEnv_flag, "prod", "Logging env (prod or dev)")

}

func Load(getter helpers.FlagGetter) {

	zapConfig = zap.NewProductionConfig()

	if getter.GetString(logEnv_flag) == "dev" {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	switch getter.GetString(logLevel_flag) {

	case "error":
		zapConfig.Level.SetLevel(zap.ErrorLevel)

	case "info":
		zapConfig.Level.SetLevel(zap.InfoLevel)

	default:
		zapConfig.Level.SetLevel(zap.DebugLevel)

	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger.Named("flipbook")
}

func Flush() {

	Logger.Sync()

}
