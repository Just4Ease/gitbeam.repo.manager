package server

import (
	"context"
	"gitbeam.baselib/utils"
	"gitbeam.repo.manager/contract/repos"
	"gitbeam.repo.manager/core"
	"gitbeam.repo.manager/models"
	"github.com/sirupsen/logrus"
)

type apiService struct {
	service *core.GitBeamService
	logger  *logrus.Logger
}

func (a apiService) HealthCheck(ctx context.Context, void *gitRepos.Void) (*gitRepos.HealthCheckResponse, error) {
	return &gitRepos.HealthCheckResponse{Code: 200}, nil
}

func (a apiService) ListGitRepositories(ctx context.Context, request *gitRepos.Void) (*gitRepos.ListGitRepositoriesResponse, error) {
	output, err := a.service.ListRepos(ctx)
	if err != nil {
		return nil, err
	}

	var repos []*gitRepos.Repo
	_ = utils.UnPack(output, &repos)
	return &gitRepos.ListGitRepositoriesResponse{Repos: repos}, nil
}

func (a apiService) GetGitRepo(ctx context.Context, request *gitRepos.GetGitRepoRequest) (*gitRepos.Repo, error) {
	output, err := a.service.GetByOwnerAndRepoName(ctx, &models.OwnerAndRepoName{
		OwnerName: request.OwnerName,
		RepoName:  request.RepoName,
	})
	if err != nil {
		return nil, err
	}

	var repo gitRepos.Repo
	_ = utils.UnPack(output, &repo)
	return &repo, nil
}

func NewApiService(core *core.GitBeamService, logger *logrus.Logger) gitRepos.GitBeamRepositoryServiceServer {
	return &apiService{
		service: core,
		logger:  logger,
	}
}
