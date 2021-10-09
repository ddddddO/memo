package usecase

import (
	"github.com/ddddddO/tag-mng/domain"
	"github.com/ddddddO/tag-mng/repository"

	"github.com/pkg/errors"
)

type memoUsecase struct {
	memoRepo repository.MemoRepository
}

func NewMemoUsecase(memoRepo repository.MemoRepository) *memoUsecase {
	return &memoUsecase{
		memoRepo: memoRepo,
	}
}

func (u *memoUsecase) FetchList(userID int, tagID int) ([]domain.Memo, error) {
	memos, err := u.memoRepo.FetchList(userID, tagID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return memos, nil
}

func (u *memoUsecase) Fetch(userID int, memoID int) (domain.Memo, error) {
	memo, err := u.memoRepo.Fetch(userID, memoID)
	if err != nil {
		return domain.Memo{}, errors.WithStack(err)
	}
	return memo, nil
}

func (u *memoUsecase) Update(memo domain.Memo) error {
	if err := u.memoRepo.Update(memo); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *memoUsecase) Create(memo domain.Memo) error {
	if err := u.memoRepo.Create(memo); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (u *memoUsecase) Delete(memo domain.Memo) error {
	if err := u.memoRepo.Delete(memo); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
