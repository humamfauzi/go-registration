package main

import (
	"context"

	"github.com/humamfauzi/go-registration/exconn"
	"github.com/humamfauzi/go-registration/utils"
	"github.com/jinzhu/gorm"
)

var (
	db            *gorm.DB
	logStore      exconn.CassandraLog
	documentDb    context.Context
	errorMap      map[string]string
	loggerFactory utils.LoggerFactory
)

func InstantiateExternalConnection() {
	db = exconn.ConnectToMySQL()
	logStore = exconn.InstantiateCassandraLog()
	documentDb = exconn.ConnectToMongo()
	errorMap = utils.InitError("./error.json")
	loggerFactory = InstantiateLoggerFactory(&logStore)
}

func InstantiateLoggerFactory(logStore utils.LogStore) utils.LoggerFactory {
	return utils.LoggerFactory{
		LogList:    make(map[string]*utils.Logger),
		LogAddress: logStore,
	}
}
