package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/ddddddO/memo/api/adapter"
)

type memoUsecase interface {
	List(userID int, tagID int, status string) ([]adapter.Memo, error)
	Detail(memoID int) (*adapter.Memo, error)
	Update(updatedMemo adapter.Memo) error
	Create(createdMemo adapter.Memo) error
	Delete(memoID int) error
}

type memoHandler struct {
	usecase memoUsecase
}

func NewMemoHandler(uc memoUsecase) *memoHandler {
	return &memoHandler{
		usecase: uc,
	}
}

func (h *memoHandler) List(w http.ResponseWriter, r *http.Request) {
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
	tagID := params.Get("tagId")
	tid, err := strconv.Atoi(tagID)
	if err != nil {
		tid = -1
	}
	status := params.Get("status")

	ams, err := h.usecase.List(uid, tid, status)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

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

func (h *memoHandler) Detail(w http.ResponseWriter, r *http.Request) {
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

	am, err := h.usecase.Detail(mid)
	if err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
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

	if err := h.usecase.Update(updatedMemo); err != nil {
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

	if err := h.usecase.Create(createdMemo); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

	if err := h.usecase.Delete(mid); err != nil {
		errResponse(w, http.StatusInternalServerError, "failed", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
