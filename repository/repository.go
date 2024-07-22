package repository

import (
	"context"
	"gitbeam.repo.manager/models"
)

//go:generate mockgen -source=repository.go -destination=../mocks/data_store_mock.go -package=mocks
type DataStore interface {
	StoreRepository(ctx context.Context, payload *models.Repo) error
	ListRepos(context.Context) ([]*models.Repo, error)
	GetRepoByOwner(ctx context.Context, owner *models.OwnerAndRepoName) (*models.Repo, error)
}
