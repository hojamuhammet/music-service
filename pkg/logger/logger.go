package logger

import (
	"io"
	"log/slog"
	"os"
)

type Loggers struct {
	InfoLogger  *slog.Logger
	DebugLogger *slog.Logger
	ErrorLogger *slog.Logger
}

func SetupLogger(env string) (*Loggers, error) {
	var infoHandler slog.Handler
	var debugHandler slog.Handler
	var errorHandler slog.Handler

	if env == "test" {
		infoHandler = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
		debugHandler = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})
		errorHandler = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	} else {
		err := os.MkdirAll("logs", 0755)
		if err != nil {
			return nil, err
		}

		infoFile, err := os.OpenFile("logs/Info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			if os.IsNotExist(err) {
				infoFile, err = os.Create("logs/Info.log")
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		debugFile, err := os.OpenFile("logs/Debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			infoFile.Close()
			if os.IsNotExist(err) {
				debugFile, err = os.Create("logs/Debug.log")
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		errorFile, err := os.OpenFile("logs/Error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			infoFile.Close()
			debugFile.Close()
			if os.IsNotExist(err) {
				errorFile, err = os.Create("logs/Error.log")
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		infoHandler = slog.NewTextHandler(infoFile, &slog.HandlerOptions{Level: slog.LevelInfo})
		debugHandler = slog.NewTextHandler(debugFile, &slog.HandlerOptions{Level: slog.LevelDebug})
		errorHandler = slog.NewTextHandler(errorFile, &slog.HandlerOptions{Level: slog.LevelError})
	}

	infoLogger := slog.New(infoHandler)
	debugLogger := slog.New(debugHandler)
	errorLogger := slog.New(errorHandler)

	return &Loggers{
		InfoLogger:  infoLogger,
		DebugLogger: debugLogger,
		ErrorLogger: errorLogger,
	}, nil
}
