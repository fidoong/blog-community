package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response[T any] struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

type PaginatedResponse[T any] struct {
	List       []T        `json:"list"`
	Pagination Pagination `json:"pagination"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Success[T any](w http.ResponseWriter, data T) {
	JSON(w, http.StatusOK, Response[T]{
		Code:      "OK",
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

func Fail(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, Response[any]{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}
