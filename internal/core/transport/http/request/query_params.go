package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not a valid integer: %v: %w",
			param, key, err, core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

const dateLayout = "2006-01-02"

// GetDateFromQueryParam парсит дату как начало календарного дня в time.Local
// (main.go выставляет time.Local из TIME_ZONE).
func GetDateFromQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	date, err := time.ParseInLocation(dateLayout, param, time.Local)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not a valid date: %v: %w",
			param, key, err, core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

// GetDateToQueryParam парсит верхнюю границу диапазона: начало следующего дня
// в time.Local (для условия task_created_at < to).
func GetDateToQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	date, err := time.ParseInLocation(dateLayout, param, time.Local)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not a valid date: %v: %w",
			param, key, err, core_errors.ErrInvalidArgument,
		)
	}

	endExclusive := date.AddDate(0, 0, 1)
	return &endExclusive, nil
}
