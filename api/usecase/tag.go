package usecase

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	"github.com/ddddddO/memo/api/adapter"
	"github.com/ddddddO/memo/models"
)

type tagRepository interface {
	FetchList(userID int) ([]*models.Tag, error)
	FetchListByMemoID(memoID int) ([]*models.Tag, error)
	Fetch(tagID int) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(tagID int) error
	Create(tag *models.Tag) error
}

type tagUsecase struct {
	repo tagRepository
}

func NewTagUsecase(repo tagRepository) *tagUsecase {
	return &tagUsecase{
		repo: repo,
	}
}

func (u *tagUsecase) List(userID int) ([]adapter.Tag, error) {
	tags, err := u.repo.FetchList(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	atags := make([]adapter.Tag, len(tags))
	for i, tag := range tags {
		atags[i] = adapter.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return atags, nil
}

func (u *tagUsecase) Detail(tagID int) (*adapter.Tag, error) {
	tag, err := u.repo.Fetch(tagID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	atag := &adapter.Tag{
		ID:   tag.ID,
		Name: tag.Name,
	}

	return atag, nil
}

func (u *tagUsecase) Update(updatedTag adapter.Tag) error {
	tag, err := u.repo.Fetch(updatedTag.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	tag.Name = updatedTag.Name

	if err := u.repo.Update(tag); err != nil {
		log.Println("failed to update tag", err)
		return errors.WithStack(err)
	}

	return nil
}

func (u *tagUsecase) Delete(tagID int) error {
	if err := u.repo.Delete(tagID); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *tagUsecase) Create(createTag adapter.Tag) error {
	tag := &models.Tag{
		Name: createTag.Name,
		UsersID: sql.NullInt64{
			Int64: int64(createTag.UserID),
			Valid: true,
		},
	}
	if err := u.repo.Create(tag); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
