package usecase

import (
	"github.com/ddddddO/tag-mng/repository"

	"github.com/pkg/errors"
)

type healthUsecase struct {
	healthRepo repository.HealthRepository
}

func NewHealth(healthRepo repository.HealthRepository) *healthUsecase {
	return &healthUsecase{
		healthRepo: healthRepo,
	}
}

func (u *healthUsecase) Check() error {
	if err := u.healthRepo.Check(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
