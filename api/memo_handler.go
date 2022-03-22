package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/adapter"
	"github.com/ddddddO/memo/models"
	"github.com/ddddddO/memo/repository"
)

type memoHandler struct {
	repo    repository.MemoRepository
	tagRepo repository.TagRepository
}

func NewMemoHandler(repo repository.MemoRepository, tagRepo repository.TagRepository) *memoHandler {
	return &memoHandler{
		repo:    repo,
		tagRepo: tagRepo,
	}
}

func (h *memoHandler) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	userID := params.Get("userId")
	if len(userID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
		return
	}
	tagID := params.Get("tagId")

	uid, err := strconv.Atoi(userID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
	tid, err := strconv.Atoi(tagID)
	if err != nil {
		tid = -1
	}

	var memos []*models.Memo
	if tid == -1 {
		memos, err = h.repo.FetchList(uid)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}
	} else {
		memos, err = h.repo.FetchListByTagID(uid, tid)
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
		}
	}

	ams := make([]adapter.Memo, len(memos))
	for i, mm := range memos {
		tags, err := h.tagRepo.FetchListByMemoID(int(mm.ID))
		if err != nil {
			errResponse(w, http.StatusInternalServerError, "failed", err)
			return
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

	res := struct {
		Memos []adapter.Memo `json:"memo_list"`
	}{
		Memos: ams,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
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

func (h *memoHandler) Detail(w http.ResponseWriter, r *http.Request) {
	memoID := chi.URLParam(r, "id")
	if len(memoID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
		return
	}

	params := r.URL.Query()
	userID := params.Get("userId")
	if len(userID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'userId'", nil)
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
	_ = uid // FIXME: 必要？
	mid, err := strconv.Atoi(memoID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	memo, err := h.repo.Fetch(mid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
	tags, err := h.tagRepo.FetchListByMemoID(mid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
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

	//ref: https://qiita.com/shohei-ojs/items/311ef080cd5cff1e0e16
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(am); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}
}

func (h *memoHandler) Update(w http.ResponseWriter, r *http.Request) {
	memoID := chi.URLParam(r, "id")
	if len(memoID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
		return
	}
	mid, err := strconv.Atoi(memoID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	updatedMemo := adapter.Memo{
		ID: mid,
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed to unmarshal json", err)
		return
	}

	memo, err := h.repo.Fetch(mid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
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

	if err := h.repo.Update(memo, tagIDs); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *memoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createdMemo adapter.Memo
	if err := json.NewDecoder(r.Body).Decode(&createdMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	memo := &models.Memo{
		Subject: createdMemo.Subject,
		Content: createdMemo.Content,
		UsersID: sql.NullInt64{
			Int64: int64(createdMemo.UserID),
			Valid: true,
		},
	}
	tagIDs := make([]int, len(createdMemo.Tags))
	for i, tag := range createdMemo.Tags {
		tagIDs[i] = tag.ID
	}

	if err := h.repo.Create(memo, tagIDs); err != nil {
		log.Println("failed to create memo", err)
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *memoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	memoID := chi.URLParam(r, "id")
	if len(memoID) == 0 {
		errResponse(w, http.StatusBadRequest, "empty value 'memoId'", nil)
		return
	}
	mid, err := strconv.Atoi(memoID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	// deleteMemo := adapter.Memo{
	// 	ID: mid,
	// }
	// if err := json.NewDecoder(r.Body).Decode(&deleteMemo); err != nil {
	// 	errResponse(w, http.StatusInternalServerError, "failed", err)
	// 	return
	// }

	if err := h.repo.Delete(mid); err != nil {
		log.Println("failed to delete memo", err)
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
