package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/amityadav9314/goinkgrid/config"
	"github.com/amityadav9314/goinkgrid/constants"
	"github.com/amityadav9314/goinkgrid/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger struct {
	zap zap.Logger
}

var appLogger Logger

func InitLogger() {
	newZapLogger()
}

func GetLogger(ctx *gin.Context) Logger {
	fillLoggingContext(ctx)
	if appLogger == nil {
		newZapLogger()
	}

	return appLogger
}

// fillLoggingContext This function is used to fill the logging context with the api name and file name.
func fillLoggingContext(ctx *gin.Context) {
	pc, _, _, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	functionName := fn.Name()
	fileName, _ := fn.FileLine(pc)
	ctx.Set(ApiName, functionName)
	ctx.Set(FileName, utils.GetLastPart(fileName))
}

func newZapLogger() {
	env := config.LoadAppEnv()
	logFile := os.Getenv(constants.VarLogFile)
	appLogger = AppLogger{
		zap: *configureZapLoggerV2(env, logFile),
	}
}

func configureZapLogger(env string, logFile string) *zap.Logger {
	var loggerConfig zap.Config
	var logLevel = zap.InfoLevel

	loggerConfig.EncoderConfig.MessageKey = "MSG"
	loggerConfig.EncoderConfig.LevelKey = "LVL"
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	loggerConfig.EncoderConfig.TimeKey = zapcore.OmitKey // Remove timestamp
	loggerConfig.EncoderConfig.EncodeTime = nil
	loggerConfig.EncoderConfig.CallerKey = zapcore.OmitKey // Remove caller
	loggerConfig.EncoderConfig.EncodeCaller = nil
	loggerConfig.EncoderConfig.StacktraceKey = ""
	loggerConfig.InitialFields = map[string]interface{}{"APP": constants.Application}

	var encoder zapcore.Encoder

	if utils.StringInSlice(env, []string{constants.EnvDev, constants.EnvPp}) {
		loggerConfig = zap.NewDevelopmentConfig()
		logLevel = zap.DebugLevel
		encoder = zapcore.NewConsoleEncoder(loggerConfig.EncoderConfig)
	} else {
		loggerConfig = zap.NewProductionConfig()
		if config.Config.GetBool("logging.debug_enabled") {
			logLevel = zap.DebugLevel
		}
		encoder = zapcore.NewJSONEncoder(loggerConfig.EncoderConfig)
	}

	// writer
	var writer zapcore.WriteSyncer
	if len(logFile) != 0 {
		loggerConfig.OutputPaths = []string{LOG_FILE.String()}
		writer = zapcore.AddSync(LogRotator(logFile))
	} else {
		loggerConfig.OutputPaths = []string{STANDARD_OUTPUT.String()}
		writer = zapcore.Lock(os.Stdout)
	}

	zapLogger, err := loggerConfig.Build()
	if err != nil {
		panic("error building zap logger")
	}

	return zapLogger.WithOptions(
		zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewCore(encoder, writer, logLevel)
		}),
	)
}

func configureZapLoggerV2(env string, logFileStr string) *zap.Logger {
	var core zapcore.Core
	var writerSync zapcore.WriteSyncer
	if len(logFileStr) != 0 {
		fileWriteSync := zapcore.AddSync(LogRotator(logFileStr))
		writerSync = fileWriteSync
	} else {
		// Create a console write sync
		consoleWriteSync := zapcore.AddSync(os.Stdout)
		writerSync = consoleWriteSync
	}

	// If you want to print logs in both console and file, then enable belwo two lines and put multiWriteSync below
	// Create a write sync for stdout
	//consoleWriteSync := zapcore.AddSync(os.Stdout)
	// Create a multi-write sync to write to both file and console
	//multiWriteSync := zapcore.NewMultiWriteSyncer(fileWriteSync, consoleWriteSync)

	// Create an encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timeStamp"
	encoderConfig.StacktraceKey = "stackTrace"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Set the logging level to capture all levels
	levelEnabler := zapcore.DebugLevel
	if utils.StringInSlice(env, []string{constants.EnvProd}) {
		// but if it is live, then only info
		levelEnabler = zapcore.InfoLevel
	}

	// Create a core for logging to the console
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writerSync,
		levelEnabler,
	)

	logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync() // flushes buffer, if any
	return logger
}

func (appLogger AppLogger) Debug(ctx context.Context, extraInfo string, request any, response any, timeTaken time.Duration, fields ...LoggingField) {
	fields = append(fields, FieldAny(Request, request))
	fields = append(fields, FieldAny(Response, response))
	fields = append(fields, FieldAny(TimeTaken, timeTaken))
	appLogger.log(ctx, LvlDebug, extraInfo, fields...)
}

func (appLogger AppLogger) Info(ctx context.Context, extraInfo string, request any, response any, timeTaken time.Duration, fields ...LoggingField) {
	fields = append(fields, FieldAny(Request, request))
	fields = append(fields, FieldAny(Response, response))
	fields = append(fields, FieldAny(TimeTaken, timeTaken))
	appLogger.log(ctx, LvlInfo, extraInfo, fields...)
}

func (appLogger AppLogger) Error(ctx context.Context, extraInfo string, request any, response any, timeTaken time.Duration, e error, fields ...LoggingField) {
	fields = append(fields, FieldAny(Request, request))
	fields = append(fields, FieldAny(Response, response))
	fields = append(fields, FieldAny(TimeTaken, timeTaken))
	fields = append(fields, FieldAny(Exception, e))
	appLogger.log(ctx, LvlError, extraInfo, fields...)
}

func (appLogger AppLogger) Flush() {
	err := appLogger.zap.Sync()
	if err != nil {
		fmt.Printf("failing while syncing zap logs %v", err)
	}
}

func (appLogger AppLogger) log(ctx context.Context, logLevel LogLevel, msg string, fields ...LoggingField) {
	zapFields := getZapFields(ctx, fields)

	switch logLevel {
	case LvlDebug:
		appLogger.zap.Debug(msg, zapFields...)
	case LvlInfo:
		appLogger.zap.Info(msg, zapFields...)
	case LvlError:
		appLogger.zap.Error(msg, zapFields...)
	default:
		appLogger.zap.Error(msg, zapFields...)
	}
}

func getZapFields(ctx context.Context, fields []LoggingField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		switch field.kind {
		case fieldAny:
			zapFields[i] = zap.Any(field.key, field.value)
		case fieldString:
			zapFields[i] = zap.String(field.key, field.value.(string))
		case fieldError:
			if field.value != nil {
				zapFields[i] = zap.Error(field.value.(error))
			} else {
				zapFields[i] = zap.Error(nil)
			}
		default:
			zapFields[i] = zap.String(field.key, "unknown field for logger")
		}
	}
	if fileNameV, ok := ctx.Value(FileName).(string); ok {
		zapFields = append(zapFields, zap.String(FileName, fileNameV))
	}
	if apiNameValue, ok := ctx.Value(ApiName).(string); ok {
		zapFields = append(zapFields, zap.String(ApiName, apiNameValue))
	}
	zapFields = append(zapFields, zap.String(Env, config.GetEnv()))
	//zapFields = append(zapFields, zap.String(TimeStamp, time.Now().Format(time.RFC3339)))
	return zapFields
}
