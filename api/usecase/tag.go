package usecase

import (
	"github.com/ddddddO/tag-mng/domain"
	"github.com/ddddddO/tag-mng/repository"

	"github.com/pkg/errors"
)

type tagUsecase struct {
	tagRepo repository.TagRepository
}

func NewTag(tagRepo repository.TagRepository) *tagUsecase {
	return &tagUsecase{
		tagRepo: tagRepo,
	}
}

func (u *tagUsecase) FetchList(userID int) ([]domain.Tag, error) {
	tags, err := u.tagRepo.FetchList(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tags, nil
}

func (u *tagUsecase) Fetch(tagID int) (domain.Tag, error) {
	tag, err := u.tagRepo.Fetch(tagID)
	if err != nil {
		return domain.Tag{}, errors.WithStack(err)
	}
	return tag, nil
}

func (u *tagUsecase) Update(tag domain.Tag) error {
	if err := u.tagRepo.Update(tag); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *tagUsecase) Create(tag domain.Tag) error {
	if err := u.tagRepo.Create(tag); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *tagUsecase) Delete(tag domain.Tag) error {
	if err := u.tagRepo.Delete(tag); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
