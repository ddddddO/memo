package usecase

import (
	"github.com/pkg/errors"
)

type healthRepository interface {
	Check() error
}

type healthUsecase struct {
	repo healthRepository
}

func NewHealthUsecase(repo healthRepository) *healthUsecase {
	return &healthUsecase{
		repo: repo,
	}
}

func (u *healthUsecase) Ping() error {
	if err := u.repo.Check(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
