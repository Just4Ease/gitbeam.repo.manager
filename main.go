package main

import (
	"fmt"
	"gitbeam.baselib/store"
	"gitbeam.repo.manager/config"
	"gitbeam.repo.manager/core"
	"gitbeam.repo.manager/repository"
	"gitbeam.repo.manager/repository/sqlite"
	"gitbeam.repo.manager/server"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	var eventStore store.EventStore
	var dataStore repository.DataStore
	var err error

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	secrets := config.GetSecrets()

	//Using SQLite as the mini persistent storage.
	//( in a real world system, this would be any production level or vendor managed db )
	if dataStore, err = sqlite.NewSqliteRepo(secrets.DatabaseName); err != nil {
		logger.WithError(err).Fatal("failed to initialize sqlite database repository for dataStore.")
	}

	// A channel based pub/sub messaging system.
	//( in a real world system, this would be apache-pulsar, kafka, nats.io or rabbitmq )
	eventStore = store.NewEventStore(logger)

	// If the dependencies were more than 3, I would use a variadic function to inject them.
	//Clarity is better here for this exercise.
	coreService := core.NewGitBeamService(logger, eventStore, dataStore, nil)

	address := fmt.Sprintf("0.0.0.0:%s", secrets.Port)
	logger.Printf("[*] %s listening on address: %s", config.ServiceName, address)
	server.ExecGRPCServer(address, coreService, logger)
}
