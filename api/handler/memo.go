package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/tag-mng/domain"
)

type memoUsecase interface {
	FetchList(userID int, tagID int) ([]domain.Memo, error)
	Fetch(userID int, memoID int) (domain.Memo, error)
	Update(domain.Memo) error
	Create(domain.Memo) error
	Delete(domain.Memo) error
}

type memoHandler struct {
	usecase memoUsecase
}

func NewMemo(usecase memoUsecase) *memoHandler {
	return &memoHandler{
		usecase: usecase,
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

	memos, err := h.usecase.FetchList(uid, tid)
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

	memo, err := h.usecase.Fetch(uid, mid)
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

	if err := h.usecase.Update(updatedMemo); err != nil {
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

	if err := h.usecase.Create(createdMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	// TODO: 201 createdへ変更
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

	if err := h.usecase.Delete(deleteMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	// TODO: 204 no content?
	w.WriteHeader(http.StatusOK)
}
