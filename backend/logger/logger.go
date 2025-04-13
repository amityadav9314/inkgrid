package logger

import (
	"context"
	"time"
)

type LogLevel uint8

const (
	LvlDebug LogLevel = iota
	LvlInfo
	LvlError
)

type fieldKind uint8

const (
	fieldString fieldKind = iota
	fieldError
	fieldAny
)

type LoggingField struct {
	key   string
	kind  fieldKind
	value any
}

type Logger interface {
	Debug(ctx context.Context, extraInfo string, request any, response any, timeTaken time.Duration, fields ...LoggingField)
	Info(ctx context.Context, extraInfo string, request any, response any, timeTaken time.Duration, fields ...LoggingField)
	Error(ctx context.Context, extraInfo string, request any, response any, timeTaken time.Duration, e error, fields ...LoggingField)
	// Flush syncs any pending logs
	Flush()
}

func FieldAny(key string, val any) LoggingField {
	return LoggingField{
		key:   key,
		kind:  fieldAny,
		value: val,
	}
}

func FieldString(key string, val string) LoggingField {
	return LoggingField{
		key:   key,
		kind:  fieldString,
		value: val,
	}
}

func FieldError(val error) LoggingField {
	return LoggingField{
		kind:  fieldError,
		value: val,
	}
}
