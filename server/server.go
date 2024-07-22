package server

import (
	gitRepos "gitbeam.repo.manager/contract/repos"
	"gitbeam.repo.manager/core"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func ExecGRPCServer(address string, core *core.GitBeamService, logger *logrus.Logger) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Recovered from err: %v", err)
		}
	}()
	api := NewApiService(core, logger)

	server := grpc.NewServer()
	gitRepos.RegisterGitBeamRepositoryServiceServer(server, api)
	reflection.Register(server)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Fatalf("failed to listen with the following errors: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}
