package utils

import (
	"time"

	"github.com/google/uuid"
)

type LoggerFactory struct {
	LogList    map[string]*Logger
	LogAddress LogStore
}

type LogStore interface {
	SendLog(input []string) bool
}

func (lf *LoggerFactory) CreateLog() *Logger {
	genUuid := uuid.New().String()
	newLog := &Logger{
		Factory: lf,
		Id:      genUuid,
	}
	lf.LogList[genUuid] = newLog
	return newLog

}

func (lf *LoggerFactory) WriteLog(loggerId string) {
	sentLog := lf.LogList[loggerId].ToExternalFormat()
	lf.LogAddress.SendLog(sentLog)
}

func (lf *LoggerFactory) DeleteLog(loggerId string) {
	delete(lf.LogList, loggerId)
}

type Logger struct {
	Factory      *LoggerFactory
	Id           string
	StartAt      time.Time
	FinishAt     time.Time
	Remarks      string
	FunctionName string
}

// Implement with channel and go func so it would not
// add more run time in a same thread

func (l *Logger) ToExternalFormat() []string {
	startAt, _ := l.StartAt.MarshalText()
	finishtAt, _ := l.StartAt.MarshalText()
	externalFormat := []string{
		l.Id,
		string(startAt),
		string(finishtAt),
		l.Remarks,
		l.FunctionName,
	}
	return externalFormat
}

func (l *Logger) SetStartTime() *Logger {
	l.StartAt = time.Now()
	return l
}

func (l *Logger) SetFinishTime() *Logger {
	l.FinishAt = time.Now()
	return l
}

func (l *Logger) SetRemarks(remarks string) *Logger {
	l.Remarks = remarks
	return l
}

func (l *Logger) SetFunctionName(functionName string) *Logger {
	l.FunctionName = functionName
	return l
}

func (l *Logger) WriteLog() {
	l.Factory.WriteLog(l.Id)
}

func (l *Logger) DeleteLog() {
	l.Factory.DeleteLog(l.Id)
}

func (l *Logger) WriteAndDeleteLog() {
	l.WriteLog()
	l.DeleteLog()
}
