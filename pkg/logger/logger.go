package logger

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
)

var Log *slog.Logger

func InitLogger() *os.File {
	var logFilePath string
	logDirPath := getFilePath()
	if err := createLogDir(logDirPath); err != nil {
		fmt.Println(err)
		logFilePath = "./thanos.log"
	} else {
		fmt.Println("here")
		logFilePath = logDirPath + "thanos.log"
	}
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	writer := io.MultiWriter(os.Stdout, logFile)
	logHandler := slog.NewJSONHandler(writer, nil)
	Log = slog.New(logHandler)
	return logFile
}

// func GetInstance() *slog.Logger {
// 	if Logger == nil {
// 		mutex.Do(InitLogger)
// 		return Logger
// 	} else {
// 		return Logger
// 	}
// }

func createLogDir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func getFilePath() string {
	var path string
	switch runtime.GOOS {
	case "windows":
		path = os.Getenv("LocalAppData") + "\\thanos\\"
	default:
		path = "/var/log/thanos/"
	}
	return path
}
