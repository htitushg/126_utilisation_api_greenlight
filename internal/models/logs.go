package models

import (
	"125_isbn_new/internal/assert"
	"125_isbn_new/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

var Logger *slog.Logger
var logs *os.File
var wg sync.WaitGroup

// closeLog checks if the log file is open and if yes, close it.
func closeLog() {
	if logs != nil {
		err := logs.Close()
		if err != nil {
			log.Println(GetCurrentFuncName(), slog.Any("output", err))
		}
	}
}

// LogInit is meant to be run as a goroutine to create a new log file every day
// appending the file's creation timestamp in its name.
func LogInit() {
	duration := utils.SetDailyTimer(0)
	var jsonHandler *slog.JSONHandler
	var err error
	var filename string
	defer closeLog()
	for {
		filename = assert.Path + "logs/logs_" + time.Now().Format(time.DateOnly) + ".log"
		closeLog()
		logs, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(GetCurrentFuncName(), slog.Any("output", err))
		}
		jsonHandler = slog.NewJSONHandler(logs, nil)
		Logger = slog.New(jsonHandler)
		time.Sleep(duration)
		duration = time.Hour * 24
	}
}

// fetchLogInfo retrieves all Log from `file` and stores it in *log.
func (log *Logs) fetchLogInfo(file string) {
	fmt.Println(GetCurrentFuncName())
	defer wg.Done()
	filename := "logs/" + file
	data, _ := os.ReadFile(filename)
	if len(data) == 0 {
		return
	}
	lines := bytes.Split(data, []byte("\n"))
	var singleLog Log
	for _, line := range lines {
		err := json.Unmarshal(line, &singleLog)
		if err != nil {
			return
		}
		*log = append(*log, singleLog)
	}
}

/* func printFileNames(files []os.DirEntry) []string {
	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result
} */

// RetrieveLogs fetches all Log from all files *.log in /logs directory
// and returns a Logs array.
func RetrieveLogs() (logArray Logs) {
	logFiles, err := os.ReadDir(assert.Path + "logs/.")
	if err != nil {
		Logger.Error(GetCurrentFuncName(), slog.Any("output", err))
	} else {
		reg := regexp.MustCompile(`^[a-zA-Z0-9-_]+\.log$`)
		for _, file := range logFiles {
			if reg.MatchString(file.Name()) {
				wg.Add(1)
				go logArray.fetchLogInfo(file.Name())
			}
		}
	}
	wg.Wait()
	logArray.sortLogs()
	return logArray
}

// sortLogs sort all Log from the newest to the oldest.
func (log *Logs) sortLogs() {
	sort.Slice(*log, func(i, j int) bool {
		return (*log)[i].Time.After((*log)[j].Time)
	})
}

// FetchLevelLogs filters Log returning only Log matching the given `level`.
func FetchAttrLogs(attr string, value string) Logs {
	attr = strings.ToLower(attr)
	logs := RetrieveLogs()
	var result Logs
	switch attr {
	case "level":
		switch strings.ToUpper(value) {
		case "INFO", "WARN", "ERROR":
			for _, singleLog := range logs {
				if singleLog.Level == strings.ToUpper(value) {
					result = append(result, singleLog)
				}
			}
			//break
		default:
			return nil
		}
	case "user", "Pseudo":
		for _, singleLog := range logs {
			if strings.EqualFold(singleLog.User.Pseudo, value) {
				result = append(result, singleLog)
			}
		}
		//break
	default:
		return nil
	}
	result.sortLogs()
	return result
}
