package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"practice5/model"
	"practice5/repository"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) GetPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	q := r.URL.Query()

	page := parseIntDefault(q.Get("page"), 1)
	pageSize := parseIntDefault(q.Get("page_size"), 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	f := model.FilterParams{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  strings.TrimSpace(q.Get("order_by")),
	}

	if v := q.Get("id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		f.ID = &id
	}
	if v := q.Get("name"); v != "" {
		f.Name = &v
	}
	if v := q.Get("email"); v != "" {
		f.Email = &v
	}
	if v := q.Get("gender"); v != "" {
		f.Gender = &v
	}
	if v := q.Get("birth_date"); v != "" {
		f.BirthDate = &v
	}

	result, err := h.repo.GetPaginatedUsers(f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	q := r.URL.Query()

	uid1, err := strconv.Atoi(q.Get("user_id1"))
	if err != nil || uid1 < 1 {
		writeError(w, http.StatusBadRequest, "invalid user_id1")
		return
	}
	uid2, err := strconv.Atoi(q.Get("user_id2"))
	if err != nil || uid2 < 1 {
		writeError(w, http.StatusBadRequest, "invalid user_id2")
		return
	}
	if uid1 == uid2 {
		writeError(w, http.StatusBadRequest, "user_id1 and user_id2 must be different")
		return
	}

	friends, err := h.repo.GetCommonFriends(uid1, uid2)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id1":       uid1,
		"user_id2":       uid2,
		"common_friends": friends,
		"count":          len(friends),
	})
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
