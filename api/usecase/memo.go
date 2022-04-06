package usecase

import (
	"database/sql"
	"log"
	"sort"

	"github.com/pkg/errors"

	"github.com/ddddddO/memo/api/adapter"
	"github.com/ddddddO/memo/models"
)

type memoRepository interface {
	FetchList(userID int) ([]*models.Memo, error)
	FetchListByTagID(userID, tagID int) ([]*models.Memo, error)
	Fetch(memoID int) (*models.Memo, error)
	Update(memo *models.Memo, tagIDs []int) error
	Create(memo *models.Memo, tagIDs []int) error
	Delete(memoID int) error
}

type memoUsecase struct {
	repo    memoRepository
	tagRepo tagRepository
}

func NewMemoUsecase(repo memoRepository, tagRepo tagRepository) *memoUsecase {
	return &memoUsecase{
		repo:    repo,
		tagRepo: tagRepo,
	}
}

func (u *memoUsecase) List(userID int, tagID int, status string) ([]adapter.Memo, error) {
	isExposed := false
	if status == "exposed" {
		isExposed = true
	}

	var (
		memos []*models.Memo
		err   error
	)
	if tagID == -1 {
		memos, err = u.repo.FetchList(userID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		memos, err = u.repo.FetchListByTagID(userID, tagID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if isExposed {
		memos = filterExposed(memos)
	}

	ams := make([]adapter.Memo, len(memos))
	for i, mm := range memos {
		tags, err := u.tagRepo.FetchListByMemoID(int(mm.ID))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		ats := make([]adapter.Tag, len(tags))
		for i, t := range tags {
			at := adapter.Tag{
				ID:   t.ID,
				Name: t.Name,
			}
			ats[i] = at
		}

		am := adapter.Memo{
			ID:          mm.ID,
			Subject:     mm.Subject,
			Content:     mm.Content,
			IsExposed:   mm.IsExposed.Bool,
			UserID:      int(mm.UsersID.Int64),
			Tags:        ats,
			NotifiedCnt: int(mm.NotifiedCnt.Int64),
			CreatedAt:   &mm.CreatedAt.Time,
			UpdatedAt:   &mm.UpdatedAt.Time,
			ExposedAt:   &mm.ExposedAt.Time,
		}
		setColor(mm, &am)
		ams[i] = am
	}

	// NOTE: NotifiedCntでメモを昇順にソート
	sort.SliceStable(ams,
		func(i, j int) bool {
			return ams[i].NotifiedCnt < ams[j].NotifiedCnt
		},
	)

	return ams, nil
}

func filterExposed(memos []*models.Memo) []*models.Memo {
	var mm []*models.Memo
	for _, m := range memos {
		if !m.IsExposed.Valid {
			continue
		}
		if m.IsExposed.Bool {
			mm = append(mm, m)
		}
	}
	return mm
}

func setColor(mm *models.Memo, am *adapter.Memo) {
	switch int(mm.NotifiedCnt.Int64) {
	case 0:
		am.RowVariant = "danger"
	case 1:
		am.RowVariant = "warning"
	case 2:
		am.RowVariant = "primary"
	case 3:
		am.RowVariant = "info"
	case 4:
		am.RowVariant = "secondary"
	case 5:
		am.RowVariant = "success"
	}
}

func (u *memoUsecase) Detail(memoID int) (*adapter.Memo, error) {
	memo, err := u.repo.Fetch(memoID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tags, err := u.tagRepo.FetchListByMemoID(memoID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ats := make([]adapter.Tag, len(tags))
	for i, t := range tags {
		at := adapter.Tag{
			ID:   t.ID,
			Name: t.Name,
		}
		ats[i] = at
	}

	am := &adapter.Memo{
		ID:          memo.ID,
		Subject:     memo.Subject,
		Content:     memo.Content,
		IsExposed:   memo.IsExposed.Bool,
		UserID:      int(memo.UsersID.Int64),
		Tags:        ats,
		NotifiedCnt: int(memo.NotifiedCnt.Int64),
		CreatedAt:   &memo.CreatedAt.Time,
		UpdatedAt:   &memo.UpdatedAt.Time,
		ExposedAt:   &memo.ExposedAt.Time,
	}

	return am, nil
}

func (u *memoUsecase) Update(updatedMemo adapter.Memo) error {
	memo, err := u.repo.Fetch(updatedMemo.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	memo.Subject = updatedMemo.Subject
	memo.Content = updatedMemo.Content
	memo.IsExposed = sql.NullBool{
		Bool:  updatedMemo.IsExposed,
		Valid: true,
	}

	tagIDs := make([]int, len(updatedMemo.Tags))
	for i, tag := range updatedMemo.Tags {
		tagIDs[i] = tag.ID
	}

	if err := u.repo.Update(memo, tagIDs); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *memoUsecase) Create(createdMemo adapter.Memo) error {
	memo := &models.Memo{
		Subject: createdMemo.Subject,
		Content: createdMemo.Content,
		IsExposed: sql.NullBool{
			Bool:  createdMemo.IsExposed,
			Valid: true,
		},
		UsersID: sql.NullInt64{
			Int64: int64(createdMemo.UserID),
			Valid: true,
		},
	}
	tagIDs := make([]int, len(createdMemo.Tags))
	for i, tag := range createdMemo.Tags {
		tagIDs[i] = tag.ID
	}

	if err := u.repo.Create(memo, tagIDs); err != nil {
		log.Println("failed to create memo", err)
		return errors.WithStack(err)
	}

	return nil
}

func (u *memoUsecase) Delete(memoID int) error {
	if err := u.repo.Delete(memoID); err != nil {
		log.Println("failed to delete memo", err)
		return errors.WithStack(err)
	}
	return nil
}
