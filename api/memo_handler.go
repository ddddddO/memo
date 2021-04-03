package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/domain"
	"github.com/ddddddO/tag-mng/repository"
)

type memoHandler struct {
	repo repository.MemoRepository
}

func NewMemoHandler(repo repository.MemoRepository) *memoHandler {
	return &memoHandler{
		repo: repo,
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
	memos, err := h.repo.FetchList(uid, tid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	res := struct {
		Memos []domain.Memo `json:"memo_list"`
	}{
		Memos: memos,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
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
	mid, err := strconv.Atoi(memoID)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	memo, err := h.repo.Fetch(uid, mid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	//ref: https://qiita.com/shohei-ojs/items/311ef080cd5cff1e0e16
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(memo); err != nil {
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

	updatedMemo := domain.Memo{
		ID: mid,
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed to unmarshal json", err)
		return
	}

	if err := h.repo.Update(updatedMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *memoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createdMemo domain.Memo
	if err := json.NewDecoder(r.Body).Decode(&createdMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	if err := h.repo.Create(createdMemo); err != nil {
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

	deleteMemo := domain.Memo{
		ID: mid,
	}
	if err := json.NewDecoder(r.Body).Decode(&deleteMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	if err := h.repo.Delete(deleteMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
