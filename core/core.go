package core

import (
	"context"
	"encoding/json"
	"errors"
	"gitbeam.baselib/store"
	"gitbeam.repo.manager/events/topics"
	"gitbeam.repo.manager/models"
	"gitbeam.repo.manager/repository"
	"github.com/google/go-github/v63/github"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	ErrGithubRepoNotFound       = errors.New("github repo not found")
	ErrOwnerAndRepoNameRequired = errors.New("owner and repo name required")
)

type GitBeamService struct {
	githubClient *github.Client
	logger       *logrus.Logger
	dataStore    repository.DataStore
	eventStore   store.EventStore
}

func NewGitBeamService(
	logger *logrus.Logger,
	eventStore store.EventStore,
	dataStore repository.DataStore,
	httpClient *http.Client, // Nullable.
) *GitBeamService {
	client := github.NewClient(httpClient) // Didn't need to pass this as a top level dependency into the git beam service.
	return &GitBeamService{
		githubClient: client,
		dataStore:    dataStore,
		eventStore:   eventStore,
		logger:       logger.WithField("serviceName", "GitBeamService").Logger,
	}
}

func (g GitBeamService) GetEventStore() store.EventStore {
	return g.eventStore
}

func (g GitBeamService) ListRepos(ctx context.Context) ([]*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "ListRepos")
	list, err := g.dataStore.ListRepos(ctx)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to list repositories")
		return make([]*models.Repo, 0), nil
	}
	return list, nil
}

func (g GitBeamService) GetByOwnerAndRepoName(ctx context.Context, owner *models.OwnerAndRepoName) (*models.Repo, error) {
	useLogger := g.logger.WithContext(ctx).WithField("methodName", "GetByOwnerAndRepoName")

	if err := owner.Validate(); err != nil {
		useLogger.WithError(err).Error("owner is invalid. please provide a valid ownerName and repoName")
		return nil, ErrOwnerAndRepoNameRequired
	}

	existingRepo, err := g.dataStore.GetRepoByOwner(ctx, owner)
	if err == nil && existingRepo != nil {
		return existingRepo, nil
	}

	gitRepo, _, err := g.githubClient.Repositories.Get(ctx, owner.OwnerName, owner.RepoName)
	if err != nil {
		useLogger.WithError(err).Errorln("failed to get repo by owner and repo name from github api")
		return nil, ErrGithubRepoNotFound
	}

	repo := &models.Repo{
		Name:          gitRepo.GetName(),
		Owner:         gitRepo.GetOwner().GetLogin(),
		Description:   gitRepo.GetDescription(),
		URL:           gitRepo.GetHTMLURL(),
		Languages:     gitRepo.GetLanguage(),
		ForkCount:     gitRepo.GetForksCount(),
		StarCount:     gitRepo.GetStargazersCount(),
		OpenIssues:    gitRepo.GetOpenIssues(),
		WatchersCount: gitRepo.GetWatchersCount(),
		TimeCreated:   gitRepo.GetCreatedAt().Time.Format(time.RFC3339),
		TimeUpdated:   gitRepo.GetUpdatedAt().Time.Format(time.RFC3339),
		Meta:          gitRepo.String(),
	}

	if err := g.dataStore.StoreRepository(ctx, repo); err != nil {
		useLogger.WithError(err).Errorln("Failed to persist repository")
		return nil, err
	}

	data, err := json.Marshal(repo)
	if err != nil {
		useLogger.WithError(err).Errorln("json.Marshal: failed to marshal repo before publishing to event store.")
		return nil, err
	}

	// This is a channel-based event store, so checking of errors aren't needed here.
	_ = g.eventStore.Publish(topics.RepoCreated, data)

	return repo, nil
}
